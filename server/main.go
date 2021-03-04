package server

import (
	"github.com/go-ldap/ldap"
)

func Initialize(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter, memberOfProperty string, searchAttributes []string) func(string) error {
	ls := searchInstance(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter, append([]string{memberOfProperty}, searchAttributes...))

	validateEntitlement := func(*ldap.Entry) bool {
		return true
	}

	return func(address string) error {
		return listen(address, ls, memberOfProperty, validateEntitlement)
	}
}
