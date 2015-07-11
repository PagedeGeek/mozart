package mozart

import (
        "encoding/json"
        "bytes"
        "net/http"
        "fmt"
)

type ActionHttpRequest struct {
        Url string `json:"url"`
        Verb string `json:"verb"`
        Header map[string]string `json:"header"`
        JsonBody *json.RawMessage `json:"json_body"`
        task *Task
}

func (a *ActionHttpRequest) Exec() (chan bool, chan error) {

        executedChan, errorChan := make(chan bool), make(chan error)

        go func() {
                if a.Verb == "" {
                        a.Verb = "POST"
                }

                httpClient := &http.Client{}
                if httpReq, err := http.NewRequest(a.Verb, a.Url, bytes.NewBuffer(*a.JsonBody)); err != nil {
                        errorChan <- err
                } else {
                        httpReq.Header.Set("Content-Type", "application/json")

                        for k, v := range a.Header {
                                httpReq.Header.Set(k, v)
                        }

                        if resp, err := httpClient.Do(httpReq); err != nil {
                                errorChan <- err
                        } else {
                                fmt.Printf("Response code: %d\n", resp.StatusCode)
                                executedChan <- true
                        }
                }
        }()
        
        return executedChan, errorChan
}

func NewActionHttpRequest(task *Task) *ActionHttpRequest {
        return &ActionHttpRequest{task: task}
}

