package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ZerQAQ/Zvpn/Zvpn"
	"github.com/ZerQAQ/Zvpn/lib"
	"github.com/ZerQAQ/Zvpn/obfus"
	"os"
)

func main() {
	genExample := flag.Bool("g", false, "generate a example config file, save as config")
	serverAddr := flag.String("s", "", "server address")
	clientAddr := flag.String("c", "", "client address")
	configFile := flag.String("f", "config", "config file path")
	end := flag.String("r", "", "role, server or client")
	flag.Parse()

	if *genExample {
		f, err := os.Create("example_config")
		//f, err := os.OpenFile("config", os.O_RDWR, 0)
		if err != nil {panic(err)}

		example, err := json.Marshal(Config{
			"0.0.0.0:1080", "0.0.0.0:1234", "TCP", "sock5", "RC4",
			hex.EncodeToString(lib.RandomBytes(obfus.RC4MaxMasterKeySize / 8)), "server or client",
		})

		if _, err = f.Write(formatJsonString(example)); err != nil {panic(err)}

		if err = f.Close(); err != nil {panic(err)}
		os.Exit(0)
	}

	var conf Config
	buf := make([]byte, 1024)
	f, err := os.Open(*configFile); if err != nil {panic(err)}
	bufLen, err := f.Read(buf); if err !=nil {panic(err)}
	if err := json.Unmarshal(buf[:bufLen], &conf); err != nil {panic(err)}

	if *serverAddr != "" {conf.SerAddr = *serverAddr}
	if *clientAddr != "" {conf.CliAddr = *clientAddr}
	if *end != "" {conf.Role = *end}

	if conf.Role != "client" && conf.Role != "server" {
		fmt.Println("end must be client or server, check you config file")
		os.Exit(0)
	}

	ptc, ok := ProtocolMap[conf.Ptc]
	if !ok {fmt.Println("protocol", conf.Ptc, "is not supported"); os.Exit(0)}
	pxy, ok := ProxyMap[conf.Pxy]
	if !ok {fmt.Println("proxy_protocol", conf.Pxy, "is not supported"); os.Exit(0)}
	obfFunc, ok := ObfusMap[conf.Alg]
	if !ok {fmt.Println("algorithm", conf.Alg, "is not supported"); os.Exit(0)}
	obfKey, err := hex.DecodeString(conf.Key); if err != nil {panic(err)}
	w := Zvpn.NewWallerCrosser(ptc, pxy, obfFunc(obfKey))

	switch conf.Role {
	case "client": w.StartClient(conf.CliAddr, conf.SerAddr)
	case "server": w.StartServer(conf.SerAddr)
	default:
		fmt.Println("role must be client or server.")
	}
}