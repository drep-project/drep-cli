package log

import (
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/drep-project/drepcli/common"
	"io"
	"os"
)

var DEBUG = false

var (
	ostream Handler
	glogger *GlogHandler
)

func init() {
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	ostream = StreamHandler(output, TerminalFormat(usecolor))
	glogger = NewGlogHandler(ostream)
}

func SetUp(cfg *LogConfig) error {
	if cfg.DataDir != "" {
		if !common.IsDirExists(cfg.DataDir) {
			err := os.MkdirAll(cfg.DataDir, 0777)
			if err != nil {
				return err
			}
		}

		rfh, err := SyncRotatingFileHandler(
			cfg.DataDir,
			262144,
			JSONFormatOrderedEx(false, true),
		)
		if err != nil {
			return err
		}
		glogger.SetHandler(MultiHandler(ostream, rfh))
	}
	glogger.Verbosity(Lvl(cfg.LogLevel))
	glogger.Vmodule(cfg.Vmodule)
	glogger.BacktraceAt(cfg.BacktraceAt)
	Root().SetHandler(glogger)
	return nil
}
