package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/pborman/getopt"
)

func main() {
	help := getopt.BoolLong("help", 'h', "display this help")
	ldap_server := getopt.StringLong("server", 's', "localhost", "ldap server")
	ldap_port := getopt.StringLong("port", 'p', "389", "ldap port")
	ldap_user := getopt.StringLong("user", 'u', "", "username")
	base_dn := getopt.StringLong("base-dn", 'd', "", "BaseDN")
	group_memberof := getopt.StringLong("group-membersof", 'm', "", "list groups that the specified group is a member, wildcards allowed")
	user_group_search := getopt.StringLong("list-groups-user", 'g', "", "list groups of a user")
	group_member_list := getopt.StringLong("list-all-members-group", 'l', "", "list all members of a group -- usually requires full distinguished name of group")
	extra_ou := getopt.StringLong("extra-ou", 'o', "", "extra organizational unit")
	modified_after := getopt.StringLong("modified-after", 'c', "", "list users created or modified after date YYYYMMDDHHMMSS")
	recursive_groups_user := getopt.StringLong("recursive-list-group", 'r', "", "list nested groups for a user")
	users_locked := getopt.BoolLong("locked", 'k', "list locked users")

	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	conn, err := ldapConnect(*ldap_server+":"+*ldap_port, *ldap_user)
	if err != nil {
		fmt.Printf("Failed to connect. %s\n", err)
		os.Exit(2)
	}
	defer conn.Close()

	if len(*group_memberof) > 0 {
		memberOfGroup(conn, *group_memberof, *base_dn, true)
		os.Exit(0)
	}

	if len(*user_group_search) > 0 {
		listGroupsUser(conn, *user_group_search, *base_dn)
		os.Exit(0)
	}

	if len(*group_member_list) > 0 {
		listAllMembersGroup(conn, *group_member_list, *base_dn, *extra_ou)
		os.Exit(0)
	}

	if len(*modified_after) > 0 {
		usersModifiedAfterDate(conn, *modified_after, *base_dn)
		os.Exit(0)
	}

	if len(*recursive_groups_user) > 0 {
		listNestedGroupsUser(conn, *recursive_groups_user, *base_dn, *extra_ou)
		os.Exit(0)
	}

	if *users_locked {
		lockedUsers(conn, *base_dn)
		os.Exit(0)
	}
	os.Exit(0)
}

func ldapConnect(server string, user string) (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", server)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to server: %s\n", err)
	}
	if err := conn.Bind(user, os.Getenv("ADP")); err != nil {
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

func printOut(sr *ldap.SearchResult, attributes []string, header bool) {
	if len(sr.Entries) <= 0 {
		fmt.Printf("No results found.\n")
	}
	for _, s := range sr.Entries {
		for _, a := range attributes {
			if header {
				fmt.Printf("%s:\n", a)
			}
			for _, i := range s.GetAttributeValues(a) {
				fmt.Printf("\t%s\n", i)
			}
		}
		fmt.Printf("\n")
	}
}

func memberOfGroup(conn *ldap.Conn, group string, base_dn string, printout bool) (*ldap.SearchResult, error) {
	sr, err := runQuery(conn, base_dn, "(cn="+group+")", []string{"cn", "sAMAccountName", "memberOf"})
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}

	// sr.Entries[0].Print()
	if printout {
		printOut(sr, []string{"cn", "memberOf"}, true)
	}
	return sr, nil
}

func listGroupsUser(conn *ldap.Conn, user string, base_dn string) {
	sr, err := runQuery(conn, base_dn, "(mail="+user+")", []string{"sAMAccountName", "memberOf"})
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}
	printOut(sr, []string{"sAMAccountName", "memberOf"}, true)
}

func listNestedGroupsUser(conn *ldap.Conn, user string, base_dn string, extra_ou string) {
	sr, err := runQuery(conn, base_dn, "(mail="+user+")", []string{"sAMAccountName", "memberOf"})
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}

	level := []string{">"}
	for _, s := range sr.Entries {
		for _, i := range s.GetAttributeValues("memberOf") {
			group := strings.Split(strings.Split(i, "=")[1], ",")[0]
			fmt.Printf(">>>>>> %s\n", group)
			recursiveSearch(conn, base_dn, group, level)
		}
	}
}

func recursiveSearch(conn *ldap.Conn, base_dn string, g string, level []string) {
	sr, err := memberOfGroup(conn, g, base_dn, false)
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}
	for _, s := range sr.Entries {
		for _, i := range s.GetAttributeValues("memberOf") {
			group := strings.Split(strings.Split(i, "=")[1], ",")[0]
			level = append(level, ">")
			fmt.Printf("Inner group %v:  %s\n", level, group)
			recursiveSearch(conn, base_dn, group, level)
		}
	}
}

func listAllMembersGroup(conn *ldap.Conn, group string, base_dn string, extra_ou string) {
	sr, err := runQuery(conn, base_dn, "(memberOf=cn="+group+","+extra_ou+","+base_dn+")", []string{"sAMAccountName"})
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}
	printOut(sr, []string{"sAMAccountName"}, false)
}

func usersModifiedAfterDate(conn *ldap.Conn, date string, base_dn string) {
	sr, err := runQuery(conn, base_dn, "(&(whenChanged>="+date+".0-0500)(objectClass=user)(objectCategory=person))", []string{"sAMAccountName"})
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}
	printOut(sr, []string{"sAMAccountName"}, false)
}

func lockedUsers(conn *ldap.Conn, base_dn string) {
	sr, err := runQuery(conn, base_dn, "(lockoutTime>=1)", []string{"sAMAccountName"})
	if err != nil {
		fmt.Printf("Query failed %s\n", err)
		os.Exit(2)
	}
	printOut(sr, []string{"sAMAccountName"}, false)
}

// TODO
// deal with optional args
// accept password via read or stdin
// deal with json
