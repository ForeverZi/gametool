package main

import (
	"fmt"
	"flag"
	"path/filepath"
	"log"
	"strings"
	"regexp"
	"encoding/json"
	"io"
	"os"
	"bytes"
)

var searchPath string
var targetPath string

func init(){
	flag.StringVar(&searchPath, "s", ".", "please enter root directory of stage resource")
	flag.StringVar(&targetPath, "d", ".", "please enter target directory of stage_parts.json")
}

func main(){
	flag.Parse()
	basepath, err := filepath.Abs(searchPath)
	if err != nil {
		log.Fatal(err)
	}
	finalBase, err := filepath.Abs(targetPath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(finalBase, 0777)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("search path: ", basepath)
	sep := string([]rune{filepath.Separator})
	if !strings.HasSuffix(basepath, sep) {
		basepath += sep
	}
	stagedirs, err := filepath.Glob(basepath + "stage_[1-9][0-9][0-9][0-9]")
	if err != nil {
		log.Fatal(err)
	}
	stageConf := make(map[string][]string)
	for _, dirpath := range stagedirs {
		fmt.Println("stagedir:", dirpath)
		stageIdReg := regexp.MustCompile(`\d{4}`)
		stageId := stageIdReg.FindString(filepath.Base(dirpath))
		if stageId == "" {
			log.Fatal(dirpath, " without stageId")
		}
		files, err := filepath.Glob(dirpath+sep+"*")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if strings.HasSuffix(file, ".prefab") {
				stageConf[stageId] = append(stageConf[stageId], filepath.Base(file[0:len(file)-7]))
			}
		}
	}
	stageConfBytes, _ := json.MarshalIndent(stageConf, "", "  ")
	wf, err := os.OpenFile(filepath.Join(finalBase, "stage_parts.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer wf.Close()
	mw := io.MultiWriter(os.Stdout, wf)
	_, err = io.Copy(mw, bytes.NewReader(stageConfBytes))
	if err != nil {
		log.Fatal(err)
	}
}