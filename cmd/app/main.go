package main

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/alecthomas/kong"
	"github.com/mwasilew2/go-service-template/cmd/app/version"
)

var programLevel = new(slog.LevelVar)

type cmdContext struct {
	Logger *slog.Logger
}

type Globals struct {
	LogLevel int                 `short:"l" help:"Log level: 0 (debug), 1 (info), 2 (warn), 3 (error)" default:"1"`
	Version  version.VersionFlag `short:"v" name:"version" help:"Print version information and quit"`
}

var kongApp struct {
	Globals

	Server    serverCmd    `cmd:"" help:"Start the app server."`
	Client    clientCmd    `cmd:"" help:"Start the app client."`
	Transform transformCmd `cmd:"" help:"Transform statistical data into a format easily digestable by an executable."`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)

	kongCtx := kong.Parse(&kongApp,
		kong.Description("A simple application."),
		kong.UsageOnError(),
	)
	switch kongApp.Globals.LogLevel {
	case 0:
		programLevel.Set(slog.LevelDebug)
	case 1:
		programLevel.Set(slog.LevelInfo)
	case 2:
		programLevel.Set(slog.LevelWarn)
	case 3:
		programLevel.Set(slog.LevelError)
	}
	err := kongCtx.Run(&cmdContext{Logger: logger})
	kongCtx.FatalIfErrorf(err)
}
