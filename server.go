package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	//Errgroup package makes it easy to handle error
	//between groutines.
	//Returned value from goroutine is necessary to test code.
	"golang.org/x/sync/errgroup"
)

func runServer(ctx context.Context, listener net.Listener, mux http.Handler) error {
	//Create ctx with stop signal.
	//Server sends response that connection is closed
	//before finishing process by these signals
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	//Run http server in another groutine
	//to be able to stop it by external action
	server := http.Server{Handler: mux}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := server.Serve(listener); err != nil &&
			//ErrServerClosed is returned after calling shutdown method
			err != http.ErrServerClosed {
			log.Printf("Failed to close: %+v", err)
			return err
		}
		return nil
	})

	//Wait until ctx is canceled
	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown http server: %+v", err)
	}

	return eg.Wait()
}
