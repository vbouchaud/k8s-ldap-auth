module vbouchaud/k8s-ldap-auth

go 1.16

require (
	github.com/adrg/xdg v0.4.0
	github.com/etherlabsio/healthcheck/v2 v2.0.0
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/gorilla/mux v1.8.1
	github.com/lestrrat-go/jwx v1.2.26
	github.com/mattn/go-isatty v0.0.20
	github.com/rs/zerolog v1.31.0
	github.com/urfave/cli/v2 v2.27.1
	github.com/zalando/go-keyring v0.2.3
	golang.org/x/term v0.16.0
	k8s.io/api v0.28.4
	k8s.io/apimachinery v0.29.3
	k8s.io/client-go v0.28.4
)
