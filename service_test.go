package mozart

import (
        "testing"
        "github.com/stretchr/testify/assert"
        "mes_tests/mozart"
)

var config *mozart.Config
// var service *mozart.Service
func init() {
        config = &mozart.Config{ ListenInfo: "localhost:1357" }
}

func TestNewService(t *testing.T) {
        s := mozart.NewService(config)
        assert.NotNil(t, s)
}

func TestAddAndRemoveTask(t *testing.T) {
        service := mozart.NewService(config)

        assert.Equal(t, 0, service.GetCount())

        task := mozart.NewTask()
        assert.NotNil(t, task)

        service.AddTask(task)
        assert.Equal(t, 1, service.GetCount())
        assert.NotNil(t, task.Service)
        assert.Equal(t, service, task.Service)

        service.RemoveTask(task.UUID)
        assert.Equal(t, 0, service.GetCount())
}

