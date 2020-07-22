package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/howeyc/fsnotify"
)

var goPaths []string
var watcher *fsnotify.Watcher

func init() {
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	goPaths = strings.Split(os.Getenv("GOPATH"), separator)

	// init watcher
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if !ev.IsAttrib() && strings.HasSuffix(ev.Name, ".go") {
					log.Printf("sending event %s\n", ev)
					eventChannel <- ev.String()
				}
			case err := <-watcher.Error:
				log.Fatalln(err)
			}
		}
	}()
}

func getDepFolders() []string {
	deps := strings.Split(strings.Trim(goList(mainPkg, `{{ join .Deps "\n" }}`), "\n"), "\n")
	var watchedFolders []string
	appendWatchedFolders := func(path string) bool {
		e, err := exists(path)
		if err != nil {
			log.Fatalln(err)
		}
		if e {
			watchedFolders = append(watchedFolders, path)
			return true
		}
		return false
	}
	for _, dep := range deps {
		path := fmt.Sprintf("vendor/%s", dep)
		ok, err := exists(path)
		if err != nil {
			log.Fatalln(err)
		}
		if ok {
			continue
		}
		for _, gopath := range goPaths {
			path := fmt.Sprintf("%s/src/%s", gopath, dep)
			if !strings.HasPrefix(path, mainPkgDir) && appendWatchedFolders(path) {
				break
			}
		}
	}

	appendWatchedFolders(mainPkgDir)
	return watchedFolders
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func watch() {
	watchedFolders := getDepFolders()
	for _, folder := range watchedFolders {
		if *verbose {
			log.Println("recursive watch:", folder)
		}
		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				err := watcher.Watch(path)
				if err != nil {
					log.Fatalln(err)
				}
			}
			return nil
		})
		err := watcher.Watch(folder)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
