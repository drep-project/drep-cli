package types

// RPCService gives meta information about the server.
// e.g. gives information about the loaded modules.
type RPCService struct {
	server *Server
}

// Modules returns the list of RPC services with their version number
func (s *RPCService) Modules() map[string]string {
	modules := make(map[string]string)
	for name := range s.server.services {
		modules[name] = "1.0"
	}
	return modules
}
