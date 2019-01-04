
package config

type LogConfig struct {
	DataDir string 		`json:"-"`
	LogLevel int 		`json:"logLevel"`
	Vmodule string		`json:"vmodule,omitempty"`
	BacktraceAt string	`json:"backtraceAt,omitempty"`
}