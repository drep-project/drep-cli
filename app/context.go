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

// CommonConfig read before app run,this fuction shared by other moudles
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

// Services can customize their own configuration, command parameters, interfaces, services
type Service interface {
	Name() string      // service  name must be unique
	Api() []API        // Interfaces required for services
	Flags() []cli.Flag // flags required for services

	Init(executeContext *ExecuteContext) error
	Start(executeContext *ExecuteContext) error
	Stop(executeContext *ExecuteContext) error
}

// ExecuteContext centralizes all the data and global parameters of application execution,
// and each service can read the part it needs.
type ExecuteContext struct {
	ConfigPath   string
	CommonConfig *CommonConfig //
	PhaseConfig  map[string]json.RawMessage
	CliContext   *cli.Context

	Services []Service

	GitCommit string
	Usage     string
}

// AddService add a service to context, The application then initializes and starts the service.
func (econtext *ExecuteContext) AddService(service Service) {
	econtext.Services = append(econtext.Services, service)
}

// GetService In addition, there is a dependency relationship between services.
// This method is used to find the dependency services you need in the context.
func (econtext *ExecuteContext) GetService(name string) Service {
	for _, service := range econtext.Services {
		if service.Name() == name {
			return service
		}
	}
	return nil
}

//	GetConfig Configuration is divided into several segments,
//	each service only needs to obtain its own configuration data,
//	and the parsing process is also controlled by each service itself.
func (econtext *ExecuteContext) GetConfig(phaseName string) json.RawMessage {
	phaseConfig, ok := econtext.PhaseConfig[phaseName]
	if ok {
		return phaseConfig
	} else {
		return nil
	}
}

// GetFlags aggregate command configuration items required for each service
func (econtext *ExecuteContext) GetFlags() []cli.Flag {
	flags := []cli.Flag{}
	for _, service := range econtext.Services {
		flags = append(flags, service.Flags()...)
	}
	return flags
}

//	GetApis aggregate interface functions for each service to provide for use by RPC services
func (econtext *ExecuteContext) GetApis() []API {
	apis := []API{}
	for _, service := range econtext.Services {
		apis = append(apis, service.Api()...)
	}
	return apis
}
