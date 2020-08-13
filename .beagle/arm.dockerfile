FROM registry.cn-qingdao.aliyuncs.com/wod/golang-arm64:1.14.6-alpine as golang

WORKDIR /go/src/github.com/kiali/kiali

Add . ./

ENV GOPROXY=https://goproxy.cn,direct  

RUN go build -o /go/bin/kiali -ldflags "-X main.version=${VERSION}"

FROM registry.cn-qingdao.aliyuncs.com/wod/alpine-arm64:3.12

ENV KIALI_HOME=/opt/kiali \
    PATH=$KIALI_HOME:$PATH

WORKDIR $KIALI_HOME

COPY --from=golang /go/bin/kiali $KIALI_HOME/

ADD console $KIALI_HOME/console/

ENTRYPOINT ["/opt/kiali/kiali"]
