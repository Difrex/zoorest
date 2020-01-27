package main

import (
	"flag"
	"github.com/Difrex/zoorest/rest"
	"os"
	"strings"
)

const (
	CorsFeatureEnableEnvVar string = `ZOOREST_CORS_ENABLE`
	CorsDebugModeEnvVar     string = `ZOOREST_CORS_DEBUG_ENABLE`
	CorsAllowedOrigins      string = `ZOOREST_CORS_ALLOWED_ORIGINS`
)

var (
	zk       string
	listen   string
	path     string
	mc       bool
	mcHosts  string
	mcPrefix string
	ok       bool
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

	// get CORS settins from environment
	// start by establishing defaults.
	var cors = rest.CorsOptions{
		Enabled:        false,
		DebugEnabled:   false,
		AllowedOrigins: []string{"*"},
	}

	var corsEnabled string
	var ok bool
	corsEnabled, ok = os.LookupEnv(CorsFeatureEnableEnvVar)
	if ok && corsEnabled != "" {
		cors.Enabled = true
	}

	var corsDebugEnabled string
	corsDebugEnabled, ok = os.LookupEnv(CorsDebugModeEnvVar)
	if ok && corsDebugEnabled != "" {
		cors.DebugEnabled = true
	}

	var corsAllowedOrigins string
	corsAllowedOrigins, ok = os.LookupEnv(CorsAllowedOrigins)
	if ok && corsAllowedOrigins != "" {
		cors.AllowedOrigins = strings.Split(corsAllowedOrigins, ",")
	}

	rest.Serve(listen, zoo, cors)
}

// getSlice returm slice
func getSlice(s string, delimeter string) []string {
	return strings.Split(s, ",")
}
