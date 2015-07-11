package mozart

type Action interface {
        Exec() (chan bool, chan error)
}

