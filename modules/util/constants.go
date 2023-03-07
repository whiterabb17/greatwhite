package util

import (
	"encoding/base64"
	"log"
	"os"
	"time"
)

const (
	IPProvider = "http://api.ipify.org"
	Soul       = "8894fafe64c043b62"
	Version    = "0.1"
	Debug      = true
)

var (
	Doze      int    = 5
	Dem       bool   = false
	Dbg       bool   = false
	ID        string = "IDKey"
	Mycellium string = "SecretAccessKey"
	StartTime        = time.Now()
)

func WriteLog(ltype string, message string) {
	f, err := os.Stat(ltype + "_Debug.Log")
	Handle(err)
	var file *os.File
	if f == nil {
		file, err = os.Create(ltype + "_Debug.Log")
		Handle(err)
	} else {
		file, err = os.OpenFile(ltype+"_Debug.Log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		Handle(err)
	}
	defer file.Close()

	_, err2 := file.WriteString(message + "\n")

	if err2 != nil {
		//	msg := tgbotapi.NewMessage(ChatID, "[<b>ERROR</b>\t Could not write <i>"+ltype+"</i> logs")
		//	msg.ParseMode = "HTML"
		//	TBApi.Send(msg)
		log.Println(err2)
	}
}

func SSt64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Tb64(text string) string {
	encodedText := base64.StdEncoding.EncodeToString([]byte(text))
	return encodedText
}
func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func FileT64(file string) string {
	bytes, err := os.ReadFile(file)
	if (err) != nil {
		return err.Error()
	}
	return ToBase64(bytes)
}

func FileF64(text string, tag string) (string, error) {
	encodedText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create("static/files/" + tag)
	Handle(err)
	defer f.Close()
	if _, err := f.Write(encodedText); err != nil {
		log.Println(err)
		return "", err
	}
	return "Successful", err
}
func SSf64(text string, tag string) (string, error) {
	encodedText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create(tag + "_Screenshot.png")
	Handle(err)
	defer f.Close()
	if _, err := f.Write(encodedText); err != nil {
		log.Println(err)
		return "", err
	}
	return "Successful", err
}
func Fb64(text string) (string, error) {
	encodedText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Println(err)
		return text, nil
	}
	return string(encodedText), err
}
