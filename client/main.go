package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	workers "github.com/jrallison/go-workers"
)

func main() {
	workers.Configure(map[string]string{
		"server":   "localhost:6379",
		"database": "0",
		"pool":     "30",
		"process":  "1",
	})

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: client [PATH]\n")
		return
	}

	filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		p, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		log.Printf("queue ... %s", p)
		workers.Enqueue("indexPath", "indexPath", []string{p})
		return nil
	})
}
