package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"hl/app"
	"hl/config"
	"hl/load"
	"hl/model"
	"log"
	"os/exec"
	_ "os/exec"
)

func main() {
	var db *gorm.DB
	var err error

	_, err = exec.Command("sh", "-c", "unzip /tmp/data/data.zip -d /root/go/src/hl/load/data").Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	config := config.GetConfig("/root/go/src/hl/config/Config2.json")

	app := &app.App{}

	db, err = gorm.Open("sqlite3", "/root/go/src/hl/load/trav.db")
	if err != nil {
		log.Fatal("Could not connect database")
	}
	fmt.Println("WORKING!!!")
	db.Exec("PRAGMA SYNCHRONOUS=OFF;")
	db.Exec("PRAGMA JOURNAL_MODE=OFF;")
	//db.LogMode(true)
	db.Exec("PRAGMA PAGE_SIZE = 65536;")
	db.Exec("PRAGMA default_cache_size=700000;")
	db.Exec("PRAGMA cache_size=700000;")
	db = model.DBMigrate(db)

	//runtime.GOMAXPROCS(runtime.NumCPU())

	load.LoadUser(db, config.DataPath)
	load.LoadLocation(db, config.DataPath)
	load.LoadVisit(db, config.DataPath)

	db.Close()

	app.Initialize(config)

	app.Run(":80")
}
