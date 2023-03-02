package main

import (
    "fmt"
    "os"
    "github.com/pborman/getopt"
    "github.com/go-ldap/ldap/v3"
)

func main () {
    help := getopt.BoolLong("help", 'h', "display this help")
    ldap_server := getopt.StringLong("server", 's', "localhost", "ldap server")
    ldap_port := getopt.StringLong("port", 'p', "389", "ldap port")
    ldap_user := getopt.StringLong("user", 'u', "", "username")
    group_search := getopt.StringLong("list-members-group", 'l', "", "list members of a group")
    getopt.Parse()

    if *help {
        getopt.Usage()
        os.Exit(0)
    }

    conn, err := ldap_connect(*ldap_server + ":" + *ldap_port)
    if err != nil {
        fmt.Printf("Failed to connect. %s\n", err)
        os.Exit(2)
    }
    defer conn.Close()

    if len(*group_search) >0 {
        listMembersGroup(*group_search)
    }

    fmt.Printf(*ldap_user, *ldap_server)
    os.Exit(0)
}

func ldap_connect(s string)(*ldap.Conn, error){
    conn, err := ldap.Dial("tcp", s)
    fmt.Printf("Connecting to LDAP...")
    if err != nil {
        return nil, fmt.Errorf("Failed to connect to server: %s\n", err)
    }
    if err := conn.Bind("somebind", "somepassword"); err != nil {
        return nil, fmt.Errorf("Failed to bind: %s\n", err)
    }
    fmt.Printf("Connected!")
    return conn, nil
}

// search groups from user
func groupSearch () {
    fmt.Printf("Searching group\n")
}

// search members from group
func listMembersGroup (s string) {
    fmt.Printf("Listing members for group: %s \n", s)
}

// Groups starting with certain string

// Users created after x date

// Users that belong to a group (Recursively)

// Disabled user accounts

// Locked user accounts
