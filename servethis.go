package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		dirname := os.Args[1]
		pwd = filepath.Join(pwd, dirname)
		dir, err := os.Stat(pwd)
		if err != nil {
			panic(err)
		}
		if !dir.IsDir() {
			fmt.Fprintln(os.Stderr, "Usage:", filepath.Base(os.Args[0]),
				"[dir]")
			os.Exit(1)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request: ", r.Method, r.URL)
		path := r.URL.Path
		localPath := filepath.Join(pwd, path)
		fi, err := os.Stat(localPath)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if fi.IsDir() {
			index := filepath.Join(localPath, "index.html")
			if _, err := os.Stat(index); err == nil {
				http.ServeFile(w, r, index)
				return
			}
		}
		http.ServeFile(w, r, localPath)
	})

	bind := ":8080"
	fmt.Fprintln(os.Stderr, "Starting server on", bind, "for", pwd)
	log.Fatal(http.ListenAndServe(bind, nil))
}
