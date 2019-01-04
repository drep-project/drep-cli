package config

import (
	"github.com/drep-project/drepcli/ethhexutil"
	"github.com/drep-project/drepcli/mycrypto"
	"github.com/drep-project/drepcli/util"
	"github.com/drep-project/drepcli/util/flags"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	ChainIdSize = 64
	defaultPort = 55555
	defaultBlockPrize = "0x1158e460913d00000"
	ClientIdentifier = "drep" // Client identifier to advertise over the network
)

var (
	RootChain ChainIdType
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TODO add config description",
	}
	nodeConfig  *NodeConfig
)

type BootNode struct {
	PubKey  *mycrypto.Point	    `json:"pubKey"`
	IP string					`json:"ip"`
	Port int					`json:"port"`
}

type NodeConfig struct {
	HomeDir string  				`json:"homeDir,omitempty"`
	Keystore string					`json:"keystore,omitempty"`
	DbPath string  					`json:"dataDir,omitempty"`
	LogDir string  					`json:"logDir,omitempty"`
	ChainId string   		   	 	`json:"chainId"`
	Blockprize *ethhexutil.Big		`json:"blockprize"`
	Boot bool						`json:"boot"`

	RpcConfig RpcConfig				`json:"rpcConfig"`
	LogConfig LogConfig				`json:"logConfig"`
}

// GetConfig
// TODO
// Temporarily as a configuration entry, removed after code refactoring
func GetConfig() *NodeConfig{
	return nodeConfig
}

func MakeConfig(ctx *cli.Context) ( *NodeConfig, error) {
	// Load defaults.
	nodeConfig = &NodeConfig{
		
		LogConfig:LogConfig{},
		RpcConfig:RpcConfig{},
	}
	
	nodeConfig.Blockprize = (*ethhexutil.Big)(ethhexutil.MustDecodeBig(defaultBlockPrize))
	//data dir setting

	//set default dir or specify
	setDataDir(ctx, nodeConfig)

	err := loadConfigFile(ctx, nodeConfig)
	if err !=  nil {
		return nil, err
	}

	// log
	setLogConfig(ctx,nodeConfig)

	//TODO
	//SetP2PConfig(ctx, &cfg.P2P)

	//rpc Config
	setRpc(ctx, nodeConfig)
	return nodeConfig, nil
}

func loadConfigFile(ctx *cli.Context, nodeConfig *NodeConfig) error {
	configFile := filepath.Join(nodeConfig.HomeDir, "config.json")

	if ctx.GlobalIsSet(ConfigFileFlag.Name) {
		file := ctx.GlobalString(ConfigFileFlag.Name)
		if util.IsFileExists(file) {
			//report error when user specify
			return errors.New("specify config file not exist")
		}
		configFile = file
	}

	if !util.IsFileExists(configFile) {
		//use default
		return nil
	}
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, nodeConfig)
}
/*
func SetP2PConfig(ctx *cli.Context, cfg *p2p.Config) {

}
*/
// setLogConfig creates an log configuration from the set command line flags,
func setLogConfig(ctx *cli.Context, cfg *NodeConfig) {
	if ctx.GlobalIsSet(flags.LogLevelFlag.Name) {
		cfg.LogConfig.LogLevel = ctx.GlobalInt(flags.LogLevelFlag.Name)
	}else{
		cfg.LogConfig.LogLevel = 3
	}

	if ctx.GlobalIsSet(flags.VmoduleFlag.Name) {
		cfg.LogConfig.Vmodule = ctx.GlobalString(flags.VmoduleFlag.Name)
	}

	if ctx.GlobalIsSet(flags.BacktraceAtFlag.Name) {
		cfg.LogConfig.BacktraceAt = ctx.GlobalString(flags.BacktraceAtFlag.Name)
	}
}


// setRpc creates an rpc configuration from the set command line flags,
func setRpc(ctx *cli.Context, cfg *NodeConfig) {
	setIPC(ctx, cfg)
	setHTTP(ctx, cfg)
	setWS(ctx, cfg)
	setRest(ctx, cfg)
}


// setIPC creates an IPC path configuration from the set command line flags,
// returning an empty string if IPC was explicitly disabled, or the set path.
func setIPC(ctx *cli.Context, cfg *NodeConfig) {
	cfg.RpcConfig.IPCEnabled = true
	if ctx.GlobalBool(flags.IPCDisabledFlag.Name) {
		cfg.RpcConfig.IPCEnabled = false
		return
	}

	checkExclusive(ctx, flags.IPCDisabledFlag, flags.IPCPathFlag)
	if ctx.GlobalIsSet(flags.IPCPathFlag.Name) {
		cfg.RpcConfig.IPCPath = ctx.GlobalString(flags.IPCPathFlag.Name)
	}else{
		cfg.RpcConfig.IPCPath = path.Join(cfg.HomeDir, DefaultIPCEndpoint(ClientIdentifier))
	}
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setHTTP(ctx *cli.Context, cfg *NodeConfig) {
	cfg.RpcConfig.HTTPEnabled = true
	if !ctx.GlobalBool(flags.HTTPEnabledFlag.Name) {
		cfg.RpcConfig.HTTPEnabled = false
		return
	}

	if ctx.GlobalIsSet(flags.HTTPListenAddrFlag.Name) {
		cfg.RpcConfig.HTTPHost = ctx.GlobalString(flags.HTTPListenAddrFlag.Name)
	} else {
		cfg.RpcConfig.HTTPHost = DefaultHTTPHost
	}

	if ctx.GlobalIsSet(flags.HTTPPortFlag.Name) {
		cfg.RpcConfig.HTTPPort = ctx.GlobalInt(flags.HTTPPortFlag.Name)
	}else{
		cfg.RpcConfig.HTTPPort = DefaultHTTPPort
	}

	if ctx.GlobalIsSet(flags.HTTPCORSDomainFlag.Name) {
		cfg.RpcConfig.HTTPCors = splitAndTrim(ctx.GlobalString(flags.HTTPCORSDomainFlag.Name))
	}

	if ctx.GlobalIsSet(flags.HTTPApiFlag.Name) {
		cfg.RpcConfig.HTTPModules = splitAndTrim(ctx.GlobalString(flags.HTTPApiFlag.Name))
	}

	if ctx.GlobalIsSet(flags.HTTPVirtualHostsFlag.Name) {
		cfg.RpcConfig.HTTPVirtualHosts = splitAndTrim(ctx.GlobalString(flags.HTTPVirtualHostsFlag.Name))
	} else {
		cfg.RpcConfig.HTTPVirtualHosts = []string{"localhost"}
	}
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setRest(ctx *cli.Context, cfg *NodeConfig) {
	cfg.RpcConfig.RESTEnabled = true
	if !ctx.GlobalBool(flags.RESTEnabledFlag.Name) {
		cfg.RpcConfig.RESTEnabled = false
		return
	}

	if ctx.GlobalIsSet(flags.RESTListenAddrFlag.Name) {
		cfg.RpcConfig.RESTHost = ctx.GlobalString(flags.RESTListenAddrFlag.Name)
	} else {
		cfg.RpcConfig.RESTHost = DefaultRestHost
	}

	if ctx.GlobalIsSet(flags.RESTPortFlag.Name) {
		cfg.RpcConfig.RESTPort = ctx.GlobalInt(flags.RESTPortFlag.Name)
	}else{
		cfg.RpcConfig.RESTPort = DefaultRestPort
	}
}

// setWS creates the WebSocket RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setWS(ctx *cli.Context, cfg *NodeConfig) {

	cfg.RpcConfig.WSEnabled = true
	if !ctx.GlobalBool(flags.WSEnabledFlag.Name) {
		cfg.RpcConfig.WSEnabled = false
		return
	}

	if ctx.GlobalIsSet(flags.WSListenAddrFlag.Name) {
		cfg.RpcConfig.WSHost = ctx.GlobalString(flags.WSListenAddrFlag.Name)
	} else{
		cfg.RpcConfig.WSHost =  DefaultWSHost
	}

	if ctx.GlobalIsSet(flags.WSPortFlag.Name) {
		cfg.RpcConfig.WSPort = ctx.GlobalInt(flags.WSPortFlag.Name)
	}else{
		cfg.RpcConfig.WSPort = DefaultWSPort
	}

	if ctx.GlobalIsSet(flags.WSAllowedOriginsFlag.Name) {
		cfg.RpcConfig.WSOrigins = splitAndTrim(ctx.GlobalString(flags.WSAllowedOriginsFlag.Name))
	}

	if ctx.GlobalIsSet(flags.WSApiFlag.Name) {
		cfg.RpcConfig.WSModules = splitAndTrim(ctx.GlobalString(flags.WSApiFlag.Name))
	}
}

func setDataDir(ctx *cli.Context, cfg *NodeConfig) {
	if ctx.GlobalIsSet(flags.HomeDirFlag.Name) {
		cfg.HomeDir = ctx.GlobalString(flags.HomeDirFlag.Name)
	} else{
		cfg.HomeDir = DefaultDataDir()
	}

	//keystore
	if ctx.GlobalIsSet(flags.KeyStoreDirFlag.Name) {
		cfg.Keystore = ctx.GlobalString(flags.KeyStoreDirFlag.Name)
	} else{
		cfg.Keystore = path.Join(cfg.HomeDir, "keystore")
	}

	//databasedir
	if ctx.GlobalIsSet(flags.DataDirFlag.Name) {
		cfg.DbPath = ctx.GlobalString(flags.DataDirFlag.Name)
	} else{
		cfg.DbPath = path.Join(cfg.HomeDir, "data")
	}

	//logdir
	if ctx.GlobalIsSet(flags.LogDirFlag.Name) {
		cfg.LogConfig.DataDir = ctx.GlobalString(flags.LogDirFlag.Name)
	} else{
		cfg.LogConfig.DataDir = path.Join(cfg.HomeDir, "log")
	}
}
// checkExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func checkExclusive(ctx *cli.Context, args ...interface{}) {
	set := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		// Make sure the next argument is a flag and skip if not set
		flag, ok := args[i].(cli.Flag)
		if !ok {
			panic(fmt.Sprintf("invalid argument, not cli.Flag type: %T", args[i]))
		}
		// Check if next arg extends current and expand its name if so
		name := flag.GetName()

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag check, make sure value set doesn't conflict with passed in option
				if ctx.GlobalString(flag.GetName()) == option {
					name += "=" + option
					set = append(set, "--"+name)
				}
				// shift arguments and continue
				i++
				continue

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.GlobalIsSet(flag.GetName()) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		flags.Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
	}
}

// splitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func splitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}


func DefaultDataDir() string {
	return util.AppDataDir("drep", false)
}

// DefaultIPCEndpoint returns the IPC path used by default.
func DefaultIPCEndpoint(clientIdentifier string) string {
	if clientIdentifier == "" {
		clientIdentifier = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if clientIdentifier == "" {
			panic("empty executable name")
		}
	}

	return clientIdentifier + ".ipc"
}

