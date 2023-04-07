#!/bin/bash
echo "Recreating rpmbuild directory"
rm -rvf /root/rpmbuild/
rpmdev-setuptree
echo "Copying over sources"
cp -rpv /project/ldapsearcher/{go.mod,go.sum,knife.go,main.go} /root/rpmbuild/SOURCES
echo "Building SRPM"
rpmbuild --undefine=_disable_source_fetch -bs /project/ldapsearcher/rpm/ldapsearcher.spec
mkdir -p ~/.config
mv /project/ldapsearcher/copr ~/.config/copr
copr-cli build ldapsearcher /root/rpmbuild/SRPMS/*.src.rpm
