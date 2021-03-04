package server

import (
	"fmt"
	"log"
	"strings"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-ldap/ldap"
	"k8s.io/api/authentication/v1"
)

func writeError(w http.ResponseWriter, err error) {
	err = fmt.Errorf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError) // 500
	fmt.Fprintln(w, err)
	log.Println(err)
}

func toLower(a []string) []string {
	var res []string

	for _, item := range a {
		res = append(res, strings.ToLower(item))
	}

	return res
}

func listen(address string, findUser func(string, string) (*ldap.Entry, error), memberOfProperty string, validateEntitlement func(*ldap.Entry) bool) error {
	handler := func(w http.ResponseWriter, r *http.Request) {

		// Read body of POST request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		log.Printf("Receiving: %s\n", string(b))

		// Unmarshal JSON from POST request to TokenReview object
		// TokenReview: https://github.com/kubernetes/api/blob/master/authentication/v1/types.go
		var tr v1.TokenReview
		err = json.Unmarshal(b, &tr)
		if err != nil {
			writeError(w, err)
			return
		}

		// Extract username and password from the token in the TokenReview object
		token, err := base64.StdEncoding.DecodeString(tr.Spec.Token)
		if err != nil {
			writeError(w, err)
			return
		}

		s := strings.SplitN(string(token), ":", 2)
		if len(s) != 2 {
			writeError(w, fmt.Errorf("badly formatted token: %s", tr.Spec.Token))
			return
		}
		username, password := s[0], s[1]

		// Make LDAP Search request with extracted username and password
		userInfo, err := findUser(username, password)
		if err != nil {
			writeError(w, fmt.Errorf("failed LDAP Search request: %v", err))
			return
		}

		// Set status of TokenReview object
		if userInfo == nil {
			tr.Status.Authenticated = false
			log.Printf("User not found.")
		} else if validateEntitlement(userInfo) == false {
			tr.Status.Authenticated = false
			log.Printf("User not entitled.")
		} else {
			log.Printf("User found: %s\n", userInfo.DN)

			tr.Status.Authenticated = true
			tr.Status.User = v1.UserInfo{
				Username: strings.ToLower(userInfo.GetAttributeValue("uid")),
				UID:      strings.ToLower(userInfo.DN),
				Groups:   toLower(userInfo.GetAttributeValues(memberOfProperty)),
			}
		}

		// Marshal the TokenReview to JSON and send it back
		b, err = json.Marshal(tr)
		if err != nil {
			writeError(w, err)
			return
		}
		w.Write(b)
		log.Printf("Returning: %s\n", string(b))
	}

	http.HandleFunc("/", handler)

	log.Printf("Now listening to '%s'\n", address)
	return http.ListenAndServe(address, nil)
}
