package storage

// much is owed to https://raw.githubusercontent.com/victorcoder/dkron/master/dkron/store.go for this implementation

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/sirupsen/logrus"

	"github.com/byxorna/flow/config"
	"github.com/byxorna/flow/types/execution"
	"github.com/byxorna/flow/types/job"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "storage"})
)

// MaxExecutions is how many executions to retain in the storage backend
const MaxExecutions = 100

// Store abstracts the DAL
type Store struct {
	// Client is the libkv client
	Client   store.Store
	keyspace string
	backend  store.Backend
}

func init() {
	etcd.Register()
	consul.Register()
	zookeeper.Register()
}

// New returns a new storage backend
func New(c config.Config) (*Store, error) {
	if len(c.EtcdEndpoints) == 0 {
		return nil, fmt.Errorf("No supported storage backend in Config")
	}
	//TODO update this if we wanna support multiple backends. For now, idgaf
	backend := store.ETCD
	machines := c.EtcdEndpoints
	keyspace := ""
	if c.EtcdPrefix != "/" {
		// "" is fine for keyspace if we have a default prefix to make joining work
		keyspace = c.EtcdPrefix
	}
	cfg := c.ToLibKVConfig()

	s, err := libkv.NewStore(store.Backend(backend), machines, &cfg)
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
	}, nil
}

// String returns a string for a Store
func (s *Store) String() string {
	return fmt.Sprintf("%s in %s", s.backend, s.keyspace)
}

// SetJob Stores a job
func (s *Store) SetJob(j *job.Spec) error {
	// Sanitize the job name
	j.ID.Name = generateSlug(j.ID.Name)
	jobKey := j.Path(s.keyspace)
	log.Debug("Storing %s to %s", j.ID, jobKey)

	if err := j.Validate(); err != nil {
		return err
	}

	// Get if the requested job already exist
	ej, err := s.GetJob(j.ID)
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
		"job":       j.ID.Name,
		"namespace": j.ID.Namespace,
		"json":      string(jobJSON),
	}).Debug("store: Setting job")

	err = s.Client.Put(jobKey, jobJSON, nil)
	return err
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

// GetJobs returns all jobs
func (s *Store) GetJobs() ([]*job.Spec, error) {
	log.Debugf("fetching keyspace %s/%s/", s.keyspace, job.StoragePath)
	res, err := s.Client.List(s.keyspace + "/" + job.StoragePath + "/")
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Debug("store: No jobs found")
			return []*job.Spec{}, nil
		}
		return nil, err
	}
	log.Debugf("got res %s", res)

	jobs := make([]*job.Spec, 0)
	for _, node := range res {
		var j job.Spec
		log.Infof("attempting to unmarshal %s as Job", node.Value)
		err := json.Unmarshal([]byte(node.Value), &j)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &j)
	}
	return jobs, nil
}

// GetJob ...
func (s *Store) GetJob(id job.ID) (*job.Spec, error) {
	res, err := s.Client.Get(job.Prefix(s.keyspace, id))
	if err != nil {
		return nil, err
	}

	var j job.Spec
	if err = json.Unmarshal([]byte(res.Value), &j); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"name":      j.ID.Name,
		"namespace": j.ID.Namespace,
	}).Debug("store: Retrieved job from datastore")

	return &j, nil
}

// DeleteJob ...
func (s *Store) DeleteJob(id job.ID) (*job.Spec, error) {
	j, err := s.GetJob(id)
	if err != nil {
		return nil, err
	}

	if err := s.DeleteExecutions(id); err != nil {
		if err != store.ErrKeyNotFound {
			return nil, err
		}
	}

	if err := s.Client.Delete(job.Prefix(s.keyspace, id)); err != nil {
		return nil, err
	}

	return j, nil
}

// GetExecutions ...
func (s *Store) GetExecutions(id job.ID) ([]*execution.Instance, error) {
	prefix := execution.Path(s.keyspace, id)
	res, err := s.Client.List(prefix)
	if err != nil {
		return nil, err
	}

	var executions []*execution.Instance

	for _, node := range res {
		if store.Backend(s.backend) != store.ZK {
			path := store.SplitKey(node.Key)
			nsDir := path[len(path)-2]
			jobDir := path[len(path)-1]
			if nsDir != id.Namespace {
				continue
			}
			if jobDir != id.Name {
				continue
			}
		}
		var e execution.Instance
		err := json.Unmarshal([]byte(node.Value), &e)
		if err != nil {
			return nil, err
		}
		executions = append(executions, &e)
	}
	return executions, nil
}

// GetLastExecutionGroup ...
func (s *Store) GetLastExecutionGroup(id job.ID) ([]*execution.Instance, error) {
	prefix := execution.Path(s.keyspace, id)
	res, err := s.Client.List(prefix)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return []*execution.Instance{}, nil
	}

	var ex execution.Instance
	err = json.Unmarshal([]byte(res[len(res)-1].Value), &ex)
	if err != nil {
		return nil, err
	}
	return s.GetExecutionGroup(&ex)
}

// GetExecutionGroup ...
func (s *Store) GetExecutionGroup(e *execution.Instance) ([]*execution.Instance, error) {
	res, err := s.Client.List(execution.Path(s.keyspace, e.Job))
	if err != nil {
		return nil, err
	}

	var executions []*execution.Instance
	for _, node := range res {
		var ex execution.Instance
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

// GetGroupedExecutions Returns executions for a job grouped and with an ordered index
// to facilitate access.
func (s *Store) GetGroupedExecutions(id job.ID) (map[int64][]*execution.Instance, []int64, error) {
	execs, err := s.GetExecutions(id)
	if err != nil {
		return nil, nil, err
	}
	groups := make(map[int64][]*execution.Instance)
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

// SetExecution Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(e *execution.Instance) (string, error) {
	exJSON, _ := json.Marshal(e)

	log.WithFields(logrus.Fields{
		"job":       e.Job.Name,
		"namespace": e.Job.Namespace,
		"key":       e.ID,
	}).Debug("store: Setting execution instance")

	err := s.Client.Put(
		fmt.Sprintf(
			"%s/%s",
			execution.Path(s.keyspace, e.Job),
			e.ID),
		exJSON,
		nil,
	)
	if err != nil {
		return "", err
	}

	execs, err := s.GetExecutions(e.Job)
	if err != nil {
		log.Errorf("store: No executions found for job %s", e.Job)
		return "", err
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
			err := s.Client.Delete(
				fmt.Sprintf(
					"%s/%s",
					execution.Path(s.keyspace, execs[i].Job),
					execs[i].ID,
				),
			)
			if err != nil {
				log.Errorf("store: Trying to delete overflowed execution %s", execs[i].ID)
			}
		}
	}

	return e.ID.String(), nil
}

// DeleteExecutions Removes all executions of a job
func (s *Store) DeleteExecutions(id job.ID) error {
	return s.Client.DeleteTree(execution.Path(s.keyspace, id))
}

// GetLeader Retrieve the leader from the store
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

// LeaderKey Retrieve the leader key used in the KV store to store the leader node
func (s *Store) LeaderKey() string {
	return s.keyspace + "/leader"
}
