package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "net/http/pprof"

	"golang.org/x/exp/slog"

	oapi_middleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	server_grpc "github.com/mwasilew2/go-service-template/gen/server-grpc"
	server_oapi "github.com/mwasilew2/go-service-template/gen/server-oapi"
	"github.com/oklog/run"
	slogecho "github.com/samber/slog-echo"
	"google.golang.org/grpc"
)

type serverCmd struct {
	// cli options
	HttpAddr string `help:"address which the http server should listen on" default:":8080"`
	GrpcAddr string `help:"address which the grpc server should listen on" default:":8081"`

	// Dependencies
	logger *slog.Logger

	// Embedded types
	server_grpc.UnimplementedAppServerServer
}

func (c *serverCmd) PostV1Message(ctx context.Context, request server_oapi.PostV1MessageRequestObject) (server_oapi.PostV1MessageResponseObject, error) {
	return nil, nil
}

func (c *serverCmd) GetV1MessageId(ctx context.Context, request server_oapi.GetV1MessageIdRequestObject) (server_oapi.GetV1MessageIdResponseObject, error) {
	return server_oapi.GetV1MessageId200JSONResponse{
		Id:      2,
		Message: "Hello World!",
	}, nil
}

func (c *serverCmd) Run(cmdCtx *cmdContext) error {
	c.logger = cmdCtx.Logger.With("component", "serverCmd")

	// create a run group
	g := run.Group{}

	// start http server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(slogecho.New(c.logger.With("subcomponent", "echo")))
	e.Use(echoprometheus.NewMiddleware("echo"))
	swagger, err := server_oapi.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to get swagger: %w", err)
	}
	e.Use(oapi_middleware.OapiRequestValidatorWithOptions(swagger, &oapi_middleware.Options{
		Skipper: func(ctx echo.Context) bool {
			path := ctx.Request().URL.Path
			if path == "/metrics" {
				return true
			}
			if strings.HasPrefix(path, "/debug") {
				return true
			}
			return false
		},
	}))
	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	strictSrv := server_oapi.NewStrictHandler(c, nil)
	server_oapi.RegisterHandlersWithBaseURL(e, strictSrv, "/api")
	g.Add(func() error {
		c.logger.Info("starting http server", "address", c.HttpAddr)
		return e.Start(c.HttpAddr)
	}, func(err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		c.logger.Debug("shutting down http server")
		if err := e.Shutdown(ctx); err != nil {
			c.logger.Error("failed to shutdown http server", "error", err)
			return
		}
		c.logger.Debug("http server stopped")
	})

	// start grpc server
	var srv *grpc.Server
	lis, err := net.Listen("tcp", c.GrpcAddr)
	if err != nil {
		return fmt.Errorf("tcp failed to listen on: %w", err)
	}
	srv = grpc.NewServer()
	server_grpc.RegisterAppServerServer(srv, c)
	g.Add(func() error {
		c.logger.Info("starting grpc server", "address", c.GrpcAddr)
		return srv.Serve(lis)
	}, func(err error) {
		c.logger.Debug("shutting down grpc server")
		srv.Stop()
		c.logger.Debug("grpc server stopped")
	})

	// listen for termination signals
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Kill, os.Interrupt)
	done := make(chan struct{})
	g.Add(func() error {
		select {
		case sig := <-osSigChan:
			c.logger.Debug("caught signal", "signal", sig.String())
			return fmt.Errorf("caught signal: %s", sig.String())
		case <-done:
			c.logger.Debug("signal catching goroutine stopped")
		}
		return nil
	}, func(err error) {
		close(done)
	})

	return g.Run()
}

func (c *serverCmd) Send(ctx context.Context, req *server_grpc.SendRequest) (*server_grpc.SendResponse, error) {
	c.logger.Debug("received Send request", "req", req)
	return &server_grpc.SendResponse{
		Status: 200,
	}, nil
}
