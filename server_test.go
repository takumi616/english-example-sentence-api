package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// Run http server.
// Try sending http request.
// Stop http server by ctx.
func TestServer_Run(t *testing.T) {
	//Create http listener.
	//Randomly select port number if passes 0 to net.Listen()
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create http listener with port %v", err)
	}

	//Create test handler.
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test handler, %s!", r.URL.Path[1:])
	})

	//Create ctx with cancel func to test if http server will be stopped
	//by external action intentionally.
	ctx, cancel := context.WithCancel(context.Background())

	//Run http server in another groutine
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return runServer(ctx, listener, mux)
	})

	//Test if http server works correctly.
	input := "test"
	url := fmt.Sprintf("http://%s/%s", listener.Addr().String(), input)
	//Check randomly selected port number.
	t.Logf("Request URL: %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("Failed to get http response: %+v", err)
	}

	//Compare response body to expected value.
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("Failed to read http response body: %v", err)
	}

	want := fmt.Sprintf("Test handler, %s!", input)
	if string(got) != want {
		t.Errorf("want value is: %q, but got: %q", want, got)
	}

	//Send cancel signal to the groutine which http server is running in
	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
