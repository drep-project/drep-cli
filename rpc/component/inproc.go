// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package component

import (
	"context"
	"net"

	rpcTypes "github.com/drep-project/drepcli/rpc/types"
)

// DialInProc attaches an in-process connection to the given RPC server.
func DialInProc(handler *rpcTypes.Server) *Client {
	initctx := context.Background()
	c, _ := newClient(initctx, func(context.Context) (net.Conn, error) {
		p1, p2 := net.Pipe()
		go handler.ServeCodec(rpcTypes.NewJSONCodec(p1), rpcTypes.OptionMethodInvocation|rpcTypes.OptionSubscriptions)
		return p2, nil
	})
	return c
}
