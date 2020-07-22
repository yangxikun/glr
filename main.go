package main

import (
	"flag"
	"os"
	"path"
	"strings"
)

var mainPkg string
var mainPkgDir string
var mainBin string
var mainGoPath string

var m = flag.String("main", "", "main package name")
var wd = flag.String("wd", "", "working directory")
var args = flag.String("args", "", "args")
var delayBuild = flag.Int("delay", 1000, "delay *ms before rebuild")
var buildFlags = flag.String("build", "", "build flags")
var verbose = flag.Bool("verbose", false, "more log")

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
		mainGoPath = goPaths[0]
	}
	mainBin = mainGoPath + "/bin/" + path.Base(strings.TrimSpace(goList(mainPkg, `{{ .Module.Path }}`)))
	mainPkgDir = strings.TrimSpace(goList(mainPkg, `{{ .Module.Dir }}`))

	autoBuild()
	eventChannel <- ""
	<-make(chan int)
}
