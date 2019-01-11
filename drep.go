package main

import (
	"fmt"

	"github.com/drep-project/drepcli/app"
	cliService "github.com/drep-project/drepcli/drepclient/service"
	"github.com/drep-project/drepcli/log"
)

func main() {
	drepApp := app.NewApp()
	drepApp.AddService(&log.LogServiice{})
	//drepApp.AddService(&rpcService.RpcService{})
	drepApp.AddService(&cliService.CliService{})

	drepApp.Name = "drep"
	drepApp.Author = ""
	//app.Authors = nil
	drepApp.Email = ""
	drepApp.Version = "1.0"
	drepApp.HideVersion = true
	drepApp.Copyright = "Copyright 2013-2018 The drep Authors"

	if err := drepApp.Run(); err != nil {
		fmt.Println(err.Error())
	}
	return
}
