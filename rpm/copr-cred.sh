#/bin/bash
cat <<EOF >copr
[copr-cli]
login = $COPR_LOGIN
username = $COPR_USERNAME
token = $COPR_TOKEN
copr_url = https://copr.fedorainfracloud.org
EOF