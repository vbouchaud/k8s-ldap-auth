package main

import (
	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
