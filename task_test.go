package mozart

import (
        "testing"
        "github.com/stretchr/testify/assert"
        "mes_tests/mozart"
        "time"
)

var config *mozart.Config
var service *mozart.Service
func init() {
        config = &mozart.Config{ ListenInfo: "localhost:1357" }
        service = mozart.NewService(config)
}

func TestNewBaseTask(t *testing.T) {
        task := mozart.NewBaseTask()
        assert.NotNil(t, task)
        assert.NotNil(t, task.GetUUID())

        task.SetService(service)
        assert.NotNil(t, task.GetService())
}

func TestParseDelay(t *testing.T) {
        task := mozart.NewBaseTask()
        task.In = "3s"
        err := task.ParseDelay()
        assert.Nil(t, err)
        assert.NotNil(t, task.GetParsedDelay())
        assert.Equal(t, time.Duration(3 * time.Second), task.GetParsedDelay())

        task2 := mozart.NewBaseTask()
        task2.In = "aa"
        err = task2.ParseDelay()
        assert.NotNil(t, err)
}

func TestAsyncExec(t *testing.T) {
        task := mozart.NewBaseTask()
        task.SetService(service)
        task.In = "5s"
        task.ParseDelay()
        assert.Nil(t, task.GetStopChannel())
        task.AsyncExec()
        assert.NotNil(t, task.GetStopChannel())

        task.Unschedule()
        _, closed := <-task.GetStopChannel()
        assert.Equal(t, true, closed)
}

