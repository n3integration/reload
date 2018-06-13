package actions

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/time/rate"
	"gopkg.in/urfave/cli.v1"

	"github.com/0xAX/notificator"
	"github.com/codegangsta/envy/lib"
	"github.com/fsnotify/fsnotify"
	"github.com/mattn/go-shellwords"
	"github.com/n3integration/reload/runtime"
)

var notifier = notificator.New(notificator.Options{
	AppName: "Reload Build",
})

func Main(c *cli.Context) {
	laddr := c.GlobalString("laddr")
	port := c.GlobalInt("port")
	all := c.GlobalBool("all")
	appPort := strconv.Itoa(c.GlobalInt("appPort"))
	immediate = c.GlobalBool("immediate")
	keyFile := c.GlobalString("keyFile")
	certFile := c.GlobalString("certFile")
	logPrefix := c.GlobalString("logPrefix")
	notifications = c.GlobalBool("notifications")

	logger.SetPrefix(fmt.Sprintf("[%s] ", logPrefix))

	envy.Bootstrap()
	os.Setenv("PORT", appPort)

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	buildArgs, err := shellwords.Parse(c.GlobalString("buildArgs"))
	if err != nil {
		logger.Fatal(err)
	}

	buildPath := c.GlobalString("build")
	if buildPath == "" {
		buildPath = c.GlobalString("path")
	}

	builder := runtime.NewBuilder(buildPath, c.GlobalString("bin"), wd, buildArgs)
	runner := runtime.NewRunner(filepath.Join(wd, builder.Binary()), c.Args()...)
	runner.SetWriter(os.Stdout)
	proxy := runtime.NewProxy(builder, runner)

	config := &runtime.Config{
		Laddr:    laddr,
		Port:     port,
		ProxyTo:  "http://localhost:" + appPort,
		KeyFile:  keyFile,
		CertFile: certFile,
	}

	err = proxy.Run(config)
	if err != nil {
		logger.Fatal(err)
	}

	if laddr != "" {
		logger.Printf("Listening at %s:%d\n", laddr, port)
	} else {
		logger.Printf("Listening on port %d\n", port)
	}

	shutdown(runner)

	// build right now
	build(builder, runner, logger)

	// scan for changes
	scanChanges(c.GlobalString("path"), c.GlobalStringSlice("excludeDir"), all, func(path string) {
		runner.Kill()
		build(builder, runner, logger)
	})
}

func build(builder runtime.Builder, runner runtime.Runner, logger *log.Logger) {
	logger.Println("Building...")
	if notifications {
		notifier.Push("Build Started", "Building "+builder.Binary()+"...", "", notificator.UR_NORMAL)
	}

	if err := builder.Build(); err == nil {
		logger.Printf("%sBuild complete%s\n", colorGreen, colorReset)
		if immediate {
			runner.Run()
		}
		if notifications {
			if err := notifier.Push("Build Succeeded", "Build Complete", "", notificator.UR_CRITICAL); err != nil {
				logger.Println("failed to publish notification")
			}
		}
	} else {
		logger.Printf("%sBuild failed%s\n", colorRed, colorReset)
		fmt.Println(builder.Errors())
		buildErrors := strings.Split(builder.Errors(), "\n")
		if notifications {
			if err := notifier.Push("Build Failed", buildErrors[1], "", notificator.UR_CRITICAL); err != nil {
				logger.Println("failed to publish notification")
			}
		}
	}

	time.Sleep(100 * time.Millisecond)
}

type scanCallback func(path string)

func scanChanges(watchPath string, excludeDirs []string, allFiles bool, cb scanCallback) {
	watcher, _ := fsnotify.NewWatcher()
	throttle := rate.NewLimiter(rate.Every(time.Second*3), 1)
	defer watcher.Close()

	if err := walk(watcher, watchPath, excludeDirs); err != nil {
		logger.Print("error:", err)
	}

	for {
		select {
		case event := <-watcher.Events:
			switch event.Op {
			case fsnotify.Create:
				if info, err := os.Stat(event.Name); err == nil {
					if info.IsDir() {
						walk(watcher, event.Name, excludeDirs)
					}
				}
			case fsnotify.Write:
				if allFiles || filepath.Ext(event.Name) == ".go" {
					if throttle.Allow() {
						cb(event.Name)
					}
				}
			case fsnotify.Remove:
				watcher.Remove(event.Name)
			}
		}
	}
}

func walk(watcher *fsnotify.Watcher, watchPath string, excludeDirs []string) error {
	return filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		for _, x := range excludeDirs {
			if x == path {
				return filepath.SkipDir
			}
		}

		if path == "." {
			return watcher.Add(path)
		}

		if (path == "vendor" || filepath.Base(path)[0] == '.') && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
}

func shutdown(runner runtime.Runner) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("received signal: ", s)
		err := runner.Kill()
		if err != nil {
			log.Print("failed to terminate: ", err)
		}
		f, _ := runner.Info()
		if err := os.Remove(f.Name()); err != nil {
			log.Println("failed to cleanup:", err)
		}
		log.Print("exiting")
		os.Exit(1)
	}()
}
