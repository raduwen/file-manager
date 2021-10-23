package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	workers "github.com/jrallison/go-workers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type File struct {
	Path   string  `gorm:"primaryKey"`
	IsDir  bool    `gorm:"index"`
	Ext    *string `gorm:"index"`
	Sha256 *string `gorm:"index"`
}

func (file *File) BeforeCreate(tx *gorm.DB) error {
	info, _ := os.Stat(file.Path)
	file.IsDir = info.IsDir()
	if file.IsDir {
		return nil
	}

	f, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	sha := fmt.Sprintf("%x", h.Sum(nil))
	file.Sha256 = &sha
	ext := filepath.Ext(file.Path)
	file.Ext = &ext

	return nil
}

func indexPath(msg *workers.Msg) {
	a, _ := msg.Args().Array()
	p, err := filepath.Abs(a[0].(string))
	if err != nil {
		log.Println(err)
	}
	// [TODO] - upsert
	db.Create(&File{Path: p})
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("index.db"), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}

	db.AutoMigrate(&File{})

	workers.Configure(map[string]string{
		"server":   "localhost:6379",
		"database": "0",
		"pool":     "30",
		"process":  "1",
	})

	workers.Process("indexPath", indexPath, 1)

	// go workers.StatsServer(8080)
	workers.Run()
}
