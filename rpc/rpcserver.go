package rpc

import (
	"fmt"
	"sync"
	"net"
	"strings"
	"github.com/drep-project/drepcli/log"
	"github.com/drep-project/drepcli/config"
)

type RpcServer struct {
	RpcAPIs       []API   // List of APIs currently provided by the node
	RestApi RestDescription
	inprocHandler *Server // In-process RPC request handler to process the API requests

	IpcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	IpcListener net.Listener // IPC RPC listener socket to serve API requests
	IpcHandler  *Server  // IPC RPC request handler to process the API requests

	HttpEndpoint  string       // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	HttpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	HttpListener  net.Listener // HTTP RPC listener socket to server API requests
	HttpHandler   *Server  // HTTP RPC request handler to process the API requests

	WsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	WsListener net.Listener // Websocket RPC listener socket to server API requests
	WsHandler  *Server  // Websocket RPC request handler to process the API requests
	
	RestEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	RestController *RestController // Websocket RPC listener socket to server API requests

	lock sync.RWMutex
	RpcConfig *config.RpcConfig
}


func NewRpcServer(apis []API, restApi RestDescription ,RpcConfig *config.RpcConfig)*RpcServer{
    return &RpcServer{
		IpcEndpoint: RpcConfig.IPCEndpoint(),
		HttpEndpoint:  RpcConfig.HTTPEndpoint(),
		WsEndpoint:  RpcConfig.WSEndpoint(),
		RestEndpoint: RpcConfig.RestEndpoint(),
		RpcConfig: RpcConfig,
		RpcAPIs: apis,
		RestApi: restApi,
    }
}

// startRPC is a helper method to start all the various RPC endpoint during node
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the node.
func (rpcserver *RpcServer) StartRPC() error {
	// All API endpoints started successfully
	//rpcserver.RpcAPIs = apis
	// Start the various API endpoints, terminating all in case of errors
	if err := rpcserver.StartInProc(rpcserver.RpcAPIs); err != nil {
		return err
	}
	if err := rpcserver.StartIPC(rpcserver.RpcAPIs); err != nil {
		rpcserver.StopInProc()
		return err
	}
	if err := rpcserver.StartHTTP(rpcserver.HttpEndpoint, rpcserver.RpcAPIs, rpcserver.RpcConfig.HTTPModules, rpcserver.RpcConfig.HTTPCors, rpcserver.RpcConfig.HTTPVirtualHosts, rpcserver.RpcConfig.HTTPTimeouts); err != nil {
		rpcserver.StopIPC()
		rpcserver.StopInProc()
		return err
	}
	if err := rpcserver.StartWS(rpcserver.WsEndpoint, rpcserver.RpcAPIs, rpcserver.RpcConfig.WSModules, rpcserver.RpcConfig.WSOrigins, rpcserver.RpcConfig.WSExposeAll); err != nil {
		rpcserver.StopHTTP()
		rpcserver.StopIPC()
		rpcserver.StopInProc()
		return err
	}

	if err := rpcserver.StartRest(rpcserver.RestEndpoint,rpcserver.RestApi); err != nil {
		rpcserver.StopREST()
		return err
	}
	return nil
}

// StartInProc initializes an in-process RPC endpoint.
func (rpcserver *RpcServer) StartInProc(apis []API) error {
	// Register all the APIs exposed by the services
	handler := NewServer()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
		log.Debug("InProc registered", "namespace", api.Namespace)
	}
	rpcserver.inprocHandler = handler
	return nil
}

// StopInProc terminates the in-process RPC endpoint.
func (rpcserver *RpcServer) StopInProc() {
	if rpcserver.inprocHandler != nil {
		rpcserver.inprocHandler.Stop()
		rpcserver.inprocHandler = nil
	}
}

// StartIPC initializes and starts the IPC RPC endpoint.
func (rpcserver *RpcServer) StartIPC(apis []API) error {
	if !rpcserver.RpcConfig.IPCEnabled {
		return nil
	}
	if rpcserver.IpcEndpoint == "" {
		return nil // IPC disabled.
	}
	listener, handler, err := StartIPCEndpoint(rpcserver.IpcEndpoint, apis)
	if err != nil {
		return err
	}
	rpcserver.IpcListener = listener
	rpcserver.IpcHandler = handler
	log.Info("IPC endpoint opened", "url", rpcserver.IpcEndpoint)
	return nil
}

// StopIPC terminates the IPC RPC endpoint.
func (rpcserver *RpcServer) StopIPC() {
	if rpcserver.IpcListener != nil {
		rpcserver.IpcListener.Close()
		rpcserver.IpcListener = nil

		log.Info("IPC endpoint closed", "endpoint", rpcserver.IpcEndpoint)
	}
	if rpcserver.IpcHandler != nil {
		rpcserver.IpcHandler.Stop()
		rpcserver.IpcHandler = nil
	}
}

// StartHTTP initializes and starts the HTTP RPC endpoint.
func (rpcserver *RpcServer) StartHTTP(endpoint string, apis []API, modules []string, cors []string, vhosts []string, timeouts config.HTTPTimeouts) error {
	if !rpcserver.RpcConfig.HTTPEnabled {
		return nil
	}
	// Short circuit if the HTTP endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, timeouts)
	if err != nil {
		return err
	}
	log.Info("HTTP endpoint opened", "url", fmt.Sprintf("http://%s", endpoint), "cors", strings.Join(cors, ","), "vhosts", strings.Join(vhosts, ","))
	// All listeners booted successfully
	rpcserver.HttpEndpoint = endpoint
	rpcserver.HttpListener = listener
	rpcserver.HttpHandler = handler

	return nil
}

// StopHTTP terminates the HTTP RPC endpoint.
func (rpcserver *RpcServer) StopHTTP() {
	if rpcserver.HttpListener != nil {
		rpcserver.HttpListener.Close()
		rpcserver.HttpListener = nil

		log.Info("HTTP endpoint closed", "url", fmt.Sprintf("http://%s", rpcserver.HttpEndpoint))
	}
	if rpcserver.HttpHandler != nil {
		rpcserver.HttpHandler.Stop()
		rpcserver.HttpHandler = nil
	}
}

// StartWS initializes and starts the websocket RPC endpoint.
func (rpcserver *RpcServer) StartWS(endpoint string, apis []API, modules []string, wsOrigins []string, exposeAll bool) error {
	if !rpcserver.RpcConfig.WSEnabled {
		return nil
	}
	// Short circuit if the WS endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := StartWSEndpoint(endpoint, apis, modules, wsOrigins, exposeAll)
	if err != nil {
		return err
	}
	log.Info("WebSocket endpoint opened", "url", fmt.Sprintf("ws://%s", listener.Addr()))
	// All listeners booted successfully
	rpcserver.WsEndpoint = endpoint
	rpcserver.WsListener = listener
	rpcserver.WsHandler = handler

	return nil
}

// StopWS terminates the websocket RPC endpoint.
func (rpcserver *RpcServer) StopWS() {
	if rpcserver.WsListener != nil {
		rpcserver.WsListener.Close()
		rpcserver.WsListener = nil

		log.Info("WebSocket endpoint closed", "url", fmt.Sprintf("ws://%s", rpcserver.WsEndpoint))
	}
	if rpcserver.WsHandler != nil {
		rpcserver.WsHandler.Stop()
		rpcserver.WsHandler = nil
	}
}

// Stop terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (rpcserver *RpcServer) Stop() error {
	rpcserver.lock.Lock()
	defer rpcserver.lock.Unlock()
	// Terminate the API, services and the p2p server.
	rpcserver.StopWS()
	rpcserver.StopHTTP()
	rpcserver.StopIPC()
	rpcserver.RpcAPIs = nil
	return nil
}


// StartHTTP initializes and starts the HTTP RPC endpoint.
func (rpcserver *RpcServer) StartRest(endpoint string,restApi RestDescription) error {
	if !rpcserver.RpcConfig.RESTEnabled {
		return nil
	}
    go func() {
        mainController := StartRest(restApi)
        rpcserver.RestEndpoint = endpoint
        rpcserver.RestController = mainController
    }()
	return nil
}

// StopHTTP terminates the HTTP RPC endpoint.
func (rpcserver *RpcServer) StopREST()  {
	if rpcserver.RestController != nil {
		rpcserver.RestController.Stop()
		rpcserver.RestController = nil
		log.Info("REST endpoint closed", "url", fmt.Sprintf("http://%s", rpcserver.HttpEndpoint))
	}
}