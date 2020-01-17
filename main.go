package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	portFlg := flag.String("p", ":53", ":<port> (default :53)")
	protocolFlg := flag.String("t", "tcp", "TCP, UDP (default TCP)")
	configFlg := flag.String("c", "cloudflare-secure", "<filename>.json (default cloudflare-secure.json)")

	if (portFlg == nil) || (protocolFlg == nil) || (configFlg == nil) {
		log.Fatal(fmt.Print("Defaults missing or not point"))
	}

	flag.Parse()

	l, err := net.Listen(*protocolFlg, *portFlg)
	if err != nil {
		log.Fatal(err)
	}

	conf := loadConfig(*configFlg + ".json")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if conf.TLS == true {
			go secureProxy(conn, conf)
		} else {
			go proxy(conn, conf)
		}
	}
}

type configuration struct {
	Protocol         string
	ConnectionString string
	TLS              bool
}

func loadConfig(config string) configuration {
	var conf configuration
	file, err := os.Open(config)
	if err != nil {
		log.Print(err)
		return conf
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Print(err)
		return conf
	}
	return conf
}

func proxy(conn net.Conn, conf configuration) {
	defer conn.Close()

	upstream, err := net.Dial(conf.Protocol, conf.ConnectionString)
	if err != nil {
		log.Print(err)
	}

	go io.Copy(upstream, conn)
	io.Copy(conn, upstream)
}

func secureProxy(conn net.Conn, conf configuration) {
	defer conn.Close()

	connTimeout, err := net.DialTimeout(conf.Protocol, conf.ConnectionString, 3*time.Second)
	if err != nil {
		log.Print(err)
		return
	}
	defer connTimeout.Close()
	upstream := tls.Client(connTimeout, &tls.Config{ServerName: "cloudflare-dns.com"})
	defer upstream.Close()
	hserr := upstream.Handshake()
	if hserr != nil {
		log.Print(hserr)
		return
	}

	go io.Copy(upstream, conn)
	io.Copy(conn, upstream)
}
