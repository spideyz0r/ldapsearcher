package main

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
)


func TestResultToJson(t *testing.T) {
	// create a sample SearchResult object
	sr := &ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: "cn=some user,ou=people,dc=sometest,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{Name: "cn", Values: []string{"some user"}},
					{Name: "mail", Values: []string{"someuser@sometest.com"}},
				},
			},
		},
	}

	// create a Search object
	s := Search{
		attributes: []string{"cn", "mail"},
	}

	// call the function with the sample SearchResult object
	jsonStr, err := s.resultToJson(sr)
	if err != nil {
		t.Errorf("Error converting SearchResult to JSON: %v", err)
	}

	// check the output JSON string
	expectedJSON := `{"cn":["some user"],"mail":["someuser@sometest.com"]}`
	if jsonStr != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, jsonStr)
	}
}





