package server

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

const (
	ScopeBaseObject   = "base"
	ScopeSingleLevel  = "single"
	ScopeWholeSubtree = "sub"
)

var scopeMap = map[string]int{
	ScopeBaseObject:   0,
	ScopeSingleLevel:  1,
	ScopeWholeSubtree: 2,
}

func searchInstance(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter string, searchAttributes []string) func(string, string) (*ldap.Entry, error) {
	return func(username, password string) (*ldap.Entry, error) {
		l, err := ldap.DialURL(ldapURL)
		if err != nil {
			return nil, err
		}

		defer l.Close()

		err = l.Bind(bindDN, bindPassword)
		if err != nil {
			return nil, err
		}

		// Execute LDAP Search request
		searchRequest := ldap.NewSearchRequest(
			searchBase,
			scopeMap[searchScope],
			ldap.NeverDerefAliases, // Dereference aliases
			0,                      // Size limit (0 = no limit)
			0,                      // Time limit (0 = no limit)
			false,                  // Types only
			fmt.Sprintf(searchFilter, username),
			searchAttributes,
			nil, // Additional 'Controls'
		)
		result, err := l.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		// If LDAP Search produced a result, return UserInfo, otherwise, return nil
		if len(result.Entries) == 0 {
			return nil, nil
		} else if len(result.Entries) > 1 {
			return nil, fmt.Errorf("Too many entries returned")
		}

		// Bind as the user to verify their password
		err = l.Bind(result.Entries[0].DN, password)
		if err != nil {
			return nil, err
		}

		// // Rebind as the read only user for any further queries
		// err = l.Bind(bindDN, bindPassword)
		// if err != nil {
		//	return nil, err
		// }

		return result.Entries[0], nil
	}
}
