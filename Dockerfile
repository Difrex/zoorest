FROM alpine

MAINTAINER Denis Zheleztsov <difrex.punk@gmail.com>

ADD out/zoorest /bin/

RUN echo -ne "zoorest\nzoorest\n" | adduser zoorest

USER zoorest

EXPOSE 8889

ENTRYPOINT ["/bin/zoorest"]
