package main

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/ldap.v3"
)

const (
	ldapServer := "your-server.domain.com"
	domain     = "domain.com"
)

func authenticate(username, password string) (User, error) {
	var u User

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, 389))
	if err != nil {
		return u, err
	}
	defer conn.Close()

	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return u, err
	}

	err = conn.Bind(fmt.Sprintf("%s@%s", username, domain), password)
	if err != nil {
		return u, err
	}

	searchRequest := ldap.NewSearchRequest(
		"dc=<domain>,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s))", username),
		[]string{"c", "cn", "displayName", "userPrincipalName"},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return u, err
	}

	if len(searchResult.Entries) != 1 {
		return u, err
	}

	u = User{
		searchResult.Entries[0].GetAttributeValue("cn"),
		searchResult.Entries[0].GetAttributeValue("userPrincipalName"),
		searchResult.Entries[0].GetAttributeValue("displayName"),
	}

	return u, nil
}
