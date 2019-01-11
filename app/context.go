package app

import (
	"encoding/json"
	"gopkg.in/urfave/cli.v1"
)

var (
	CommandHelpTemplate = `{{.cmd.Name}}{{if .cmd.Subcommands}} command{{end}}{{if .cmd.Flags}} [command options]{{end}} [arguments...]
{{if .cmd.Description}}{{.cmd.Description}}
{{end}}{{if .cmd.Subcommands}}
SUBCOMMANDS:
	{{range .cmd.Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	{{end}}{{end}}{{if .categorizedFlags}}
{{range $idx, $categorized := .categorizedFlags}}{{$categorized.Name}} OPTIONS:
{{range $categorized.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}{{end}}`
)

func init() {
	cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = CommandHelpTemplate
}

type CommonConfig struct {
	HomeDir string `json:"homeDir,omitempty"`
}

// API describes the set of methods offered over the RPC interface
type API struct {
	Namespace string      // namespace under which the rpc methods of Service are exposed
	Version   string      // api version for DApp's
	Service   interface{} // receiver instance which holds the methods
	Public    bool        // indication if the methods must be considered safe for public use
}

type Service interface {
	Name() string
	Api() []API
	Flags() []cli.Flag

	Init(executeContext *ExecuteContext) error
	Start(executeContext *ExecuteContext) error
	Stop(executeContext *ExecuteContext) error
}

type ExecuteContext struct {
	ConfigPath   string
	CommonConfig *CommonConfig //
	PhaseConfig  map[string]json.RawMessage
	CliContext   *cli.Context

	Services []Service

	GitCommit string
	Usage     string
}

func (econtext *ExecuteContext) AddService(service Service) {
	econtext.Services = append(econtext.Services, service)
}

func (econtext *ExecuteContext) GetService(name string) Service {
	for _, service := range econtext.Services {
		if service.Name() == name {
			return service
		}
	}
	return nil
}

func (econtext *ExecuteContext) GetConfig(phaseName string) json.RawMessage {
	phaseConfig, ok := econtext.PhaseConfig[phaseName]
	if ok {
		return phaseConfig
	} else {
		return nil
	}
}

func (econtext *ExecuteContext) GetFlags() []cli.Flag {
	flags := []cli.Flag{}
	for _, service := range econtext.Services {
		flags = append(flags, service.Flags()...)
	}
	return flags
}

func (econtext *ExecuteContext) GetApis() []API {
	apis := []API{}
	for _, service := range econtext.Services {
		apis = append(apis, service.Api()...)
	}
	return apis
}
