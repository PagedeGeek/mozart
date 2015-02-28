package mozart

import (
        "os"
        "encoding/json"
        "io/ioutil"
)

type ActionWriteFile struct {
        Filename string `json:"filename"`
        OsFileMode int `json:"filemode"`
        JsonBody *json.RawMessage `json:"json_body"`
        task *Task
}

func (a *ActionWriteFile) Exec() (chan bool, chan error) {
        executedChan, errorChan := make(chan bool), make(chan error)

        go func() {
                if a.OsFileMode == 0 {
                        a.OsFileMode = 0644
                }
                if err := ioutil.WriteFile(a.Filename, *a.JsonBody, os.FileMode(a.OsFileMode)); err != nil {
                        errorChan <- err
                } else {
                        executedChan <- true
                }
        }()

        return executedChan, errorChan
}

func NewActionWriteFile(task *Task) *ActionWriteFile {
        return &ActionWriteFile{task: task}
}


