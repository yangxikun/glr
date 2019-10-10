package main

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/kballard/go-shellquote"
)

func run() error {
	log.Println("Running...")
	log.Println(mainBin, *args)

	cmdArga, err := shellquote.Split(*args)
	if err != nil {
		return err
	}
	cmd := exec.Command(mainBin, cmdArga...)
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
		err := cmd.Process.Kill()
		if err != nil {
			log.Panicln("Killing err", err)
		}
		state, err := cmd.Process.Wait()
		if err != nil {
			log.Panicln("Wait err", err)
		}
		log.Println(state)
		waitChannel <- struct{}{}
	}()
	return nil
}
