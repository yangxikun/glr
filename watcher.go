package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
	cmd := exec.Command("go", "list", "-f", `{{ join .Deps "\n" }}`, mainPkg)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	depsText, _ := ioutil.ReadAll(stdout)
	io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}

	deps := strings.Split(strings.Trim(string(depsText), "\n"), "\n")
	var watchedFolders []string
	for _, dep := range deps {
		for _, gopath := range goPaths {
			path := fmt.Sprintf("%s/src/%s", gopath, dep)
			e, err := exists(path)
			if err != nil {
				log.Fatalln(err)
			}
			if e {
				watchedFolders = append(watchedFolders, path)
				break
			}
		}
	}

	watchedFolders = append(watchedFolders, mainGoPath+"/src/"+mainPkg)
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
		err := watcher.Watch(folder)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
