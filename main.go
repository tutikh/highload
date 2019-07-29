package main

import (
	"fmt"
	"hl/app"
	"hl/config"
	"hl/load"
	"os/exec"
	_ "os/exec"
)

func main() {

	_, err := exec.Command("sh", "-c", "unzip /tmp/data/data.zip -d /root/go/src/hl/load/data").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	config := config.GetConfig("/root/go/src/hl/config/Config2.json")

	app := &app.App{}
	app.Initialize(config)

	fmt.Println("WORKING!!!")

	//runtime.GOMAXPROCS(runtime.NumCPU())
	//
	load.LoadUser(app.DB, config.DataPath)
	load.LoadLocation(app.DB, config.DataPath)
	load.LoadVisit(app.DB, config.DataPath)

	app.Run(":80")
}
