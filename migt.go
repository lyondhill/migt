package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/influxdata/idpe/etcd"
	_ "github.com/influxdata/idpe/query/builtin"
	tasketcd "github.com/influxdata/idpe/task/store/etcd"
	"github.com/influxdata/influxdb/kv"
	"github.com/influxdata/influxdb/task/backend"
)

func main() {
	etcdConfig := etcd.Config{
		Hosts:       []string{"http://localhost:2379"},
		DialTimeout: time.Minute,
	}
	kvStore := etcd.NewKVStore(etcdConfig)
	if err := kvStore.Open(); err != nil {
		panic(err)
	}
	svc := kv.NewService(kvStore)
	if err := svc.Initialize(context.Background()); err != nil {
		panic(err)
	}

	cl, err := clientv3.New(clientv3.Config{Endpoints: []string{"http://localhost:2379"}, DialTimeout: time.Minute})
	if err != nil {
		panic(err)
	}
	st, err := tasketcd.NewEtcdStore(cl, 30)
	if err != nil {
		panic(err)
	}

	tasks, err := st.ListTasks(context.Background(), backend.TaskSearchParams{})

	for len(tasks) > 0 {
		for _, tm := range tasks {
			// fmt.Printf("old: %+v\n", tm)
			t, err := toPlatformTask(tm.Task, &tm.Meta)
			if err != nil {
				panic(err)
			}
			fmt.Printf("task: %+v\n", t)

			if err := createTask(kvStore, svc, t); err != nil {
				panic(err)
			}
			fmt.Printf("newt: %+v\n", t)
		}

		tasks, err = st.ListTasks(context.Background(), backend.TaskSearchParams{After: tasks[len(tasks)-1].Task.ID})
	}
}
