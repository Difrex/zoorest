#!/bin/bash

case $1 in docker)
               docker build -t zoorest -f Dockerfile .
               ;;
            binary)
                go build -o zoorest .
                ;;
            docker-binary)
                docker build --rm -t zoorest -f Dockerfile .
                docker run --rm -u root -v $(pwd)/out:/tmp/out --entrypoint '/bin/cp' zoorest /bin/zoorest /tmp/out
                ;;
            *)
                ;;
esac
