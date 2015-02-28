package mozart

import (
        "encoding/json"
)

type Config struct {
        ListenInfo string `json:"listen_info"`
}

func NewConfigFromJSON(jsonRaw []byte) (*Config, error) {
        var config Config
        if err := json.Unmarshal(jsonRaw, &config); err != nil {
                return nil, err
        }
        return &config, nil
}

