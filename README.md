# k8s-ldap-auth

## What
This is a webhook token authentication plugin implementation for ldap backend heavily inspired from Daniel Weibel article and own implementation in "Implementing LDAP authentication for Kubernetes" at https://itnext.io/implementing-ldap-authentication-for-kubernetes-732178ec2155

k8s-ldap-auth is actually providing authentication with a token such as `base64(username:password)`, populating v1.UserInfo with: 
```
v1.UserInfo{
  Username: ldapUser.uid,
  UID:      ldapUser.dn,
  Groups:   ldapUser[memberOfProperty],
}
```

## Usage
```
NAME:
   k8s-ldap-auth server - start the authentication server

USAGE:
   k8s-ldap-auth server [command options] [arguments...]

OPTIONS:
   --host HOST                    The HOST the server will listen on. [$HOST]
   --port PORT                    The PORT the server will listen to. (default: 3000) [$PORT]
   --ldap-host HOST               The ldap HOST (and scheme) the server will authenticate against. (default: "ldap://localhost") [$LDAP_ADDR]
   --bind-dn DN                   The service account DN to do the ldap search. [$LDAP_BINDDN]
   --bind-credentials PASSWORD    The service account PASSWORD to do the ldap search, can be located in '/etc/k8s-ldap-auth/ldap/password'. [$LDAP_BINDCREDENTIALS]
   --search-base DN               The DN where the ldap search will take place. [$LDAP_USER_SEARCHBASE]
   --search-filter FILTER         The FILTER to select users. (default: "(&(objectClass=inetOrgPerson)(uid=%s))") [$LDAP_USER_SEARCHFILTER]
   --member-of-property PROPERTY  The PROPERTY where group entitlements are located. (default: "ismemberof") [$LDAP_USER_MEMBEROFPROPERTY]
   --search-attributes PROPERTY   Repeatable. User PROPERTY to fetch. Everything beside 'uid', 'dn', 'cn' (mandatory fields) will be stored in extra values in the UserInfo object. (default: "uid", "dn", "cn") [$LDAP_USER_SEARCHATTR]
   --search-scope SCOPE           The SCOPE of the search. Can take to values base object: 'base', single level: 'single' or whole subtree: 'sub'. (default: "sub") [$LDAP_USER_SEARCHSCOPE]
```

## What's next
 - Group search for ldap not supporting memberof attribute
 - Entitlement check with ldap groups based on a configuration file
 - Some kind of authorization check based on groups and dn ?
