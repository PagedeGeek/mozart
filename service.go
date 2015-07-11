package mozart

import (
        "net/http"
        "github.com/gorilla/mux"
        "encoding/json"
        "fmt"
        "sync"
        "errors"
        "github.com/satori/go.uuid"
)

type Service struct {
        config *Config
        sync.Mutex
        tasks map[uuid.UUID]*Task
        tasksCount int
}

func NewService(config *Config) (*Service) {
        return &Service{ config: config }
}

func (s *Service) AddTask(t *Task) error {
        s.Lock()
        defer s.Unlock()

        if s.tasks == nil {
                s.tasks = make(map[uuid.UUID]*Task)
        }

        if _, exists := s.tasks[t.UUID]; exists {
                return errors.New("Task uuid already exists")
        }

        t.SetService(s)
        s.tasks[t.UUID] = t
        s.tasksCount += 1
        return nil
}

func (s *Service) UnscheduleTask(_uuid uuid.UUID) error {
        s.Lock()
        defer s.Unlock()

        if s.tasks == nil {
                s.tasks = make(map[uuid.UUID]*Task)
        }

        task, exists := s.tasks[_uuid]
        if !exists {
                return errors.New("Task uuid do not exists")
        }

        task.Unschedule()
        return nil
}

func (s *Service) RemoveTask(_uuid uuid.UUID) {
        s.Lock()
        defer s.Unlock()

        if s.tasks == nil {
                s.tasks = make(map[uuid.UUID]*Task)
        }

        delete(s.tasks, _uuid)
        s.tasksCount -= 1
}

func (s *Service) GetCount() int {
        s.Lock()
        defer s.Unlock()
        return s.tasksCount
}

func (s *Service) SignalTaskTerminated(_uuid uuid.UUID) {
        s.RemoveTask(_uuid)
}

func (s *Service) Start() error {
        router := mux.NewRouter()
        router.HandleFunc("/tasks/schedule", s.handleScheduleTask).Methods("POST")
        router.HandleFunc("/tasks/unschedule/{uuid}", s.handleUnscheduleTask).Methods("DELETE")
        router.HandleFunc("/tasks/count", s.handleCountTasks).Methods("GET")
        router.HandleFunc("/tasks", s.handleListTasks).Methods("GET")
        if err := http.ListenAndServe(s.config.ListenInfo, router); err != nil {
                return err
        }
        return nil
}

type RespMsg struct {
        Message string `json:"message"`
        TaskUUID uuid.UUID `json:"task_uuid,omitempty"`
        Errors []string `json:"errors"`
}

func (s *Service) handleScheduleTask(w http.ResponseWriter, req *http.Request) {
        defer req.Body.Close()
        w.Header().Set("Content-Type", "application/json")
        respMsg := &RespMsg{}
        var statusHeader int

        if task, errList:= NewTaskFromJSON(req.Body); errList != nil {
                statusHeader = http.StatusBadRequest
                respMsg.Message = "errors"
                respMsg.Errors = make([]string, 0)
                for _, err := range errList {
                        fmt.Printf("Task error: %s\n", err)
                        respMsg.Errors = append(respMsg.Errors, err.Error())
                }
        } else {
                s.AddTask(task)
                task.AsyncExec()
                respMsg.Message = "OK"
                respMsg.TaskUUID = task.UUID
                statusHeader = http.StatusAccepted
        }

        if respJson, err := json.Marshal(respMsg); err != nil {
                http.Error(w, err.Error(), 500)
        } else {
                w.WriteHeader(statusHeader)
                w.Write(respJson)
        }
}

func (s *Service) handleUnscheduleTask(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        rawUUID := vars["uuid"]
        if _uuid, err := uuid.FromString(rawUUID); err != nil {
                http.Error(w, err.Error(), 500)
        } else {
                s.UnscheduleTask(_uuid)
                w.WriteHeader(http.StatusCreated)
        }
}

func (s *Service) handleListTasks(w http.ResponseWriter, req *http.Request) {
        acc := make([]*Task, 0)
        for _, t := range s.tasks {
                acc = append(acc, t)
        }

        if rawResp, err := json.Marshal(acc); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusAccepted)
                w.Write(rawResp)
        }
}

func (s *Service) handleCountTasks(w http.ResponseWriter, req *http.Request) {
        respMsg := struct {
                Count int `json:"count"`
        }{
                Count: s.GetCount(),
        }
        if rawResp, err := json.Marshal(respMsg); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write(rawResp)
        }
}

