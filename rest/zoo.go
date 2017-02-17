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
	log.Print("rm: ", path)
	c, _, err := z.Conn.Children(path)
	if err != nil {
		log.Print("[zk ERROR] ", err)
	}
	log.Print("[WARNING] Trying delete ", path)
	if len(c) > 0 {
		for _, child := range c {
			childPath := strings.Join([]string{path, child}, "/")
			z.RMR(childPath)
		}
	}
	err = z.Conn.Delete(path, -1)
	if err != nil {
		log.Print("[zk ERROR] ", err)
	}
	log.Print("[WARNING] ", path, " deleted")
}

// CreateNode ...
func (z ZooNode) CreateNode(path string, content []byte) string {
	createPath := strings.Join([]string{z.Path, path}, "")
	if strings.Contains(createPath, "//") {
		createPath = strings.Replace(createPath, "//", "/", 1)
	}
	if path == "/" {
		return "ERROR: Not creating root path\n"
	}
	log.Print("Creating ", createPath)
	_, err := z.EnsureZooPath(createPath)
	if err != nil {
		return err.Error()
	}

	return z.UpdateNode(path, content)
}

// UpdateNode ...
func (z ZooNode) UpdateNode(path string, content []byte) string {
	upPath := strings.Join([]string{z.Path, path}, "")
	if strings.Contains(upPath, "//") {
		upPath = strings.Replace(upPath, "//", "/", 1)
	}
	if upPath == "/" {
		return "Not updating root path"
	}

	_, err := z.Conn.Set(upPath, content, -1)
	if err != nil {
		return err.Error()
	}

	return upPath
}

//EnsureZooPath create zookeeper path
func (z ZooNode) EnsureZooPath(path string) (string, error) {
	flag := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	s := strings.Split(path, "/")
	var p []string
	var fullnodepath string

	for i := 1; i < len(s); i++ {
		p = append(p, strings.Join([]string{"/", s[i]}, ""))
	}

	for i := 0; i < len(p); i++ {
		fullnodepath = strings.Join([]string{fullnodepath, p[i]}, "")
		exists, _, _ := z.Conn.Exists(fullnodepath)
		if !exists {
			z.Conn.Create(fullnodepath, []byte(""), flag, acl)
		}
	}

	return fullnodepath, nil
}
