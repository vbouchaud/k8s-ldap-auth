package cmd

import (
	"fmt"
	"os"

	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/version"
	"github.com/urfave/cli/v2"
)

type action func(*cli.Context) error

// Start the cmd application
func Start() error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(version.Version())
	}

	app := cli.NewApp()
	app.Name = "ldap-auth"
	app.Version = version.VERSION
	app.Compiled = version.Compiled()
	app.Authors = []*cli.Author{
		{
			Name:  "Vianney Bouchaud",
			Email: "vianney@bouchaud.org",
		},
	}

	app.UseShortOptionHandling = true
	app.Commands = []*cli.Command{
		getServerCmd(),
	}

	return app.Run(os.Args)
}
