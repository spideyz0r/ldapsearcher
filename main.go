package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
	"github.com/pborman/getopt"
)

func main () {
    help := getopt.BoolLong("help", 'h', "display this help")
    ldap_server := getopt.StringLong("server", 's', "localhost", "ldap server")
    ldap_port := getopt.StringLong("port", 'p', "389", "ldap port")
    ldap_user := getopt.StringLong("user", 'u', "", "username")
    base_dn := getopt.StringLong("base-dn", 'd', "", "BaseDN")
    group_search := getopt.StringLong("list-members-group", 'l', "", "list members of a group")
    getopt.Parse()

    if *help {
        getopt.Usage()
        os.Exit(0)
    }

    conn, err := ldapConnect(*ldap_server + ":" + *ldap_port, *ldap_user)
    if err != nil {
        fmt.Printf("Failed to connect. %s\n", err)
        os.Exit(2)
    }
    defer conn.Close()

    if len(*group_search) >0 {
        memberOfGroup(conn, *group_search, *base_dn)
    }

    fmt.Printf("%s %s \n", *ldap_user, *ldap_server)
    os.Exit(0)
}

func ldapConnect(server string, user string)(*ldap.Conn, error){

    conn, err := ldap.Dial("tcp", server)
    fmt.Printf("Connecting to LDAP...\n")
    if err != nil {
        return nil, fmt.Errorf("Failed to connect to server: %s\n", err)
    }
    if err := conn.Bind(user, os.Getenv("ADP")); err != nil {
        return nil, fmt.Errorf("Failed to bind: %s\n", err)
    }
    fmt.Printf("Connected!\n")

    return conn, nil
}

// search members from group
func memberOfGroup (conn *ldap.Conn, group string, base_dn string) {
    fmt.Printf("bah\n")
    searchRequest := ldap.NewSearchRequest(
        base_dn,
        2, 0, 0, 0, false,
        "(cn="+ group +")",
        []string{"cn", "sAMAccountName", "memberOf"},
        nil,
    )

    sr, err := conn.Search(searchRequest)
    if err != nil {
        log.Fatal(err)
    }

    for _, s:= range sr.Entries {
        fmt.Printf("----\n")
        s.Print()
        fmt.Printf("----\n")
    }
    // sr.Entries[0].Print()
    // fmt.Printf("%+v\n", sr)
}

// search groups from user
func groupSearch () {
    fmt.Printf("Searching group\n")
}


// TODO
// deal with optional args
// accept password via read or stdin

// Groups starting with certain string

// Users created after x date

// Users that belong to a group (Recursively)

// Disabled user accounts

// Locked user accounts
