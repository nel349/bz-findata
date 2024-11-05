package scheduler

import "github.com/robfig/cron/v3"

// Task represents a scheduled task
type Task struct {
	ID       cron.EntryID
	Schedule string
	Hours    int
	Limit    int
}

// Scheduler defines the interface for task scheduling operations
type Scheduler interface {
	StartTask(schedule string, hours, limit int) (Task, error)
	StopTask(taskID cron.EntryID) error
	ListTasks() ([]Task, error)
}
