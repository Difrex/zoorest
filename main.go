package main

import (
	"flag"
	"github.com/Difrex/zoorest/rest"
	"strings"
)

var (
	zk     string
	listen string
	path   string
)

// init ...
func init() {
	flag.StringVar(&zk, "zk", "127.0.0.1:2181", "Zk servers. Comma separated")
	flag.StringVar(&listen, "listen", "127.0.0.1:8889", "Address to listen")
	flag.StringVar(&path, "path", "/", "Zk root path")
	flag.Parse()
}

// main ...
func main() {
	var z rest.Zk
	hosts := strings.Split(zk, ",")
	z.Hosts = hosts
	conn, err := z.InitConnection()
	if err != nil {
		panic(err)
	}

	zoo := rest.ZooNode{
		path,
		conn,
		z,
	}

	rest.Serve(listen, zoo)
}
