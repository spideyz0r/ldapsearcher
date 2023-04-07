package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type Search struct {
	base_dn     string
	filter      string
	attributes  []string
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

func formatResult(sr *ldap.SearchResult, attributes []string) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(sr.Entries) <= 0 {
		return result, nil
	}
	for _, e := range sr.Entries {
		for _, a := range attributes {
			var results []string
			for _, r := range e.GetAttributeValues(a) {
				results = append(results, r)
			}
			if len(result[a]) > 0 {
				result[a] = append(result[a], results...)
			} else {
				result[a] = results
			}
		}
	}
	return result, nil
}

func printResult(result map[string][]string, json_output bool) (error) {
	j, err := json.Marshal(result)
	if err != nil {
		return err
	}

	if json_output {
		fmt.Printf("%s", string(j))
		return nil
	}
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, j, "", "\t")
	fmt.Println(prettyJSON.String())
	return nil
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
