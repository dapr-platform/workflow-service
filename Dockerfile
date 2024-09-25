FROM golang:1.23.1-alpine3.19 as build
RUN  apk add --no-cache git upx \
    && rm -rf /var/cache/apk/* \
    && rm -rf /root/.cache \
    && rm -rf /tmp/*
RUN mkdir /app
WORKDIR /app
#ENV GOPATH=/app/vendor:$GOPATH
COPY go.mod .
COPY go.sum .
#RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2
ENV GOPROXY=https://proxy.golang.com.cn,direct
ENV GOSUMDB=off
RUN go mod tidy
COPY . .
#RUN go clean -modcache && rm -f go.sum && go mod tidy
#RUN  rm -f go.sum && go mod tidy
#RUN swag init --parseDependency --parseInternal --parseDepth 1
RUN go build -ldflags "-s -w" -o  workflow-service && upx -9 workflow-service

FROM alpine:3.19
RUN  apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf /var/cache/apk/* \
    && rm -rf /root/.cache \
    && rm -rf /tmp/*

RUN mkdir /app
WORKDIR /app
COPY --from=build /app/workflow-service .

EXPOSE 80
CMD ["sh","-c","./workflow-service"]
