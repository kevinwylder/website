package main

import (
    "log"
)

func (server *serverState) logUrl(path, method string) {
    _, err := server.db.Exec(`
INSERT INTO requests (url, method) VALUES (?, ?);
    `, path, method)
    if err != nil {
        log.Println(err)
    }
}
