// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// drep is the official command-line client for Ethereum.
package main

import (
	"fmt"
	"math"
	"os"
	"path"
	godebug "runtime/debug"
	"strconv"
	"strings"

	"github.com/drep-project/drepcli/config"
	"github.com/drep-project/drepcli/drepclient/console"
	"github.com/drep-project/drepcli/log"
	"github.com/drep-project/drepcli/rpc"

	"github.com/drep-project/drepcli/util/flags"
	"github.com/elastic/gosigar"
	"gopkg.in/urfave/cli.v1"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = flags.NewApp(gitCommit, "the drep command line interface")
	nCfg *config.NodeConfig
	consoleFlags = []cli.Flag{flags.JSpathFlag, flags.ExecFlag, flags.PreloadJSFlag}
)

func init() {
	// Initialize the CLI app and start Drep
	app.Action = remoteConsole
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2013-2018 The drep Authors"
	app.Flags = append(app.Flags, consoleFlags...)

	app.Before = func(ctx *cli.Context) error {
		var err error
		nCfg, err = config.MakeConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		err = log.SetUp(&nCfg.LogConfig)  //logDir config here
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		// Cap the cache allowance and tune the garbage collector
		var mem gosigar.Mem
		if err := mem.Get(); err == nil {
			allowance := int(mem.Total / 1024 / 1024 / 3)
			if cache := ctx.GlobalInt(flags.CacheFlag.Name); cache > allowance {
				log.Warn("Sanitizing cache to Go's GC limits", "provided", cache, "updated", allowance)
				ctx.GlobalSet(flags.CacheFlag.Name, strconv.Itoa(allowance))
			}
		}
		// Ensure Go's GC ignores the database cache for trigger percentage
		cache := ctx.GlobalInt(flags.CacheFlag.Name)
		gogc := math.Max(20, math.Min(100, 100/(float64(cache)/1024)))

		log.Debug("Sanitizing Go's GC trigger", "percent", int(gogc))
		godebug.SetGCPercent(int(gogc))
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// remoteConsole will connect to a remote drep instance, attaching a JavaScript
// console to it.
func remoteConsole(ctx *cli.Context) error {
	// Attach to a remotely running drep instance and start the JavaScript console
	endpoint := ctx.Args().First()
	path := config.DefaultDataDir()
	if endpoint == "" {
		if ctx.GlobalIsSet(flags.DataDirFlag.Name) {
			path = ctx.GlobalString(flags.DataDirFlag.Name)
		}
		endpoint = fmt.Sprintf("%s/drep.ipc", path)
	}
	client, err := dialRPC(nCfg, endpoint)
	if err != nil {
		flags.Fatalf("Unable to attach to remote drep: %v", err)
	}
	config := console.Config{
		HomeDir: path,
		DocRoot: ctx.GlobalString(flags.JSpathFlag.Name),
		Client:  client,
		Preload: flags.MakeConsolePreloads(ctx),
	}

	console, err := console.New(config)
	if err != nil {
		flags.Fatalf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	if script := ctx.GlobalString(flags.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}

	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}

// dialRPC returns a RPC client which connects to the given endpoint.
// The check for empty endpoint implements the defaulting logic
// for "drep attach" and "drep monitor" with no argument.
func dialRPC(cfg *config.NodeConfig, endpoint string) (*rpc.Client, error) {
	if endpoint == "" {
		endpoint = path.Join(cfg.HomeDir, config.DefaultIPCEndpoint(config.ClientIdentifier))
	} else if strings.HasPrefix(endpoint, "rpc:") || strings.HasPrefix(endpoint, "ipc:") {
		// Backwards compatibility with drep < 1.5 which required
		// these prefixes.
		endpoint = endpoint[4:]
	}
	return rpc.Dial(endpoint)
}