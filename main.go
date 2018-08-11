package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var mainPkg string
var mainBin string
var mainGoPath string

var m = flag.String("main", "", "main package name")
var wd = flag.String("wd", "", "working directory")
var args = flag.String("args", "", "args")
var delayBuild = flag.Int("delay", 1000, "delay *ms before rebuild")
var buildFlags = flag.String("build", "", "build flags")

func main() {
	flag.Parse()
	mainPkg = *m
	for _, goPath := range goPaths {
		if _, err := os.Stat(goPath + "/src/" + mainPkg); os.IsNotExist(err) {
			continue
		}
		mainGoPath = goPath
		break
	}
	if mainGoPath == "" {
		log.Fatalln(mainPkg, "not found")
	}
	mainBin = mainGoPath + "/bin/" + path.Base(mainPkg)

	autoBuild()
	startChannel <- ""
	<-make(chan int)
}
