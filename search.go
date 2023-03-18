package main

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type Search struct {
	base_dn    string
	filter     string
	attributes []string
	result     *ldap.SearchResult
}

func (s Search) newSearch(conn *ldap.Conn) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(
		s.base_dn,
		2, 0, 0, 0, false,
		s.filter,
		s.attributes,
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to run search: %s\n", err)
	}
	return sr, nil
}

func ldapConnect(server string, user string, pass string) (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", server)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to server: %s\n", err)
	}
	if err := conn.Bind(user, pass); err != nil {
		return nil, fmt.Errorf("Failed to bind: %s\n", err)
	}

	return conn, nil
}

func runQuery(conn *ldap.Conn, base_dn string, filter string, attributes []string) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(
		base_dn,
		2, 0, 0, 0, false,
		filter,
		attributes,
		nil,
	)
	sr, err := conn.Search(searchRequest)

	if err != nil {
		return nil, fmt.Errorf("Failed to run search: %s\n", err)
	}
	return sr, nil
}
