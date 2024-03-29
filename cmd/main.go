package main

import (
	"context"
	"fmt"
	"gosocialgraph/pkg/handler"
	"gosocialgraph/server"
	"os"
	"os/signal"
	"syscall"
)

var APIVersion = "nonset"
var Environment = "nonset"

func main() {
	fmt.Println("??")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer cancel()

	appContext := handler.NewAppContext()
	server.RegisterHandlers(appContext.Router, appContext)

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Gracefully shutting down...")

				err := appContext.Router.Shutdown(ctx)
				if err != nil {
					panic("could not shutdown")
				}

				return
			}
		}
	}()

	appContext.Router.Logger.Fatal(appContext.Router.Start(":3010"))
}
