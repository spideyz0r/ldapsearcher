# ldapsearcher [![CI](https://github.com/spideyz0r/ldapsearcher/workflows/gotester/badge.svg)][![CI](https://github.com/spideyz0r/ldapsearcher/workflows/goreleaser/badge.svg)]
A ldap search tool. Run pre-defined or custom queries
```
$ ./ldapsearcher -h
Usage: ldapsearcher [-hkpt] [-a value] [-c value] [-d value] [-f value] [-g value] [-l value] [-m value] [-o value] [-r value] [-s value] [-u value] [parameters ...]
 -a, --custom-attributes=value
                   list of attributes delimited by ',' for custom search
 -c, --modified-after=value
                   list users created or modified after date YYYYMMDDHHMMSS
 -d, --base-dn=value
                   BaseDN
 -f, --custom-filter=value
                   filter for custom search
 -g, --list-groups-user=value
                   list groups of a user
 -h, --help        display this help
 -k, --locked      list locked users
 -l, --list-all-members-group=value
                   list all members of a group -- usually requires full
                   distinguished name of group
 -m, --group-membersof=value
                   list groups that the specified group is a member, wildcards
                   allowed
 -o, --extra-ou=value
                   extra organizational unit
 -p, --passwd
 -r, --port=value  ldap port
 -s, --server=value
                   ldap server
 -t, --custom-search
 -u, --user=value  username
```
