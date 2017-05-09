%global import_path github.com/coreos/matchbox
%global repo matchbox
%global debug_package %{nil}

Name:		matchbox
Version:	0.6.0
Release:	2%{?dist}
Summary:	Network boot and provision CoreOS machines
License:	ASL 2.0
URL:		https://%{import_path}
Source0:        https://%{import_path}/archive/v%{version}/%{name}-%{version}.tar.gz


BuildRequires: golang
BuildRequires: systemd
%{?systemd_requires}

Requires(pre): shadow-utils

%description
matchbox is a service that matches machines to profiles to PXE boot and provision
clusters. Profiles specify the kernel/initrd, kernel args, iPXE config, GRUB
config, Container Linux config, Cloud-config, or other configs. matchbox provides
a read-only HTTP API for machines and an authenticated gRPC API for clients.

# Limit to architectures supported by golang or gcc-go compilers
ExclusiveArch: %{go_arches}
# Use golang or gcc-go compiler depending on architecture
BuildRequires: compiler(golang)

%prep
%setup -q -n %{repo}-%{version}

%build
# create a Go workspace with a symlink to builddir source
mkdir -p src/github.com/coreos
ln -s ../../../ src/github.com/coreos/matchbox
export GOPATH=$(pwd):%{gopath}
export GO15VENDOREXPERIMENT=1
function gobuild { go build -a -ldflags "-w -X github.com/coreos/matchbox/matchbox/version.Version=v%{version}" "$@"; }
gobuild -o bin/matchbox %{import_path}/cmd/matchbox

%install
install -d %{buildroot}/%{_bindir}
install -d %{buildroot}%{_sharedstatedir}/%{name}
install -p -m 0755 bin/matchbox %{buildroot}/%{_bindir}
# systemd service unit
mkdir -p %{buildroot}%{_unitdir}
cp contrib/systemd/%{name}.service %{buildroot}%{_unitdir}/

%files
%doc README.md CHANGES.md CONTRIBUTING.md LICENSE NOTICE DCO
%{_bindir}/matchbox
%{_sharedstatedir}/%{name}
%{_unitdir}/%{name}.service

%pre
getent group matchbox >/dev/null || groupadd -r matchbox
getent passwd matchbox >/dev/null || \
    useradd -r -g matchbox -s /sbin/nologin matchbox

%post
%systemd_post matchbox.service

%preun
%systemd_preun matchbox.service

%postun
%systemd_postun_with_restart matchbox.service

%changelog
* Mon Apr 24 2017 <dalton.hubble@coreos.com> - 0.6.0-1
- New support for terraform-provider-matchbox plugin
- Add ProfileDelete, GroupDelete, IgnitionGet and IgnitionDelete gRPC endpoints
- Generate code with gRPC v1.2.1 and matching Go protoc-gen-go plugin
- Update Ignition to v0.14.0 and coreos-cloudinit to v1.13.0
- New documentation at https://coreos.com/matchbox/docs/latest
* Wed Jan 25 2017 <dalton.hubble@coreos.com> - 0.5.0-1
- Rename project from bootcfg to matchbox
* Sat Dec 3 2016 <dalton.hubble@coreos.com> - 0.4.1-3
- Add missing ldflags which caused bootcfg -version to report wrong version
* Fri Dec 2 2016 <dalton.hubble@coreos.com> - 0.4.1-2
- Fix bootcfg user creation
* Fri Dec 2 2016 <dalton.hubble@coreos.com> - 0.4.1-1
- Initial package

