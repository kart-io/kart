package image

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kart-io/kart/cmd/kart/app"
)

var (
	fileName   = "Dockerfile"
	dockerfile = `# build go code image
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
# use alpine:3 image run image
FROM alpine:3
WORKDIR /build
COPY --from=builder  /build/kart.yaml /build/kart.yaml
COPY --from=builder kart /build/kart

ENTRYPOINT ["/build/kart"]`
)

func Command() *cobra.Command {
	command := app.NewCommand("image", "This is the image command", func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	})
	return command
}

func Run(_ *cobra.Command, _ []string) {
	filePwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join(filePwd, fileName)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // 关闭文件

	_, err = f.Write([]byte(dockerfile))
	cobra.CheckErr(err)
}

func isFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
