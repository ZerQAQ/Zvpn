package core

import (
	"github.com/ZerQAQ/Zvpn/obfus"
	"github.com/ZerQAQ/Zvpn/protocol"
	"github.com/ZerQAQ/Zvpn/proxy"
)

type WallerCrosser interface {
	StartClient(localAddr, ServerAddr string)
	StartServer(addr string)
}

func NewWallerCrosser(protocol protocol.Protocol,
	proxy proxy.Proxy, obfuscate obfus.Obfuscate) WallerCrosser {
	return &WallCrosserImply{
		ProtocolWrapper{protocol, obfuscate},
		proxy, "", "",
	}
}
