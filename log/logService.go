package log

import (
	"encoding/json"
	"github.com/drep-project/drepcli/app"
	"gopkg.in/urfave/cli.v1"
	"path"
)

type LogServiice struct {
	config *LogConfig
}

func (logService *LogServiice) Name() string {
	return "log"
}
func (logService *LogServiice) Api() []app.API {
	return []app.API{}
}
func (logService *LogServiice) Flags() []cli.Flag {
	return []cli.Flag{LogDirFlag, LogLevelFlag, VmoduleFlag, BacktraceAtFlag}
}

func (logService *LogServiice) Init(executeContext *app.ExecuteContext) error {
	phase := executeContext.GetConfig(logService.Name())
	logService.config = &LogConfig{}
	err := json.Unmarshal(phase, logService.config)
	if err != nil {
		return err
	}
	logService.setLogConfig(executeContext.CliContext, executeContext.CommonConfig.HomeDir)
	return SetUp(logService.config)
}

func (logService *LogServiice) Start(executeContext *app.ExecuteContext) error {
	return nil
}

func (logService *LogServiice) Stop(executeContext *app.ExecuteContext) error {
	return nil
}

// setLogConfig creates an log configuration from the set command line flags,
func (logService *LogServiice) setLogConfig(ctx *cli.Context, homeDir string) {
	logService.config = &LogConfig{}
	if ctx.GlobalIsSet(LogLevelFlag.Name) {
		logService.config.LogLevel = ctx.GlobalInt(LogLevelFlag.Name)
	} else {
		logService.config.LogLevel = 3
	}

	if ctx.GlobalIsSet(VmoduleFlag.Name) {
		logService.config.Vmodule = ctx.GlobalString(VmoduleFlag.Name)
	}

	if ctx.GlobalIsSet(BacktraceAtFlag.Name) {
		logService.config.BacktraceAt = ctx.GlobalString(BacktraceAtFlag.Name)
	}

	//logdir
	if ctx.GlobalIsSet(LogDirFlag.Name) {
		logService.config.DataDir = ctx.GlobalString(LogDirFlag.Name)
	} else {
		logService.config.DataDir = path.Join(homeDir, "log")
	}
}
