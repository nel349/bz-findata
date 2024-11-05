package scheduler

import (
	// "context"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nel349/bz-findata/internal/analysis/application/ports"
	"github.com/nel349/bz-findata/internal/analysis/task"
	"github.com/robfig/cron/v3"
)

type TaskManager struct {
	cron    *cron.Cron
	tasks   map[cron.EntryID]scheduler.Task
	mutex   sync.RWMutex
	service *task.Service
}

func NewTaskManager(service *task.Service) *TaskManager {
	c := cron.New(cron.WithSeconds())
	c.Start() // Start the cron scheduler

	return &TaskManager{
		cron:    c,
		tasks:   make(map[cron.EntryID]scheduler.Task),
		service: service,
	}
}

func (tm *TaskManager) StartTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Schedule string `json:"schedule"`
		Hours    int    `json:"hours"`
		Limit    int    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	id, err := tm.cron.AddFunc(req.Schedule, func() {
		// Log a task is running
		log.Printf("Running scheduled task at %s - Hours: %d, Limit: %d",
			time.Now().Format(time.RFC3339), req.Hours, req.Limit)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := tm.service.StoreMatchOrders(ctx, req.Hours, req.Limit); err != nil {
			log.Printf("Error executing scheduled task: %v", err)
		}
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task := scheduler.Task{
		ID:       id,
		Schedule: req.Schedule,
		Hours:    req.Hours,
		Limit:    req.Limit,
	}

	tm.tasks[id] = task

	log.Printf("Cron task scheduled with ID %d at %s", id, time.Now().Format(time.RFC3339))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (tm *TaskManager) StopTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	id, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if _, exists := tm.tasks[cron.EntryID(id)]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	tm.cron.Remove(cron.EntryID(id))
	delete(tm.tasks, cron.EntryID(id))
	log.Printf("Cron task with ID %d stopped at %s", id, time.Now().Format(time.RFC3339))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Task %s stopped successfully", taskID),
	})
}

func (tm *TaskManager) ListTasks(w http.ResponseWriter, r *http.Request) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	tasks := make([]scheduler.Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
