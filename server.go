package main

import (
	"flag"
	"fmt"
    "os"
	"log"
	"net/http"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func main() {
	port := flag.Int("port", 8080, "the port to serve on")
	path := flag.String("path", "./static", "the path to serve static content")
	flag.Parse()

    fmt.Println("Starting Server")
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), CreateServer(*path))
    if err != nil {
        log.Fatal(err)
    }
}

type serverState struct {
	static http.Handler
    db *sql.DB
}

// CreateServer creates an http handler for the website
func CreateServer(static string) http.Handler {
	files := http.FileServer(http.Dir(static))
    db, err := sql.Open("mysql", os.Getenv("DB_CREDS"))
    if err != nil {
        log.Fatal(err)
    }
	return &serverState{
		static: files,
        db: db,
	}
}

func (server *serverState) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    go server.logUrl(req.URL.Path, req.Method)
	server.static.ServeHTTP(res, req)
}
