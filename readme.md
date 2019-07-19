# Task Migration

Migrate tasks from old task store to new kv store

## Expectations

During the migration tasks will exist in both systems. This migration is not transactional but it is repeatable and reversable.

We can expect to lose our run history. Migrating that data would be slow and has a much greater potential to lose data and cause the migration to fail.

It is possible for a short period of time the task could be run by both the old system and the new. This can only happen if the deployment in k8s hangs during deploy.

## Required Steps

1. Merge the PR that integrates the new system [3537](https://github.com/influxdata/idpe/pull/3537)
2. Get the image of the currently running task service (in case we need to roll back).
3. Get the image that needs to be pushed to the live servers.
4. Establish port forward to the production etcd server.
5. Ensure new task data location is empty. `etcdctl get --prefix=true tasksv1`
6. Run Migration tool.
7. Deploy new task image `two-prod set image deployment/tasks tasks=quay.io/influxdb/tasks:<tag>`
8. Run Migration tool again (ensure data is up to date).
9. Manually confirm thats are running using a user logged into production and by viewing task service logs.

## Failure Procedure

1. Roll back to previous taskd image. `two-prod rollout undo deploy/task`
2. clean out cruft created by migration tool. `etcdctl del --prefix=true tasksv1`