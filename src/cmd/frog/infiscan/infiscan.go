package main

import (
	"flag"
	"log"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"strings"
	"regexp"
	"encoding/json"
	"os"
)

var searchPath string
var targetPath string

func init(){
	flag.StringVar(&searchPath, "s", ".", "please enter root directory of stage resource")
	flag.StringVar(&targetPath, "d", ".", "please enter target directory of stage_parts.json")
}

func main(){
	flag.Parse()
	absSrcDir, err := filepath.Abs(searchPath)
	if err != nil {
		log.Fatal(err)
	}
	absDestDir, err := filepath.Abs(targetPath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(absDestDir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	rawFiles, err := filepath.Glob(filepath.Join(absSrcDir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	partSeqRe := regexp.MustCompile(`\d+$`)
	infiConf := make(map[string][]string, len(rawFiles)/2)
	for _, rawFilePath := range rawFiles {
		if !strings.HasSuffix(rawFilePath, ".meta") {
			partseqNum := partSeqRe.FindString(rawFilePath)
			if partseqNum == "" {
				log.Fatal("unknow part seq with ", rawFilePath)
			}
			files, err := filepath.Glob(filepath.Join(rawFilePath, "*.prefab"))
			if err != nil {
				log.Fatal(err)
			}
			infiConf[partseqNum] = make([]string, 0, len(files))
			for _, file := range files {
				fmt.Println(filepath.Base(file[0:len(file)-7]))
				infiConf[partseqNum] = append(infiConf[partseqNum], filepath.Base(file[0:len(file)-7]))
			}
		}
	}
	data, err := json.MarshalIndent(infiConf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(absDestDir, "infinite.json"), data, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

