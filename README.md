# Flow - Periodic Job Execution


# Notes

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