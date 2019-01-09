package service

import (
	"fmt"
	dAPP "github.com/drep-project/drepcli/app"
	"github.com/drep-project/drepcli/drepclient/component/console"
	cliTypes "github.com/drep-project/drepcli/drepclient/types"
	"github.com/drep-project/drepcli/drepclient/types/flags"
	rpcComponent "github.com/drep-project/drepcli/rpc/component"
	"gopkg.in/urfave/cli.v1"
)

var (
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TODO add config description",
	}
)

type CliService struct {
	config *cliTypes.Config
}

func (cliService *CliService)Name() string{
	return "cli"
}

func (cliService *CliService)Api() []dAPP.API{
	return []dAPP.API{}
}

func (cliService *CliService)Flags() []cli.Flag{
	return []cli.Flag{flags.JSpathFlag, flags.ExecFlag, flags.PreloadJSFlag}
}

func (cliService *CliService)Init(executeContext *dAPP.ExecuteContext) error{
	endpoint := executeContext.CliContext.Args().First()
	if len(endpoint) == 0 {
		return fmt.Errorf("You have to specify an address")
	}
	client, err := rpcComponent.Dial(endpoint)
	if err != nil {
		return fmt.Errorf("Unable to attach to remote drep: %v", err)
	}

	path := executeContext.CommonConfig.HomeDir
	cliService.config = &cliTypes.Config{}
	cliService.config.Config = console.Config{
		HomeDir: path,
		DocRoot: executeContext.CliContext.GlobalString(flags.JSpathFlag.Name),
		Client:  client,
		Preload: flags.MakeConsolePreloads(executeContext.CliContext),
	}
	return nil
}

func (cliService *CliService)Start(executeContext *dAPP.ExecuteContext) error{
	return cliService.remoteConsole(executeContext)
}

func (cliService *CliService)Stop(executeContext *dAPP.ExecuteContext) error{
	console.Stdin.Close()
	return  nil
}

// remoteConsole will connect to a remote drep instance, attaching a JavaScript
// console to it.
func  (cliService *CliService) remoteConsole(executeContext *dAPP.ExecuteContext) error {
	console, err := console.New(cliService.config.Config)
	if err != nil {
		return fmt.Errorf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	if script := executeContext.CliContext.GlobalString(flags.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}

	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}