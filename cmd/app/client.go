package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/exp/slog"

	"github.com/muesli/cancelreader"
	server_grpc "github.com/mwasilew2/go-service-template/gen/server-grpc"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clientCmd struct {
	// cli options
	HttpAddr string `help:"address of the http server which the client should connect to" default:":8080"`
	GrpcAddr string `help:"address of the grpc server which the client should connect to" default:":8081"`

	// Dependencies
	logger *slog.Logger
}

func (c *clientCmd) Run(cmdCtx *cmdContext) error {
	c.logger = cmdCtx.Logger.With("component", "clientCmd")

	// create a run group
	g := run.Group{}

	// listen for termination signals
	osSigChan := make(chan os.Signal)
	signal.Notify(osSigChan, os.Kill, os.Interrupt)
	done := make(chan struct{})
	g.Add(func() error {
		select {
		case sig := <-osSigChan:
			c.logger.Debug("caught signal", "signal", sig.String())
			return fmt.Errorf("received signal: %s", sig.String())
		case <-done:
			c.logger.Debug("closing signal catching goroutine")
		}
		return nil
	}, func(err error) {
		close(done)
	})

	// grpc error printing goroutine
	errChan := make(chan error)
	g.Add(func() error {
		c.logger.Debug("started error printing goroutine")
		for e := range errChan {
			c.logger.Error("error", "err", e)
		}
		return nil
	}, func(err error) {
		c.logger.Debug("closing error printing goroutine")
		close(errChan)
	})

	// initialize a grpc client stub
	var err error
	conn, err := grpc.Dial(c.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to initialize client stub: %w", err)
	}
	defer conn.Close()
	pbClient := server_grpc.NewAppServerClient(conn)

	// read user input and send it to the grpc server
	var cReader cancelreader.CancelReader // bufio.Scanner.Scan() is a blocking call, and it's impossible to close os.Stdin, so linux epoll has to be used, a library for that is used here instead of implementing it myself
	g.Add(func() error {
		cReader, err = cancelreader.NewReader(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to create cancel reader: %w", err)
		}
		scanner := bufio.NewScanner(cReader)
		c.logger.Info("enter message to send")
		for scanner.Scan() {
			line := scanner.Text()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			resp, err := pbClient.Send(ctx, &server_grpc.SendRequest{Message: line})
			if err != nil {
				errChan <- fmt.Errorf("failed to send message: %w", err)
				cancel()
				continue
			}
			c.logger.Info("message sent", "response", resp)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		return nil
	}, func(err error) {
		c.logger.Debug("closing input reading goroutine")
		cReader.Cancel()
	})

	return g.Run()
}
