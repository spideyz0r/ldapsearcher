package main

import (
	"testing"
	"reflect"
	"github.com/go-ldap/ldap/v3"
)

func TestFormatResult(t *testing.T) {
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
	attributes := []string{"cn", "email"}

	// Call the function being tested
	result, err := formatResult(sr, attributes)

	// Check the result
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResult := map[string][]string{
		"cn":    {"some user"},
		"email": {"someuser@sometest.com"},
	}

	if reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result: %s got: %v", expectedResult, result)
	}

}



