package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/pborman/getopt"
)

func main() {
	const (
		emailFilter    = "(mail=%s)"
		memberOfFilter = "(memberOf=cn=%s,%s,%s)"
	)

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
	custom_search := getopt.BoolLong("custom-search", 't', "", "custom search")
	custom_filter := getopt.StringLong("custom-filter", 'f', "", "filter for custom search")
	custom_attributes := getopt.StringLong("custom-attributes", 'a', "", "list of attributes delimited by ',' for custom search")

	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	conn, err := ldapConnect(*ldap_server+":"+*ldap_port, *ldap_user)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if *custom_search {
		if len(*custom_filter) == 0 || len(*custom_attributes) == 0 {
			log.Fatal("Custom search requires valid filter and attributes")
		}
		s := Search{
			base_dn:    *base_dn,
			filter:     *custom_filter,
			attributes: strings.Split(*custom_attributes, ","),
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		os.Exit(0)
	}

	// switch {
	// case  len(*user_group_search) > 0:
	// ...
	// case len(*group_member_list) > 0:
	// ...
	// case

	// len(*modified_after) > 0

	if len(*recursive_groups_user) > 0 {
		s := Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf("(mail=%s)", *recursive_groups_user),
			attributes: []string{"sAMAccountName", "memberOf"},
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		listNestedGroupsUser(conn, s, sr)

		// if len(*recursive_groups_user) > 0 {
		// 	err = listNestedGroupsUser(conn, *recursive_groups_user, *base_dn, *extra_ou)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	os.Exit(0)
	}

	if len(*user_group_search) > 0 {
		s := Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf("(mail=%s)", *user_group_search),
			attributes: []string{"sAMAccountName", "memberOf"},
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		os.Exit(0)
	}

	if len(*group_member_list) > 0 {
		s := Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf("(memberOf=cn=%s, %s, %s)", *group_member_list, *extra_ou, *base_dn),
			attributes: []string{"sAMAccountName"},
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		os.Exit(0)
	}

	if len(*modified_after) > 0 {
		s := Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf("(&(whenChanged>=%s.0-0500)(objectClass=user)(objectCategory=person))", *modified_after),
			attributes: []string{"sAMAccountName"},
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		os.Exit(0)
	}

	if *users_locked {
		s := Search{
			base_dn:    *base_dn,
			filter:     "(lockoutTime>=1)",
			attributes: []string{"sAMAccountName"},
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		os.Exit(0)
	}

	if len(*group_memberof) > 0 {
		s := Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf("(cn=%s)", *group_memberof),
			attributes: []string{"cn", "sAMAccountName", "memberOf"},
		}
		sr, err := s.newSearch(conn)
		if err != nil {
			log.Fatal(err)
		}
		printOut(sr, s.attributes, true)
		os.Exit(0)
	}
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

func listGroupsUser(conn *ldap.Conn, user string, base_dn string) error {
	attrs := []string{"sAMAccountName", "memberOf"}

	sr, err := runQuery(conn, base_dn, "(mail="+user+")", attrs)
	if err != nil {
		return fmt.Errorf("query failed %s", err)
	}
	printOut(sr, attrs, true)
	return nil
}

type M map[string]interface{}

func listNestedGroupsUser(conn *ldap.Conn, s Search, sr *ldap.SearchResult) error {
	var myM []M
	for _, e := range sr.Entries {
		for _, i := range e.GetAttributeValues("memberOf") {
			group := strings.Split(strings.Split(i, "=")[1], ",")[0]
			fmt.Printf("%s\n", group)
			myM = append(myM, recursiveSearch(conn, s.base_dn, group))
			// if err != nil {
			//     return fmt.Errorf("query failed %s", err)
			// }
		}
	}
	fmt.Printf("%v", myM)
	return nil
}

func recursiveSearch(conn *ldap.Conn, b string, g string) M {
	var myM M
	s := Search{
		base_dn:    b,
		filter:     fmt.Sprintf("(cn=%s)", g),
		attributes: []string{"memberOf"},
	}
	sr, err := s.newSearch(conn)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range sr.Entries {
		for _, i := range e.GetAttributeValues("memberOf") {
			group := strings.Split(strings.Split(i, "=")[1], ",")[0]
			// fmt.Printf(">>>>>>>>>Inner group %s\n",group)
			myM[group] = recursiveSearch(conn, group, s.base_dn)
			// if err != nil {
			//     return nil, fmt.Errorf("query failed %s", err)
			// }
		}
	}
	// fmt.Printf("%v", all_groups)
	return myM
}

// func recursiveSearch(conn *ldap.Conn, b string, g string, m map[string][]string) ([]string error) {
// 	s := Search {
// 		base_dn: b,
// 		filter: fmt.Sprintf("(cn=%s)", g),
// 		attributes: []string{"memberOf"},
// 	}
// 	sr, err := s.newSearch(conn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, s := range sr.Entries {
// 		for _, i := range s.GetAttributeValues("memberOf") {
// 			group := strings.Split(strings.Split(i, "=")[1], ",")[0]
// 			all_groups = append(all_groups, group)
// 			// fmt.Printf(">>>>>>>>>Inner group %s\n", group)
// 			m[group], err = recursiveSearch(conn, s.base_dn, group, all_groups)

// 			err := recursiveSearch(conn, b, group, all_groups)
// 			if err != nil {
//                 return nil, fmt.Errorf("query failed %s", err)
//             }
// 		}
// 	}
// 	fmt.Printf("%v", all_groups)
//     return nil,
// }

// TODO
// deal with optional args
// accept password via read or stdin
// deal with json
// add custom queries
