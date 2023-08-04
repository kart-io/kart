package image

import (
	"fmt"
	"github.com/kart-io/kart/cmd/kart/app"
	"github.com/spf13/cobra"
	"os"
)

var dockerfile = `# build go code image
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


ENTRYPOINT ["/build/kart"]`

func Command() *cobra.Command {
	command := app.NewCommand("image", "This is the image command", func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	})
	return command
}

func Run(_ *cobra.Command, _ []string) {
	fileName := "Dockerfile"
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(dockerfile))
		cobra.CheckErr(err)
	}
}
