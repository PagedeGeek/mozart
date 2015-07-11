package mozart

import (
        "testing"
        "github.com/stretchr/testify/assert"
        "mes_tests/mozart"
)

func TestNewConfigFromJSON(t *testing.T) {
        jsonRaw := `{
                "listen_info": "localhost:1357"
        }`

        config, err := mozart.NewConfigFromJSON([]byte(jsonRaw))
        assert.Nil(t, err)
        assert.NotNil(t, config)

        assert.Equal(t, "localhost:1357", config.ListenInfo)
}

