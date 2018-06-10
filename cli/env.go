package cli

import (
	"fmt"

	"github.com/codegangsta/envy/lib"
	"gopkg.in/urfave/cli.v1"
)

func EnvAction(c *cli.Context) {
	logPrefix := c.GlobalString("logPrefix")
	logger.SetPrefix(fmt.Sprintf("[%s] ", logPrefix))

	// Bootstrap the environment
	env, err := envy.Bootstrap()
	if err != nil {
		logger.Fatalln(err)
	}

	for k, v := range env {
		fmt.Printf("%s: %s\n", k, v)
	}
}
