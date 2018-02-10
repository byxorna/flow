# Flow - Periodic Job Execution

super WIP. just throwing things at the wall and seeing if something sticks. I am unhappy
with the design of airflow, and wanted to try my hand at something

* more container native
* more extensible
  * supports arbitrary executors like k8s, mesos, mapreduce, etc
* keeps business logic out of job definitions

# TODO

- [ ] define job spec schema
  - [ ] support N instances?
  - [ ] templating job?
- [ ] define job status
  - [ ] fields for last fire, first fire, completion
- [ ] define API
  - [ ] job spec submission API
  - [ ] job status API
  - [ ] job delete API
  - [ ] ...
- [ ] define scheduling projection - how do we know which jobs should get queued up next interval?
  - [ ] includes missed timers
- [ ] define schema for etcd
  - [ ] how will we store and retrieve models?
  - [ ] define DAL abstraction
- [ ] handle missed job firing - (replay? skip? track jitter of submission?)
- [ ] metrics exposition
  - [ ] how jittery is a handler for expected vs actual submission?
  - [ ] track job completion histograms/launch times?

# Notes

## API

## Data Model

## Scheduling

TODO: what happens if we miss a firing of a job? (flow down? executor backpressure?)

## Braindump

DAG scheduler system
like airflow, but not shitty
description language for creating job dags
jobs have constraints/labels
queue system for workers who advertise some capabilities/labels
support for external systems like docker containers, mapreduce, k8s
API driven storage for jobs
https://landing.google.com/sre/book/chapters/distributed-periodic-scheduling.html

• pluggable executors
    • local, shell, docker, kubernetes, mapreduce
• API driven job submission, support .d directory too
• label based job targeting
• backed by etcd for storage
• use etcd for queue?

• API stores state in etcd
• executors handle processing jobs and shipping to backend executor systems
• output is stored in executors native system
    • does executor have system to extract output? (nah, too complex)
    • link to logs though... should implement "get logs for this job instance"
