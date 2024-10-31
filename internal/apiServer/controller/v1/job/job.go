package job

import clientv3 "go.etcd.io/etcd/client/v3"

type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}