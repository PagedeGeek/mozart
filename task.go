package mozart

import (
        "github.com/satori/go.uuid"
        "time"
        "encoding/json"
        "fmt"
        "io"
        "errors"
)

type Task struct {
        UUID uuid.UUID `json:"uuid"`
        In string `json:"in"`
        Timeout string `json:"timeout"`
        Do string `json:"do"`
        RawParams json.RawMessage `json:"params"`

        service *Service
        parsedDelay time.Duration
        parsedTimeout time.Duration
        action Action
        unscheduleChannel chan bool
        executedChan chan bool
        errorChan chan error
}

var uuidChannel chan uuid.UUID

func generateUUID() uuid.UUID {
        if uuidChannel == nil {
                uuidChannel = make(chan uuid.UUID)
                go func() {
                        for { uuidChannel <-uuid.NewV4() }
                }()
        }
        return <-uuidChannel
}

func NewTask() *Task {
        return &Task{ UUID: generateUUID() }
}

func NewTaskFromJSON(reader io.ReadCloser) (*Task, []error) {
        errList := make([]error, 0)

        task := NewTask()
        decoder := json.NewDecoder(reader)
        if err := decoder.Decode(task); err != nil {
                errList = append(errList, err)
        } else {
                for _, err:= range task.prepare() {
                        errList = append(errList, err)
                }
        }

        if len(errList) > 0 {
                return nil, errList
        } else {
                return task, nil
        }
}

func (t *Task) SetService(s *Service) {
        t.service = s
}

func (t *Task) Unschedule() {
        t.unscheduleChannel <-true
}

func (t *Task) ParseDelay() (error) {
        delay, err := time.ParseDuration(t.In)
        if err != nil {
                return err
        }
        t.parsedDelay = delay
        return nil
}

func (t *Task) ParseTimeout() (error) {
        if t.Timeout == "" {
                t.Timeout = "15s"
        }
        delay, err := time.ParseDuration(t.Timeout)
        if err != nil {
                return err
        }
        t.parsedTimeout = delay
        return nil
}

func (t *Task) prepare() []error {
        errList := make([]error, 0)
        if err := t.ParseDelay(); err != nil {
                errList = append(errList, err)
        }

        if err := t.ParseTimeout(); err != nil {
                errList = append(errList, err)
        }

        switch t.Do {
        case "http_request":
                t.action = NewActionHttpRequest(t)
        case "write_file":
                t.action = NewActionWriteFile(t)
        default:
                err := errors.New(fmt.Sprintf("Unknow action: '%s'", t.Do))
                errList = append(errList, err)
        }

        if t.action != nil {
                if err := json.Unmarshal(t.RawParams, t.action); err != nil {
                        errList = append(errList, err)
                }
        }

        return errList
}

func (t *Task) Finalize() {
        close(t.unscheduleChannel)
        if t.executedChan != nil { close(t.executedChan) }
        if t.errorChan != nil { close(t.errorChan) }
}

func (t *Task) AsyncExec() {
        fmt.Printf("Task created. Executed in %q %s\n", t.In, t.UUID)
        t.unscheduleChannel = make(chan bool)
        go func() {
                select {
                case <-t.unscheduleChannel:
                        fmt.Printf("Task unscheduled: %s\n", t.UUID)
                case <-time.After(t.parsedDelay):
                        fmt.Printf("Task exec: %s\n", t.UUID)
                        t.executedChan, t.errorChan = t.action.Exec()
                        select {
                        case <-time.After(t.parsedTimeout):
                                fmt.Printf("Task timeout: %s\n", t.UUID)
                        case err := <-t.errorChan:
                                fmt.Printf("Action error: %s\n", err)
                        case <-t.executedChan:
                                fmt.Printf("Task executed: %s\n", t.UUID)
                        }
                }
                t.Finalize()
                t.service.SignalTaskTerminated(t.UUID)
        }()
}

