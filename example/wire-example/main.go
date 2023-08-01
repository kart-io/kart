package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	apiServer, cleanup, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	err = apiServer.Run()
	if err != nil {
		return
	}
	//
	//chSig := make(chan os.Signal, 1)
	//signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//if err = app.Run(); err != nil {
	//	chSig <- syscall.SIGTERM
	//}
	//logger.Infow("Server exiting")
	//<-chSig
}
