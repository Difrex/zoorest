# Zoorest 

Zookeeper HTTP rest API

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-generate-toc again -->
**Table of Contents**

- [Zoorest](#zoorest)
    - [Usage](#usage)
    - [API](#api)
        - [List node childrens](#list-node-childrens)
            - [Errors](#errors)
        - [Get node data](#get-node-data)
            - [Errors](#errors)
        - [Create node recursive](#create-node-recursive)
        - [Update node](#update-node)
            - [Errors](#errors)
        - [Delete node recursive](#delete-node-recursive)
            - [Errors](#errors)
    - [Build](#build)
        - [Binary](#binary)
        - [Docker build](#docker-build)
            - [Binary file](#binary-file)
            - [Docker image](#docker-image)
- [AUTHORS](#authors)
- [LICENSE](#license)

<!-- markdown-toc end -->


## Usage

```
Usage of ./zoorest:
  -listen string
    	Address to listen (default "127.0.0.1:8889")
  -path string
    	Zk root path (default "/")
  -zk string
    	Zk servers. Comma separated (default "127.0.0.1:2181")
```

## API

### List node childrens

Method: **GET**

Location: **/v1/ls**

Return JSON
```json
curl -s -XGET http://127.0.0.1:8889/v1/ls/ | jq
{
  "childrens": [
    "two",
    "three",
    "one"
  ],
  "path": "/zoorest",
  "state": "OK",
  "error": ""
}
```

#### Errors

```json
curl -s -XGET http://127.0.0.1:8889/v1/ls/does/not/exist | jq
{
  "childrens": null,
  "path": "",
  "state": "ERROR",
  "error": "zk: node does not exist"
}
```

### Get node data

Method: **GET**

Location: **/v1/get**

Return JSON
```
curl -s -XGET http://127.0.0.1:8889/v1/get/one/data | jq
{
  "path": "/zoorest/one/data",
  "state": "OK",
  "error": "",
  "data": "eyJzb21lIjogImpzb24ifQ=="
}
```
Node data stored in *data* field as base64 encoded string
```
echo eyJzb21lIjogImpzb24ifQ== | base64 -d
{"some": "json"}
```

#### Errors

```json
 curl -s -XGET http://127.0.0.1:8889/v1/get/does/not/exist | jq
{
  "path": "",
  "state": "ERROR",
  "error": "zk: node does not exist",
  "data": null
}
```

### Create node recursive

Method: **PUT**

Location: **/v1/up**

Return string with created path
```
curl -XPUT http://127.0.0.1:8889/v1/up/two/three/four -d '{"four": "json"}'
/zoorest/two/three/four
```

### Update node

Method: **POST**

Location: **/v1/up**

Return string with updated path
```
curl -XPOST http://127.0.0.1:8889/v1/up/two -d '{"two": "json"}'
/zoorest/two
```

#### Errors

```
curl -XPOST http://127.0.0.1:8889/v1/up/twa -d '{"two": "json"}'
zk: node does not exist
```

### Delete node recursive
Method: **DELETE**

Location: **/v1/rmr**

Return string with removed path
```
curl -XDELETE http://127.0.0.1:8889/v1/rmr/two
/zoorest/two
```

#### Errors
```
curl -XPOST http://127.0.0.1:8889/v1/rmr/two
Method POST not alowed
```

## Build

### Binary

Set GOPATH variable
```
export GOPATH ${HOME}/.local
```

Get source code
```
go get github.com/Difrex/zoorest/rest
```

Get dependencies
```
cd ${GOPATH}/src/github.com/Difrex/zoorest
go get -t -v ./...
```

Build staticaly linked binary
```
go build -ldflags "-linkmode external -extldflags -static"
```

Build dynamicaly linked binary
```
go build
```

Build dynamicaly linked binary with gcc
```
go build -compile gccgo
```

### Docker build

#### Binary file

Build binary file
```
git clone https://github.com/Difrex/zoorest.git
cd zoorest
./build.sh
```
Result binary file will be placed in out/ dir

#### Docker image

Build Alpine based docker image
```
git clone https://github.com/Difrex/zoorest.git
cd zoorest
./build.sh alpine
```

Image will be tagged as zoorest:latest

# AUTHORS

Denis Zheleztsov <difrex.punk@gmail.com>

# LICENSE 

GPLv3 see [LICENSE](LICENSE)
