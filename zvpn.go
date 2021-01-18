package Zvpn

import "github.com/ZerQAQ/Zvpn/core"

func NewWallerCrosser(protocol Protocol,
	proxy Proxy, obfuscate Obfuscate) core.WallerCrosser {
	return core.NewWallerCrosser(protocol, proxy, obfuscate)
}