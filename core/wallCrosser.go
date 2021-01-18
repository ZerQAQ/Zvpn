package core

import (
	"github.com/ZerQAQ/Zvpn"
)

type WallerCrosser interface {
	StartClient(localAddr, ServerAddr string)
	StartServer(addr string)
}

func NewWallerCrosser(protocol Zvpn.Protocol,
	proxy Zvpn.Proxy, obfuscate Zvpn.Obfuscate) WallerCrosser {
	return &WallCrosserImply{
		ProtocolWrapper{protocol, obfuscate},
		proxy, "", "",
	}
}
