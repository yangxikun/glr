package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func run() {
	log.Println("Running...")

	cmd := exec.Command(mainBin)
	cmd.Dir = *wd

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

	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdout, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		log.Println("Killing PID", pid)
		cmd.Process.Kill()
	}()
}
