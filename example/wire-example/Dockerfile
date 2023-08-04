# build go code image
# usr golang version 1.19 image
FROM golang:1.19-alpine as builder
LABEL stage=gobuilder

ENV GOPROXY https://goproxy.cn,direct

# build namespace
WORKDIR /build

COPY go.mod .
# download go package
RUN go mod download

COPY . .
# build go
# RUN go build -o /main main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /kart main.go wire_gen.go
RUN ls
# use alpine:3 image run image
FROM alpine:3
WORKDIR /build
COPY --from=builder  /build/kart.yaml /build/kart.yaml
COPY --from=builder kart /build/kart


ENTRYPOINT ["/build/kart"]
