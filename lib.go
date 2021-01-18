package main

import (
	"github.com/ZerQAQ/Zvpn/Zvpn"
	"github.com/ZerQAQ/Zvpn/obfus"
	"github.com/ZerQAQ/Zvpn/protocol"
	"github.com/ZerQAQ/Zvpn/proxy"
)

type ObfusFuncType func([]byte) Zvpn.Obfuscate

var ObfusMap = map[string]ObfusFuncType{
	"RC4": obfus.NewRC4,
	"AES": obfus.NewAES256,
	"nothing": obfus.NewNothing,
}

var ProxyMap = map[string]Zvpn.Proxy{
	"SOCK5": proxy.Sock5, "Sock5": proxy.Sock5, "sock5": proxy.Sock5,
}

var ProtocolMap = map[string]Zvpn.Protocol{
	"TCP": protocol.TCP, "tcp": protocol.TCP, "Tcp": protocol.TCP,
}

type Config struct {
	CliAddr string `json:"client_address"`
	SerAddr string `json:"server_address"`
	Ptc string `json:"protocol"`
	Pxy string `json:"proxy"`
	Alg string `json:"algorithm"`
	Key string	`json:"key"`
	Role string `json:"role"`
}

func formatJsonString(s []byte) []byte {
	ret := make([]byte, 0, len(s))
	printEnd := func(spaceNum int) {
		ret = append(ret, '\n')
		for i := 0; i < spaceNum; i++ {ret = append(ret, ' ')}
	}

	spaceNum := 0
	gap := 4
	inStr := false
	for _, elm := range s{
		if !inStr {
			if elm == '}' {spaceNum -= gap}
			if elm == '}' {printEnd(spaceNum)}
		}
		if elm == '"' {inStr = !inStr}

		ret = append(ret, byte(elm))

		if !inStr {
			if elm == ':' {ret = append(ret, ' ')}
			if elm == '{' {spaceNum += gap}
			if elm == '{' || elm == '}' || elm == ',' {printEnd(spaceNum)}
		}
	}
	return ret
}