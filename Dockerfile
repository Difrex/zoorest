FROM alpine

MAINTAINER Denis Zheleztsov <difrex.punk@gmail.com>

ADD zoorest /usr/bin/

EXPOSE 8889

ENTRYPOINT ["/usr/bin/zoorest"]
