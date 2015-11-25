#!/bin/bash -e
# Setup a minimal PXE Server

# PXE Server IP should be the static IP set in the Vagrantfile.
export NODE_IP=192.168.32.10

# dnsmasq - your all in one DHCP, TFTP, and DNS
dnf install -yq dnsmasq

cp /etc/dnsmasq.conf /etc/dnsmasq.old
cat << EOF > "/etc/dnsmasq.conf"
dhcp-range=192.168.32.2,192.168.32.254,12h
dhcp-boot=pxelinux.0
enable-tftp
tftp-root=/var/lib/tftpboot
dhcp-authoritative
log-queries
log-dhcp
conf-dir=/etc/dnsmasq.d,.rpmnew,.rpmsave,.rpmorig
EOF

# TFTP

# Create TFTP root directory
if [ ! -d "/var/lib/tftpboot/pxelinux.cfg" ]; then
	mkdir -p "/var/lib/tftpboot/pxelinux.cfg"
fi

# TFTP pxelinux.cfg
cat << EOF > "/var/lib/tftpboot/pxelinux.cfg/default"
default coreos
prompt 1
timeout 15

display boot.msg

label coreos
  menu default
  kernel coreos_production_pxe.vmlinuz
  append initrd=coreos_production_pxe_image.cpio.gz cloud-config-url=http://$NODE_IP/pxe-cloud-config.yml
EOF

# TFTP ldlinux.c32 pxelinux.0
dnf install -yq syslinux
cp /usr/share/syslinux/pxelinux.0 /var/lib/tftpboot/pxelinux.0
cp /usr/share/syslinux/ldlinux.c32 /var/lib/tftpboot/ldlinux.c32

# TFTP kernel image and options
dnf install -yq wget
wget -q -O /var/lib/tftpboot/coreos_production_pxe.vmlinuz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz
wget -q -O /var/lib/tftpboot/coreos_production_pxe_image.cpio.gz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz
restorecon -R /var/lib/tftpboot

systemctl enable dnsmasq
systemctl start dnsmasq

# HTTP

# static HTTP server
dnf install -yq httpd

# TODO - this static config is exactly what we can improve upon
cat << EOF > "/var/www/html/pxe-cloud-config.yml"
#cloud-config
coreos:
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
ssh_authorized_keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDp9dypy87zDQgoMm9tk9aKVNDXemN2bp+taQzY4I2BS74H850SjvxVqMHbTcZVi+gMmUev9DPo5xfDOeTPeaT5fnNaKTLNGLezGM0pjMtKWOrZSQqK2YywjT/VdnBvNG6zKHE/MgpexLdkqiPJGPSHhRB+rQ5JX1FL3Qpt9prD2tW8audLmZ6PorC8/HLIwgBIrlh8JFbVFYGEy4dFImz2eyTqNUxDhP2hx17Drvssy/D6HR4zehgP1MGJWOKGfFlE2xe+hy+S6fjkPK5Nc7MR0uqhCkJHpCfVFjbHyZsQK3oi30Afe9t/CdRfPgQW+daeA7yBd1d/56A8QEbRf7W8wO5/CaSBoomA3bufJwwtT50oQwi3r0+hapivvmlUo3WhXwFoPV/XLyQUNGMWKmOWEqqmREFpi/NVnfItU8CenowvDoly0FkreztO4bAAhMd0w4/bAUisRHlucXXifqFEZr5NUt0BIaZmjavbXDadajzhwbK5mWpijuLoiuYCsa+5pQjFVlZRo5f3F3dQ5nREgRSWJPkxzllHygdEeKqxcE4slZJwZbrBgV4OT2nwGjBpfHsbBKBZRM6eKOiDIOVf/hh4+9+GirHWKosv3HRKHH9TisfIFlZoxzxY6E87oPvlf+7rEwHhuaPUoThSXxB6JS9LiOcqxjgL1ltfapgLqQ== dghubble@gmail.com
EOF

systemctl enable httpd
systemctl start httpd

echo "Done"