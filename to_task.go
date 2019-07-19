package main

import (
	"time"

	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/task/backend"
	"github.com/influxdata/influxdb/task/options"
)

func toPlatformTask(t backend.StoreTask, m *backend.StoreTaskMeta) (*influxdb.Task, error) {
	opts, err := options.FromScript(t.Script)
	if err != nil {
		return nil, err
	}

	pt := &influxdb.Task{
		ID:             t.ID,
		OrganizationID: t.Org,
		Name:           t.Name,
		Flux:           t.Script,
		Cron:           opts.Cron,
	}
	if !opts.Every.IsZero() {
		pt.Every = opts.Every.String()
	}
	if opts.Offset != nil && !(*opts.Offset).IsZero() {
		pt.Offset = opts.Offset.String()
	}
	if m != nil {
		pt.Status = string(m.Status)
		pt.LatestCompleted = time.Unix(m.LatestCompleted, 0).Format(time.RFC3339)
		if m.CreatedAt != 0 {
			pt.CreatedAt = time.Unix(m.CreatedAt, 0).Format(time.RFC3339)
		}
		if m.UpdatedAt != 0 {
			pt.UpdatedAt = time.Unix(m.UpdatedAt, 0).Format(time.RFC3339)
		}
		pt.AuthorizationID = influxdb.ID(m.AuthorizationID)
	}
	return pt, nil
}
