package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

const (
	LocalMysqlDSN         = "reg:pass@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	InsertRackQueryPrefix = "INSERT INTO racks (name, ipaddr, live) VALUES (?, ?, ?)"
)

var (
	startInd    = 1
	endInd      = 200
	namePrefix  = "my_rack_go2_"
	fruitNames  = []string{"papaya", "apple", "banana", "watermelon", "raspberry", "blueberry", "strawberry", "lemon", "avocado", "grape"}
	ipDomain    = "10.20.0."
	truthValues = []bool{true, false}
	serverStatusList = []string{"Live", "Maintenance", "Retired", "Reserve", "ProvisioningOS", "HWClassified", "Prep", "Deallocation",
		"Prototype", "Migration", "Defective", "Hibernation"}
	hostnamePrefix = "my_server_"
)

type Rack struct {
	Id     int    `sql:"AUTO_INCREMENT"`
	Name   string `sql:"varchar(20)"`
	Ipaddr string `sql:"varchar(20)"`
	Live   bool   `sql:"boolean"`
}

func InsertRackData() {
	db, err := gorm.Open(mysql.Open(LocalMysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't open handle to MySQL")
	}

	rand.Seed(time.Now().UnixNano())
	for i := startInd; i <= endInd; i++ {
		var rack Rack
		rack.Name = fmt.Sprintf("%s%d", namePrefix, i)
		rack.Ipaddr = fmt.Sprintf("%s%d", ipDomain, i)
		rack.Live = truthValues[rand.Intn(2)]
		db.Create(&rack)
	}
}

func InsertFruitData() {
	db, err := gorm.Open(mysql.Open(LocalMysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't open handle to MySQL")
	}

	rand.Seed(time.Now().UnixNano())
	for i := startInd; i <= endInd; i++ {
		var f entity.Fruit
		f.Name = fmt.Sprintf("%s%d", fruitNames[rand.Intn(len(fruitNames))], rand.Intn(10))
		f.Quantity = rand.Intn(2000)
		db.Create(&f)
	}
}

func InsertServerStatusData() {
	db, err := gorm.Open(mysql.Open(LocalMysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't open handle to MySQL")
	}

	rand.Seed(time.Now().UnixNano())
	for i, ss := range serverStatusList {
		var f entity.ServerStatus
		f.Id =i+1	// 0 indexing will be an issue in DBs
		f.Name = ss
		db.Table("ServerStatus").Create(&f)
	}
}

func InsertServerData() {
	db, err := gorm.Open(mysql.Open(LocalMysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't open handle to MySQL")
	}

	rand.Seed(time.Now().UnixNano())
	for i:= startInd; i <= endInd; i++ {
		var f entity.Server
		f.Id =i+1	// 0 indexing will be an issue in DBs
		f.HostName = fmt.Sprintf("%s%d", hostnamePrefix, i+1)
		f.PublicIpAddress = fmt.Sprintf("%s%d", ipDomain, i)
		f.ServerStatus = rand.Intn(len(serverStatusList)+1)
		db.Create(&f)
	}
}

func main() {
	//InsertRackData()
	//InsertFruitData()
	//InsertServerStatusData()
	InsertServerData()
}
