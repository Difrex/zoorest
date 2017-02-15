package rest

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
	"time"
)

//ZooNode zookeeper node
type ZooNode struct {
	Path string
	Conn *zk.Conn
	Zoo  Zk
}

//Zk Zookeeper connection settings
type Zk struct {
	Hosts []string
	Path  string
	Tick  int
}

//InitConnection Initialize Zookeeper connection
func (conf Zk) InitConnection() (*zk.Conn, error) {
	conn, _, err := zk.Connect(conf.Hosts, time.Second)
	if err != nil {
		log.Panic("[ERROR] ", err)
	}

	return conn, err
}

// GetChildrens ...
func (z ZooNode) GetChildrens(path string) Ls {
	var lsPath string
	lsPath = strings.Join([]string{z.Path, path}, "")
	if path == "/" {
		lsPath = z.Path
	}

	if strings.Contains(lsPath, "//") {
		lsPath = strings.Replace(lsPath, "//", "/", 1)
	}

	log.Print("ls: ", lsPath)

	var l Ls
	l.State = "OK"

	childrens, _, err := z.Conn.Children(lsPath)
	if err != nil {
		l.State = "ERROR"
		l.Error = err.Error()
		return l
	}

	l.Error = ""
	l.Childrens = childrens
	l.Path = lsPath

	return l
}

// GetNode ...
func (z ZooNode) GetNode(path string) Get {
	var getPath string
	getPath = strings.Join([]string{z.Path, path}, "")
	if path == "/" {
		getPath = z.Path
	}

	if strings.Contains(getPath, "//") {
		getPath = strings.Replace(getPath, "//", "/", 1)
	}

	log.Print("get: ", getPath)

	var g Get
	g.State = "OK"

	data, _, err := z.Conn.Get(getPath)
	if err != nil {
		g.State = "ERROR"
		g.Error = err.Error()
		return g
	}

	g.Error = ""
	g.Data = data
	g.Path = getPath

	return g
}

//RMR remove Zk node recursive
func (z ZooNode) RMR(path string) {
	var rmPath string
	rmPath = strings.Join([]string{z.Path, path}, "")
	if path == "/" {
		return
	}

	if strings.Contains(rmPath, "//") {
		rmPath = strings.Replace(rmPath, "//", "/", 1)
	}

	log.Print("rm: ", rmPath)
	c, _, err := z.Conn.Children(rmPath)
	if err != nil {
		log.Print("[zk ERROR] ", err)
	}
	log.Print("[WARNING] Trying delete ", path)
	if len(c) > 0 {
		for _, child := range c {
			childPath := strings.Join([]string{rmPath, child}, "/")
			z.RMR(childPath)
		}
	}
	err = z.Conn.Delete(path, -1)
	if err != nil {
		log.Print("[zk ERROR] ", err)
	}
	log.Print("[WARNING] ", path, " deleted")
}
