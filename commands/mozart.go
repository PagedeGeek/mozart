package main

import "github.com/pagedegeek/mozart"

func main() {

        jsonConfig := `{
                "listen_info": "localhost:1357"
        }`

        if config, err := mozart.NewConfigFromJSON([]byte(jsonConfig)); err != nil {
                panic(err)
        } else {
                s := mozart.NewService(config)
                if err := s.Start(); err != nil {
                        panic(err)
                }
        }
}