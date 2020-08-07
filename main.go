package main

import (
	"flag"
	"log"
	"path"
	"path/filepath"
	"strings"
)

var mainPkg string
var mainPkgDir string
var mainBin string
var mainGoPath string
var modulePath string
var moduleDir string

var m = flag.String("main", ".", "main package name")
var wd = flag.String("wd", "", "working directory")
var args = flag.String("args", "", "args")
var delayBuild = flag.Int("delay", 1000, "delay *ms before rebuild")
var buildFlags = flag.String("build", "", "build flags")
var verbose = flag.Bool("verbose", false, "more log")

func main() {
	flag.Parse()
	var err error
	mainPkg, err = filepath.Abs(*m)
	if err != nil {
		log.Fatal(err)
	}
	for _, goPath := range goPaths {
		if strings.HasPrefix(mainPkg, goPath) {
			mainGoPath = goPath
			break
		}
	}
	if mainGoPath == "" {
		mainGoPath = goPaths[0]
	}
	mainBin = mainGoPath + "/bin/" + path.Base(mainPkg)
	mainPkgDir = strings.TrimSpace(goList(mainPkg, `{{ .Dir }}`))
	moduleDir = strings.TrimSpace(goList(mainPkg, `{{ .Module.Dir }}`))
	modulePath = strings.TrimSpace(goList(mainPkg, `{{ .Module.Path }}`))

	autoBuild()
	eventChannel <- ""
	<-make(chan int)
}
