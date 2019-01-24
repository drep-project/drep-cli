package main

import (
	"fmt"
	"github.com/drep-project/drepcli/log"
	"reflect"

	"github.com/drep-project/drepcli/app"
	rpcService "github.com/drep-project/drepcli/rpc/service"
	accountService "github.com/drep-project/drepcli/accounts/service"
	cliService "github.com/drep-project/drepcli/drepclient/service"
)

func main() {
	drepApp := app.NewApp()
	err := drepApp.AddServiceType(
		reflect.TypeOf(log.LogService{}),
		reflect.TypeOf(accountService.AccountService{}),
		reflect.TypeOf(rpcService.RpcService{}),
		reflect.TypeOf(cliService.CliService{}),
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	drepApp.Name = "drep"
	drepApp.Author = "Drep-Project"
	drepApp.Email = ""
	drepApp.Version = "0.1"
	drepApp.HideVersion = true
	drepApp.Copyright = "Copyright 2018 - now The drep Authors"

	if err := drepApp.Run(); err != nil {
		fmt.Println(err.Error())
	}
	return
}
