package Zvpn

type WallCrosser interface {
	StartClient(localAddr, ServerAddr string)
	StartServer(addr string)
}

func NewWallCrosser(protocol Protocol,
	proxy Proxy, obfuscate Obfuscate) WallCrosser {
	return &WallCrosserImply{
		ProtocolWrapper{protocol, obfuscate},
		proxy, "", "",
	}
}
