package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/kballard/go-shellquote"
)

var eventChannel = make(chan string, 7)
var stopChannel = make(chan struct{})
var waitChannel = make(chan struct{})

func autoBuild() {
	started := false
	go func() {
		for {
			eventName := <-eventChannel

			log.Println("receiving first event", eventName)
			log.Printf("sleeping for %d milliseconds\n", *delayBuild)
			time.Sleep(time.Duration(*delayBuild) * time.Millisecond)
			log.Println("flushing events")

			flushEvents()

			buildFail := false
			err := build(mainPkg)
			if err != nil {
				log.Printf("Build Failed:\n%s", err)
				buildFail = true
			}

			if !buildFail && started {
				stopChannel <- struct{}{}
				<-waitChannel
			}
			watch()

			if !buildFail {
				err = run()
				if err != nil {
					log.Println(err)
					started = false
					continue
				}
				started = true
			}
		}
	}()
}

func flushEvents() {
	for {
		select {
		case eventName := <-eventChannel:
			log.Println("receiving event", eventName)
		default:
			return
		}
	}
}

func build(mainPkg string) error {
	log.Println("Building...")

	bFlags, err := shellquote.Split(*buildFlags)
	if err != nil {
		log.Println(err)
	}
	cmdArgs := []string{"install"}
	cmdArgs = append(cmdArgs, bFlags...)
	cmdArgs = append(cmdArgs, mainPkg)
	cmd := exec.Command("go", cmdArgs...)
	// don't use pipe, seems cmd pipe will stuck
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
