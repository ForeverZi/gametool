package main

import (
	"flag"
	"strings"
	"path/filepath"
	"fmt"
	"log"
	"sync"
	"encoding/json"
	"os"
	"io/ioutil"
)

type EnvSet []string

func (set *EnvSet)String()string{
	return strings.Join(*set, ",")
}

func (set *EnvSet)Set(str string)error{
	*set = strings.Split(str, ",")
	return nil
}

var (
	srcDir string
	destDir string
	envSet EnvSet = []string{"production", "development"}
)

var wg sync.WaitGroup

func init(){
	flag.StringVar(&srcDir, "s", ".", "please input src directory")
	flag.StringVar(&destDir, "d", "dest", "please input dest directory")
	flag.Var(&envSet, "envset", "please input envset for pomelo, default like '-envset production,devlopment'")
}

func main(){
	flag.Parse()
	fmt.Println(envSet)
	absSrcDir, err := filepath.Abs(srcDir)
	if err != nil {
		log.Fatal(err)
	}
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(absDestDir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	rawFiles, err := filepath.Glob(filepath.Join(absSrcDir, "*.json"))
	if err != nil {
		log.Fatal(err)
	}
	for _, rawFilePath := range rawFiles {
		wg.Add(1)
		go cvtFile(rawFilePath, filepath.Join(absDestDir, filepath.Base(rawFilePath)))
	}
	wg.Wait()
}

func cvtFile(src string, dest string){
	defer func(){
		if err := recover(); err == nil {
			fmt.Println("cvt complete:", dest)
		}else{
			fmt.Println("cvt error:", err)
		}
		wg.Done()
	}()
	fmt.Println("cvt begin:", src)
	data, err := readRawData(src)
	if err != nil {
		panic(err)
	}
	dataWithEnv := cvt2envFormat(data)
	err = writeEnvFormat(dest, dataWithEnv)
	if err != nil {
		panic(err)
	}
}

func writeEnvFormat(path string, data map[string]interface{})error{
	bdata, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bdata, 0666)
}

func readRawData(path string)(raw map[string]interface{}, err error){
	bdata, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(bdata, &data)
	if err == nil {
		// 一些转换工具末行会出现null
		delete(data, "null")
	}
	return data, err 
}

func cvt2envFormat(raw map[string]interface{})map[string]interface{}{
	data := make(map[string]interface{})
	for _, envKey := range envSet {
		data[envKey] = raw
	}
	return data
}

