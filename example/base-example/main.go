// https://colobu.com/2015/10/09/Linux-Signals/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	handler := gin.Default()
	app, cleanup, err := wireApp(handler)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err = app.Run(); err != nil {
		chSig <- syscall.SIGTERM
	}
	fmt.Println("Server exiting")
	<-chSig
}
