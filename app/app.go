package app

import (
	"encoding/json"
	"github.com/drep-project/drepcli/common"
	"io/ioutil"
	"os"
	"errors"
	"path/filepath"
	"gopkg.in/urfave/cli.v1"
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

type DrepApp struct {
	Context *ExecuteContext
	*cli.App
}

func NewApp () *DrepApp{
	return &DrepApp{
		Context: &ExecuteContext{},
		App:cli.NewApp(),
	}
}
func (mApp DrepApp) AddService(service Service)  {
	mApp.Context.AddService(service)
}

func (mApp DrepApp) Run() error {
	mApp.Before = mApp.before
	mApp.Flags =append(mApp.Flags, ConfigFileFlag)
	mApp.Flags =append(mApp.Flags, mApp.Context.GetFlags()...)
	mApp.Action = mApp.action
	if err := mApp.App.Run(os.Args); err != nil {
		return err
	}
	return nil
}
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

func (mApp DrepApp) before(ctx *cli.Context) error {
	mApp.Context.CliContext = ctx

	homeDir := ""
	if ctx.GlobalIsSet(HomeDirFlag.Name) {
		homeDir = ctx.GlobalString(HomeDirFlag.Name)
	} else{
		homeDir = common.AppDataDir(mApp.Name, false)
	}
	mApp.Context.ConfigPath = homeDir

	mApp.Context.CommonConfig = &CommonConfig{
		HomeDir:homeDir,
	}
	phaseConfig, err := loadConfigFile(ctx,homeDir)

	if err != nil {
		return err
	}
	mApp.Context.PhaseConfig = phaseConfig

	return nil
}

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
		return nil, errors.New("file not exits")
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