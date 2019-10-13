package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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
	secret []byte
	static http.Handler
	db     *sql.DB
}

// CreateServer creates an http handler for the website
func CreateServer(static string) http.Handler {
	files := http.FileServer(http.Dir(static))
	db, err := sql.Open("mysql", os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
	}
	if secret, exists := os.LookupEnv("SECRET"); exists {
		return &serverState{
			static: files,
			db:     db,
			secret: []byte(secret),
		}
	}
	log.Fatal("Missing HMAC secret")
	return nil
}

func (server *serverState) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	go server.logURL(req.URL.Path, req.Method)
	if req.URL.Path == "/_" && req.Method == "POST" {
		// Handle Event posts
		http.Error(res, "", server.OnEvent(req.Body))
		return
	}
	server.static.ServeHTTP(res, req)
}
