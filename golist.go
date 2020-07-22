package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func goList(pkg, template string) string {
	cmd := exec.Command("go", "list", "-f", template, pkg)

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

	text, _ := ioutil.ReadAll(stdout)
	io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}

	return string(text)
}
