package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var startChannel = make(chan string)
var stopChannel = make(chan struct{})

func autoBuild() {
	started := false
	go func() {
		for {
			eventName := <-startChannel

			log.Println("receiving first event", eventName)
			log.Printf("sleeping for %d milliseconds\n", *delayBuild)
			time.Sleep(time.Duration(*delayBuild) * time.Millisecond)
			log.Println("flushing events")

			flushEvents()

			buildFailed := false
			errorMessage, ok := build(mainPkg)
			if !ok {
				buildFailed = true
				log.Printf("Build Failed:\n%s", errorMessage)
			}

			if buildFailed {
				continue
			}
			if started {
				stopChannel <- struct{}{}
			}
			run()
			watch()
			started = true
		}
	}()
}

func flushEvents() {
	for {
		select {
		case eventName := <-startChannel:
			log.Println("receiving event", eventName)
		default:
			return
		}
	}
}

func build(mainPkg string) (string, bool) {
	log.Println("Building...")

	cmd := exec.Command("go", "install", "-x", mainPkg)

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

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}
