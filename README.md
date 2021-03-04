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
   --host HOST                  The HOST the server will listen on. [$HOST]
   --port PORT                  The PORT the server will listen to. (default: 443) [$PORT]
   --ldap-host HOST             The ldap HOST (and scheme) the server will authenticate against. (default: "ldap://localhost") [$LDAP_ADDR]
   --bind-dn DN                 The service account DN to do the ldap search [$LDAP_BINDDN]
   --bind-credentials PASSWORD  The service account PASSWORD to do the ldap search [$LDAP_BINDCREDENTIALS]
   --search-base value           [$LDAP_USER_SEARCHBASE]
   --search-filter value        (default: "(&(objectClass=inetOrgPerson)(uid=%s))") [$LDAP_USER_SEARCHFILTER]
   --member-of-property value   (default: "ismemberof") [$LDAP_USER_MEMBEROFPROPERTY]
   --search-attributes value    (default: "uid", "dn", "cn") [$LDAP_USER_SEARCHATTR]
   --search-scope value         (default: "sub") [$LDAP_USER_SEARCHSCOPE]
```

Despite the default port being 443, for now the server only knows http.

The bindDN password value will be fetch from `/etc/k8s-ldap-auth/ldap/password`.

## What's next
 - SSL
 - Group search for ldap not supporting memberof attribute
 - Entitlement check with ldap groups based on a configuration file
 - Some kind of RoleBinding based on groups and dn ?
