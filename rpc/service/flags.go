package service

import (
	rpcTypes "github.com/drep-project/drepcli/rpc/types"
	"github.com/drep-project/drepcli/common"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

var (
	// RPC settings
	HTTPEnabledFlag = cli.BoolFlag{
		Name:  "http",
		Usage: "Enable the HTTP-RPC server",
	}
	HTTPListenAddrFlag = cli.StringFlag{
		Name:  "httpaddr",
		Usage: "HTTP-RPC server listening interface",
		Value: rpcTypes.DefaultHTTPHost,
	}
	HTTPPortFlag = cli.IntFlag{
		Name:  "httpport",
		Usage: "HTTP-RPC server listening port",
		Value: rpcTypes.DefaultHTTPPort,
	}
	HTTPCORSDomainFlag = cli.StringFlag{
		Name:  "httpcorsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value: "",
	}
	HTTPVirtualHostsFlag = cli.StringFlag{
		Name:  "httpvhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value: strings.Join([]string{"localhost"}, ","),
	}
	HTTPApiFlag = cli.StringFlag{
		Name:  "httpapi",
		Usage: "API's offered over the HTTP-RPC interface",
		Value: "",
	}
	IPCDisabledFlag = cli.BoolFlag{
		Name:  "ipcdisable",
		Usage: "Disable the IPC-RPC server",
	}
	IPCPathFlag = common.DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}
	WSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable the WS-RPC server",
	}
	WSListenAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface",
		Value: rpcTypes.DefaultWSHost,
	}
	WSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port",
		Value: rpcTypes.DefaultWSPort,
	}
	WSApiFlag = cli.StringFlag{
		Name:  "wsapi",
		Usage: "API's offered over the WS-RPC interface",
		Value: "",
	}
	WSAllowedOriginsFlag = cli.StringFlag{
		Name:  "wsorigins",
		Usage: "Origins from which to accept websockets requests",
		Value: "",
	}
	// RPC settings
	RESTEnabledFlag = cli.BoolFlag{
		Name:  "rest",
		Usage: "Enable the REST-RPC server",
	}
	RESTListenAddrFlag = cli.StringFlag{
		Name:  "restaddr",
		Usage: "REST-RPC server listening interface",
		Value: rpcTypes.DefaultRestHost,
	}
	RESTPortFlag = cli.IntFlag{
		Name:  "restport",
		Usage: "REST-RPC server listening port",
		Value: rpcTypes.DefaultRestPort,
	}
)
