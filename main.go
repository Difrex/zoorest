package main

import (
	"flag"
	"github.com/Difrex/zoorest/rest"
	"strings"
)

var (
	zk       string
	listen   string
	path     string
	mc       bool
	mcHosts  string
	mcPrefix string
)

// init ...
func init() {
	flag.StringVar(&zk, "zk", "127.0.0.1:2181", "Zk servers. Comma separated")
	flag.StringVar(&listen, "listen", "127.0.0.1:8889", "Address to listen")
	flag.StringVar(&path, "path", "/", "Zk root path")
	flag.BoolVar(&mc, "mc", false, "Enable memcached support")
	flag.StringVar(&mcHosts, "mchosts", "127.0.0.1:11211", "Memcached servers. Comma separated")
	flag.StringVar(&mcPrefix, "mcprefix", "zoorest", "Memcached key prefix")
	flag.Parse()
}

// main ...
func main() {
	var z rest.Zk
	hosts := getSlice(zk, ",")
	z.Hosts = hosts
	conn, err := z.InitConnection()
	if err != nil {
		panic(err)
	}

	var zoo rest.ZooNode
	zoo.Path = path
	zoo.Conn = conn
	zoo.Zoo = z

	var MC rest.MC
	MC.Hosts = getSlice(mcHosts, ",")
	MC.Prefix = mcPrefix
	MC.Enabled = mc
	MC.Client = MC.InitConnection()

	zoo.MC = MC

	rest.Serve(listen, zoo)
}

// getSlice returm slice
func getSlice(s string, delimeter string) []string {
	return strings.Split(s, ",")
}
