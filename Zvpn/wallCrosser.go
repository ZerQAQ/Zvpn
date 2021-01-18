package Zvpn

type WallerCrosser interface {
	StartClient(localAddr, ServerAddr string)
	StartServer(addr string)
}

func NewWallerCrosser(protocol Protocol,
	proxy Proxy, obfuscate Obfuscate) WallerCrosser {
	return &WallCrosserImply{
		ProtocolWrapper{protocol, obfuscate},
		proxy, "", "",
	}
}
