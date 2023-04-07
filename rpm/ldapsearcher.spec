%global go_version 1.18.10
%global go_release go1.18.10

Name:           ldapsearcher
Version:        0.1.5
Release:        1%{?dist}
Summary:        This is a test
License:        GPLv3
URL:            https://github.com/spideyz0r/ldapsearcher
Source0:        %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= %{go_version}
BuildRequires:  git

%description
A ldap search tool. Run pre-defined or custom queries

%global debug_package %{nil}

%prep
%autosetup -n %{name}-%{version}

%build
go build -v -o %{name} -ldflags=-linkmode=external

%check
go test

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}


%files
%{_bindir}/ldapsearcher

%license LICENSE

%changelog
* Fri Apr  7 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.1.5-1
- Initial build

