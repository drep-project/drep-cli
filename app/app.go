package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"

	"github.com/drep-project/drepcli/common"
)

var (
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TODO add config description",
	}
	// General settings
	HomeDirFlag = common.DirectoryFlag{
		Name:  "homedir",
		Usage: "Home directory for the datadir logdir and keystore",
	}
)

// DrepApp based on the cli.App, the module service operation is encapsulated.
// The purpose is to achieve the independence of each module and reduce dependencies as far as possible.
type DrepApp struct {
	Context *ExecuteContext
	*cli.App
}

// NewApp create a new app
func NewApp() *DrepApp {
	return &DrepApp{
		Context: &ExecuteContext{},
		App:     cli.NewApp(),
	}
}

// AddService add a server into context
func (mApp DrepApp) AddService(service Service) {
	mApp.Context.AddService(service)
}

//Run read the global configuration, parse the global command parameters,
// initialize the main process one by one, and execute the service before the main process starts.
func (mApp DrepApp) Run() error {
	mApp.Before = mApp.before
	mApp.Flags = append(mApp.Flags, ConfigFileFlag)
	mApp.Flags = append(mApp.Flags, mApp.Context.GetFlags()...)
	mApp.Action = mApp.action
	if err := mApp.App.Run(os.Args); err != nil {
		return err
	}
	return nil
}

// action used to init and run each services
func (mApp DrepApp) action(ctx *cli.Context) error {
	defer func() {
		for _, service := range mApp.Context.Services {
			err := service.Stop(mApp.Context)
			if err != nil {
				return
			}
		}
	}()

	for _, service := range mApp.Context.Services {
		err := service.Init(mApp.Context)
		if err != nil {
			return err
		}
	}

	for _, service := range mApp.Context.Services {
		err := service.Start(mApp.Context)
		if err != nil {
			return err
		}
	}
	return nil
}

//  read global config before main process
func (mApp DrepApp) before(ctx *cli.Context) error {
	mApp.Context.CliContext = ctx

	homeDir := ""
	if ctx.GlobalIsSet(HomeDirFlag.Name) {
		homeDir = ctx.GlobalString(HomeDirFlag.Name)
	} else {
		homeDir = common.AppDataDir(mApp.Name, false)
	}
	mApp.Context.ConfigPath = homeDir

	mApp.Context.CommonConfig = &CommonConfig{
		HomeDir: homeDir,
	}
	phaseConfig, err := loadConfigFile(ctx, homeDir)

	if err != nil {
		return err
	}
	mApp.Context.PhaseConfig = phaseConfig

	return nil
}

//	loadConfigFile sed to read configuration files
func loadConfigFile(ctx *cli.Context, configPath string) (map[string]json.RawMessage, error) {
	configFile := filepath.Join(configPath, "config.json")

	if ctx.GlobalIsSet(ConfigFileFlag.Name) {
		file := ctx.GlobalString(ConfigFileFlag.Name)
		if common.IsFileExists(file) {
			//report error when user specify
			return nil, errors.New("specify config file not exist")
		}
		configFile = file
	}

	if !common.IsFileExists(configFile) {
		//use default
		return nil, errors.New("config file not found")
	}else{

	}
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	jsonPhase := map[string]json.RawMessage{}
	err = json.Unmarshal(content, &jsonPhase)
	if err != nil {
		return nil, err
	}
	return jsonPhase, nil
}
