package storage

// much is owed to https://raw.githubusercontent.com/victorcoder/dkron/master/dkron/store.go for this implementation

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"

	"github.com/byxorna/flow/types/execution"
	"github.com/byxorna/flow/types/job"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "storage"})
)

// MaxExecutions is how many executions to retain in the storage backend
const MaxExecutions = 100

type Store struct {
	// Client is the libkv client
	Client   store.Store
	keyspace string
	backend  string
}

func init() {
	etcd.Register()
	consul.Register()
	zookeeper.Register()
}

func NewStore(backend string, machines []string, keyspace string) *Store {
	s, err := libkv.NewStore(store.Backend(backend), machines, nil)
	if err != nil {
		log.Error(err)
	}

	log.WithFields(logrus.Fields{
		"backend":  backend,
		"machines": machines,
		"keyspace": keyspace,
	}).Debug("store: Backend config")

	_, err = s.List(keyspace)
	if err != store.ErrKeyNotFound && err != nil {
		log.WithError(err).Fatal("store: Store backend not reachable")
	}

	return &Store{
		Client:   s,
		keyspace: keyspace,
		backend:  backend,
	}
}

// Store a job
func (s *Store) SetJob(j *job.Spec) error {
	// Sanitize the job name
	j.Name = generateSlug(job.Name)
	jobKey := fmt.Sprintf("%s/jobs/%s", s.keyspace, j.Name)

	if err := s.validateJob(j); err != nil {
		return err
	}

	// Get if the requested job already exist
	ej, err := s.GetJob(j.Namespace, j.Name)
	if err != nil && err != store.ErrKeyNotFound {
		return err
	}
	if ej != nil {
		// When the job runs, these status vars are updated
		// otherwise use the ones that are stored
		if ej.LastError.After(j.LastError) {
			j.LastError = ej.LastError
		}
		if ej.LastSuccess.After(j.LastSuccess) {
			j.LastSuccess = ej.LastSuccess
		}
		if ej.SuccessCount > j.SuccessCount {
			j.SuccessCount = ej.SuccessCount
		}
		if ej.ErrorCount > j.ErrorCount {
			j.ErrorCount = ej.ErrorCount
		}
	}

	jobJSON, _ := json.Marshal(j)

	log.WithFields(logrus.Fields{
		"job":  j.Name,
		"json": string(jobJSON),
	}).Debug("store: Setting job")

	if err := s.Client.Put(jobKey, jobJSON, nil); err != nil {
		return err
	}

	return nil
}

/*
// Set the depencency tree for a job given the job and the previous version
// of the Job or nil if it's new.
func (s *Store) SetJobDependencyTree(j *job.Spec, previousJob *job.Spec) error {
	// Existing job that doesn't have parent job set and it's being set
	if previousJob != nil && previousJob.ParentJob == "" && j.ParentJob != "" {
		pj, err := j.GetParent()
		if err != nil {
			return err
		}
		pj.Lock()
		defer pj.Unlock()

		pj.DependentJobs = append(pj.DependentJobs, j.Name)
		if err := s.SetJob(pj); err != nil {
			return err
		}
	}

	// Existing job that has parent job set and it's being removed
	if previousJob != nil && previousJob.ParentJob != "" && j.ParentJob == "" {
		pj, err := previousJob.GetParent()
		if err != nil {
			return err
		}
		pj.Lock()
		defer pj.Unlock()

		ndx := 0
		for i, djn := range pj.DependentJobs {
			if djn == j.Name {
				ndx = i
				break
			}
		}
		pj.DependentJobs = append(pj.DependentJobs[:ndx], pj.DependentJobs[ndx+1:]...)
		if err := s.SetJob(pj); err != nil {
			return err
		}
	}

	// New job that has parent job set
	if previousJob == nil && job.ParentJob != "" {
		pj, err := j.GetParent()
		if err != nil {
			return err
		}
		pj.Lock()
		defer pj.Unlock()

		pj.DependentJobs = append(pj.DependentJobs, j.Name)
		if err := s.SetJob(pj); err != nil {
			return err
		}
	}

	return nil
}
*/

func (s *Store) validateJob(j *job.Spec) error {
	// TODO perform validation on Spec before storage
	if j.ParentJob == j.Name {
		return ErrSameParent
	}

	// Only validate the schedule if it doesn't have a parent
	if j.ParentJob == "" {
		if err := j.Schedule.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// GetJobs returns all jobs
func (s *Store) GetJobs() ([]*job.Spec, error) {
	res, err := s.Client.List(s.keyspace + "/jobs/")
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Debug("store: No jobs found")
			return []*job.Spec{}, nil
		}
		return nil, err
	}

	jobs := make([]*job.Spec, 0)
	for _, node := range res {
		var j job.Spec
		err := json.Unmarshal([]byte(node.Value), &j)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &j)
	}
	return jobs, nil
}

// Get a job
func (s *Store) GetJob(ns string, name string) (*job.Spec, error) {
	res, err := s.Client.Get(fmt.Sprintf("%s/jobs/%s/%s", s.keyspace, ns, name))
	if err != nil {
		return nil, err
	}

	var j job.Spec
	if err = json.Unmarshal([]byte(res.Value), &j); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"job": j.Name,
	}).Debug("store: Retrieved job from datastore")

	return &j, nil
}

func (s *Store) DeleteJob(ns string, name string) (*job.Spec, error) {
	j, err := s.GetJob(ns, name)
	if err != nil {
		return nil, err
	}

	if err := s.DeleteExecutions(ns, name); err != nil {
		if err != store.ErrKeyNotFound {
			return nil, err
		}
	}

	if err := s.Client.Delete(fmt.Sprintf("%s/jobs/%s/%s", s.keyspace, ns, name)); err != nil {
		return nil, err
	}

	return j, nil
}

func (s *Store) GetExecutions(ns string, name string) ([]*execution.Execution, error) {
	prefix := fmt.Sprintf("%s/executions/%s/%s", s.keyspace, ns, name)
	res, err := s.Client.List(prefix)
	if err != nil {
		return nil, err
	}

	var executions []*execution.Execution

	for _, node := range res {
		if store.Backend(s.backend) != store.ZK {
			path := store.SplitKey(node.Key)
			dir := path[len(path)-2]
			if dir != jobName {
				continue
			}
		}
		var e execution.Execution
		err := json.Unmarshal([]byte(node.Value), &e)
		if err != nil {
			return nil, err
		}
		executions = append(executions, &e)
	}
	return executions, nil
}

func (s *Store) GetLastExecutionGroup(ns string, j string) ([]*execution.Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, ns, j))
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return []*execution.Execution{}, nil
	}

	var ex execution.Execution
	err = json.Unmarshal([]byte(res[len(res)-1].Value), &ex)
	if err != nil {
		return nil, err
	}
	return s.GetExecutionGroup(&ex)
}

func (s *Store) GetExecutionGroup(e *execution.Execution) ([]*execution.Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, e.Namespace, e.Job))
	if err != nil {
		return nil, err
	}

	var executions []*execution.Execution
	for _, node := range res {
		var ex execution.Execution
		err := json.Unmarshal([]byte(node.Value), &ex)
		if err != nil {
			return nil, err
		}

		if ex.Group == e.Group {
			executions = append(executions, &ex)
		}
	}
	return executions, nil
}

// Returns executions for a job grouped and with an ordered index
// to facilitate access.
func (s *Store) GetGroupedExecutions(ns string, j string) (map[int64][]*execution.Execution, []int64, error) {
	execs, err := s.GetExecutions(ns, j)
	if err != nil {
		return nil, nil, err
	}
	groups := make(map[int64][]*execution.Execution)
	for _, exec := range execs {
		groups[exec.Group] = append(groups[exec.Group], exec)
	}

	// Build a separate data structure to show in order
	var byGroup int64arr
	for key := range groups {
		byGroup = append(byGroup, key)
	}
	sort.Sort(sort.Reverse(byGroup))

	return groups, byGroup, nil
}

// Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(e *execution.Execution) (string, error) {
	exJson, _ := json.Marshal(e)
	key := e.Key()

	log.WithFields(logrus.Fields{
		"job":       e.Job,
		"namespace": e.Namespace,
		"execution": key,
	}).Debug("store: Setting key")

	err := s.Client.Put(fmt.Sprintf("%s/executions/%s/%s/%s", s.keyspace, e.Namespace, e.Job, key), exJson, nil)
	if err != nil {
		return "", err
	}

	execs, err := s.GetExecutions(e.Namespace, e.Job)
	if err != nil {
		log.Errorf("store: No executions found for job %s", execution.JobName)
	}

	// Get and ordered array of all execution groups
	var byGroup int64arr
	for _, ex := range execs {
		byGroup = append(byGroup, ex.Group)
	}
	sort.Sort(byGroup)

	// Delete all execution results over the limit, starting from olders
	if len(byGroup) > MaxExecutions {
		for i := range byGroup[MaxExecutions:] {
			err := s.Client.Delete(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, execs[i].JobName, execs[i].Key()))
			if err != nil {
				log.Errorf("store: Trying to delete overflowed execution %s", execs[i].Key())
			}
		}
	}

	return key, nil
}

// Removes all executions of a job
func (s *Store) DeleteExecutions(ns string, j string) error {
	return s.Client.DeleteTree(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, ns, j))
}

// Retrieve the leader from the store
func (s *Store) GetLeader() []byte {
	res, err := s.Client.Get(s.LeaderKey())
	if err != nil {
		if err == store.ErrNotReachable {
			log.Fatal("store: Store not reachable, be sure you have an existing key-value store running is running and is reachable.")
		} else if err != store.ErrKeyNotFound {
			log.Error(err)
		}
		return nil
	}

	log.WithField("node", string(res.Value)).Debug("store: Retrieved leader from datastore")

	return res.Value
}

// Retrieve the leader key used in the KV store to store the leader node
func (s *Store) LeaderKey() string {
	return s.keyspace + "/leader"
}