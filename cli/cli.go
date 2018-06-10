package cli

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

var (
	logger        = log.New(os.Stdout, "[reload] ", 0)
	immediate     = false
	colorGreen    = string([]byte{27, 91, 57, 55, 59, 51, 50, 59, 49, 109})
	colorRed      = string([]byte{27, 91, 57, 55, 59, 51, 49, 59, 49, 109})
	colorReset    = string([]byte{27, 91, 48, 109})
	notifications = false
)

func Exec() {
	app := cli.NewApp()
	app.Name = "reload"
	app.Usage = "A live reload utility for Go web applications."
	app.Action = MainAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "laddr,l",
			Value:  "",
			EnvVar: "RELOAD_LADDR",
			Usage:  "Listening address for the proxy server",
		},
		cli.IntFlag{
			Name:   "port,p",
			Value:  3000,
			EnvVar: "RELOAD_PORT",
			Usage:  "Port for the proxy server",
		},
		cli.IntFlag{
			Name:   "appPort,a",
			Value:  3001,
			EnvVar: "BIN_APP_PORT",
			Usage:  "Port for the Go web server",
		},
		cli.StringFlag{
			Name:   "bin,b",
			Value:  "reload-bin",
			EnvVar: "RELOAD_BIN",
			Usage:  "Name of generated binary file",
		},
		cli.StringFlag{
			Name:   "path,t",
			Value:  ".",
			EnvVar: "RELOAD_PATH",
			Usage:  "Path to watch files from",
		},
		cli.StringFlag{
			Name:   "build,d",
			Value:  "",
			EnvVar: "RELOAD_BUILD",
			Usage:  "Path to build files from (defaults to same value as --path)",
		},
		cli.StringSliceFlag{
			Name:   "excludeDir,x",
			Value:  &cli.StringSlice{},
			EnvVar: "RELOAD_EXCLUDE_DIR",
			Usage:  "Relative directories to exclude",
		},
		cli.BoolFlag{
			Name:   "immediate,i",
			EnvVar: "RELOAD_IMMEDIATE",
			Usage:  "Run the server immediately after it's built",
		},
		cli.BoolFlag{
			Name:   "all",
			EnvVar: "RELOAD_ALL",
			Usage:  "Reloads whenever any file changes, as opposed to reloading only on .go file change",
		},
		cli.StringFlag{
			Name:   "buildArgs",
			EnvVar: "RELOAD_BUILD_ARGS",
			Usage:  "Additional go build arguments",
		},
		cli.StringFlag{
			Name:   "certFile",
			EnvVar: "RELOAD_CERT_FILE",
			Usage:  "TLS Certificate",
		},
		cli.StringFlag{
			Name:   "keyFile",
			EnvVar: "RELOAD_KEY_FILE",
			Usage:  "TLS Certificate Key",
		},
		cli.StringFlag{
			Name:   "logPrefix",
			EnvVar: "RELOAD_LOG_PREFIX",
			Usage:  "Log prefix",
			Value:  "reload",
		},
		cli.BoolFlag{
			Name:   "notifications",
			EnvVar: "RELOAD_NOTIFICATIONS",
			Usage:  "Enables desktop notifications",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the reload proxy in the current working directory",
			Action:    MainAction,
		},
		{
			Name:      "env",
			ShortName: "e",
			Usage:     "Display environment variables set by the .env file",
			Action:    EnvAction,
		},
	}
	app.Run(os.Args)
}
