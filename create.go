package main

import (
	"context"
	"encoding/json"

	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/kv"
	"github.com/influxdata/influxdb/task/backend"
)

func createTask(store kv.Store, svc *kv.Service, task *influxdb.Task) error {

	org, err := svc.FindOrganizationByID(context.Background(), task.OrganizationID)
	if err != nil {
		return err
	}
	task.Organization = org.Name

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if task.Status == "" {
		task.Status = string(backend.TaskActive)
	}

	return store.Update(context.Background(), func(tx kv.Tx) error {
		taskBucket, err := tx.Bucket([]byte("tasksv1"))
		if err != nil {
			return err
		}

		indexBucket, err := tx.Bucket([]byte("taskIndexsv1"))
		if err != nil {
			return err
		}

		taskKey, err := taskKey(task.ID)
		if err != nil {
			return err
		}
		orgKey, err := taskOrgKey(task.OrganizationID, task.ID)
		if err != nil {
			return err
		}

		// write the task
		err = taskBucket.Put(taskKey, taskBytes)
		if err != nil {
			return err
		}

		// write the org index
		err = indexBucket.Put(orgKey, taskKey)
		if err != nil {
			return err
		}
		return nil
	})
}

func taskKey(taskID influxdb.ID) ([]byte, error) {
	encodedID, err := taskID.Encode()
	if err != nil {
		return nil, err
	}
	return encodedID, nil
}

func taskOrgKey(orgID, taskID influxdb.ID) ([]byte, error) {
	encodedOrgID, err := orgID.Encode()
	if err != nil {
		return nil, err
	}
	encodedID, err := taskID.Encode()
	if err != nil {
		return nil, err
	}

	return []byte(string(encodedOrgID) + "/" + string(encodedID)), nil
}
