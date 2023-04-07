package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pborman/getopt"
	"golang.org/x/term"
)

func main() {
	const (
		userGroupSearchFilter  = "(mail=%s)"
		groupMemberListFilter  = "(memberOf=cn=%s, %s, %s)"
		modifiedAfterFilter    = "(&(whenChanged>=%s.0-0500)(objectClass=user)(objectCategory=person))"
		usersLockedFilter      = "(lockoutTime>=1)"
		groupMemberOfFilter    = "(cn=%s)"
	)

	help := getopt.BoolLong("help", 'h', "display this help")
	ldap_server := getopt.StringLong("server", 's', "localhost", "ldap server. default: localhost")
	ldap_port := getopt.StringLong("port", 'P', "389", "ldap port")
	ldap_user := getopt.StringLong("user", 'u', "", "username")
	ldap_passwd := getopt.BoolLong("passwd", 'p', "", "password")
	base_dn := getopt.StringLong("base-dn", 'b', "", "BaseDN")
	group_memberof := getopt.StringLong("groups-of-a-group", 'G', "", "list groups that the specified group is a member, wildcards allowed")
	user_group_search := getopt.StringLong("user-list-groups", 'l', "", "list groups of a user")
	group_member_list := getopt.StringLong("group-list-members", 'g', "", "list all members of a group -- usually requires full distinguished name of group")
	extra_ou := getopt.StringLong("extra-ou", 'o', "", "extra organizational unit")
	modified_after := getopt.StringLong("modified-after", 'm', "", "list users created or modified after date YYYYMMDDHHMMSS")
	// recursive_groups_user := getopt.StringLong("recursive-list-group", 'r', "", "list nested groups for a user")
	users_locked := getopt.BoolLong("locked", 'L', "list locked users")
	custom_search := getopt.BoolLong("custom-search", 'c', "", "custom search")
	custom_filter := getopt.StringLong("custom-filter", 'f', "", "filter for custom search")
	custom_attributes := getopt.StringLong("custom-attributes", 'a', "", "list of attributes delimited by ',' for custom search")

	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	var pass string
	if *ldap_passwd {
		fmt.Print("Enter your password: ")
		pass_b, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		pass = string(pass_b)
	} else {
		pass = os.Getenv("ADP")
	}

	conn, err := ldapConnect(*ldap_server+":"+*ldap_port, *ldap_user, pass)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var s Search
	switch {
	case *custom_search:
		if len(*custom_filter) == 0 || len(*custom_attributes) == 0 {
			log.Fatal("Custom search requires valid filter and attributes")
		}
		s = Search{
			base_dn:    *base_dn,
			filter:     *custom_filter,
			attributes: strings.Split(*custom_attributes, ","),
		}
	case len(*user_group_search) > 0:
		s = Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf(userGroupSearchFilter, *user_group_search),
			attributes: []string{"sAMAccountName", "memberOf"},
		}
	case len(*group_member_list) > 0:
		s = Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf(groupMemberListFilter, *group_member_list, *extra_ou, *base_dn),
			attributes: []string{"sAMAccountName"},
		}
	case len(*modified_after) > 0:
		s = Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf(modifiedAfterFilter, *modified_after),
			attributes: []string{"sAMAccountName"},
		}
	case *users_locked:
		s = Search{
			base_dn:    *base_dn,
			filter:     usersLockedFilter,
			attributes: []string{"sAMAccountName"},
		}
	case len(*group_memberof) > 0:
		s = Search{
			base_dn:    *base_dn,
			filter:     fmt.Sprintf(groupMemberOfFilter, *group_memberof),
			attributes: []string{"cn", "sAMAccountName", "memberOf"},
		}
	default:
		log.Fatal("No search criteria provided")
	}

	result_json, err := s.newSearch(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result_json)
	os.Exit(0)
}


// TODO
// re-implement the recursive search
// rpm package pipeline