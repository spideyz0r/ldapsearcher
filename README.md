# ldapsearcher [![CI](https://github.com/spideyz0r/ldapsearcher/workflows/gotester/badge.svg)][![CI](https://github.com/spideyz0r/ldapsearcher/workflows/goreleaser/badge.svg)]
A ldap search tool. Run pre-defined or custom queries
```
$ ldapsearcher -h
Usage: search [-chjLp] [-a value] [-b value] [-f value] [-g value] [-G value] [-l value] [-m value] [-o value] [-P value] [-s value] [-u value] [parameters ...]
 -a, --custom-attributes=value
                   list of attributes delimited by ',' for custom search
 -b, --base-dn=value
                   BaseDN
 -c, --custom-search
 -f, --custom-filter=value
                   filter for custom search
 -g, --group-list-members=value
                   list all members of a group -- usually requires full
                   distinguished name of group
 -G, --groups-of-a-group=value
                   list groups that the specified group is a member, wildcards
                   allowed
 -h, --help        display this help
 -j, --single-line-output
                   single-line json output
 -L, --locked      list locked users
 -l, --user-list-groups=value
                   list groups of a user
 -m, --modified-after=value
                   list users created or modified after date YYYYMMDDHHMMSS
 -o, --extra-ou=value
                   extra organizational unit
 -p, --passwd
 -P, --port=value  ldap port
 -s, --server=value
                   ldap server. default: localhost
 -u, --user=value  username
```
