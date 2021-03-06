package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/server"
)

func getServerCmd() *cli.Command {
	return &cli.Command{
		Name:     "server",
		Aliases:  []string{"s", "serve"},
		Usage:    "start the authentication server",
		HideHelp: true,
		Flags: []cli.Flag{
			// server configuration
			&cli.StringFlag{
				Name:    "host",
				Value:   "",
				EnvVars: []string{"HOST"},
				Usage:   "The `HOST` the server will listen on.",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   3000,
				EnvVars: []string{"PORT"},
				Usage:   "The `PORT` the server will listen to.",
			},

			// ldap server configuration
			&cli.StringFlag{
				Name:    "ldap-host",
				Value:   "ldap://localhost",
				EnvVars: []string{"LDAP_ADDR"},
				Usage:   "The ldap `HOST` (and scheme) the server will authenticate against.",
			},

			// bind dn configuration
			&cli.StringFlag{
				Name:     "bind-dn",
				EnvVars:  []string{"LDAP_BINDDN"},
				Required: true,
				Usage:    "The service account `DN` to do the ldap search.",
			},
			&cli.StringFlag{
				Name:     "bind-credentials",
				EnvVars:  []string{"LDAP_BINDCREDENTIALS"},
				FilePath: "/etc/k8s-ldap-auth/ldap/password",
				Usage:    "The service account `PASSWORD` to do the ldap search, can be located in '/etc/k8s-ldap-auth/ldap/password'.",
			},

			// user search configuration
			&cli.StringFlag{
				Name:    "search-base",
				EnvVars: []string{"LDAP_USER_SEARCHBASE"},
				Usage:   "The `DN` where the ldap search will take place.",
			},
			&cli.StringFlag{
				Name:    "search-filter",
				Value:   "(&(objectClass=inetOrgPerson)(uid=%s))",
				EnvVars: []string{"LDAP_USER_SEARCHFILTER"},
				Usage:   "The `FILTER` to select users.",
			},
			&cli.StringFlag{
				Name:    "member-of-property",
				Value:   "ismemberof",
				EnvVars: []string{"LDAP_USER_MEMBEROFPROPERTY"},
				Usage:   "The `PROPERTY` where group entitlements are located.",
			},
			&cli.StringSliceFlag{
				Name:    "search-attributes",
				Value:   cli.NewStringSlice("uid", "dn", "cn"),
				EnvVars: []string{"LDAP_USER_SEARCHATTR"},
				Usage:   "Repeatable. User `PROPERTY` to fetch. Everything beside 'uid', 'dn', 'cn' (mandatory fields) will be stored in extra values in the UserInfo object.",
			},
			&cli.StringFlag{
				Name:    "search-scope",
				Value:   "sub",
				EnvVars: []string{"LDAP_USER_SEARCHSCOPE"},
				Usage:   "The `SCOPE` of the search. Can take to values base object: 'base', single level: 'single' or whole subtree: 'sub'.",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				port = c.Int("port")
				host = c.String("host")

				ldapURL          = c.String("ldap-host")
				bindDN           = c.String("bind-dn")
				bindPassword     = c.String("bind-credentials")
				searchBase       = c.String("search-base")
				searchScope      = c.String("search-scope")
				searchFilter     = c.String("search-filter")
				searchAttributes = c.StringSlice("search-attributes")
				memberOfProperty = c.String("member-of-property")
			)
			//	return fmt.Errorf("There was an error starting the server, %w", err)

			listen := server.Initialize(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter, memberOfProperty, searchAttributes)

			return listen(fmt.Sprintf("%s:%d", host, port))
		},
	}
}
