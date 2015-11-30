#!/bin/bash -e
# Usage: Setup a minimal PXE Server
# ./pxe.sh IP DHCP_RANGE SSH_KEY
# ./pxe.sh "192.168.32.10" "192.168.32.2,192.168.32.254,12h" "AABC.... name"

PXE_SERVER_IP=$1
DHCP_RANGE=$2
SSH_AUTHORIZED_KEYS=$3

# dnsmasq - your all in one DHCP, TFTP, and DNS
dnf install -yq dnsmasq

cp /etc/dnsmasq.conf /etc/dnsmasq.old
cat << EOF > "/etc/dnsmasq.conf"
dhcp-range=$DHCP_RANGE
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
  append initrd=coreos_production_pxe_image.cpio.gz cloud-config-url=http://$PXE_SERVER_IP/cloud-config.yml
EOF

# TFTP ldlinux.c32 pxelinux.0
dnf install -yq syslinux
ln -s /usr/share/syslinux/pxelinux.0 /var/lib/tftpboot/pxelinux.0
ln -s /usr/share/syslinux/ldlinux.c32 /var/lib/tftpboot/ldlinux.c32

# TFTP kernel image and init RAM disk
dnf install -yq wget
wget -q -O /var/lib/tftpboot/coreos_production_pxe.vmlinuz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz
wget -q -O /var/lib/tftpboot/coreos_production_pxe_image.cpio.gz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz
# Add cobbler_var_lib_t and tftpdir_rw_t SELinux context as appropriate
restorecon -R /var/lib/tftpboot

systemctl enable dnsmasq
systemctl start dnsmasq

# HTTP

# static cloud-config HTTP server
dnf install -yq httpd

cat << EOF > "/var/www/html/cloud-config.yml"
#cloud-config
coreos:
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
ssh_authorized_keys:
    - ssh-rsa $SSH_AUTHORIZED_KEYS
EOF

systemctl enable httpd
systemctl start httpd

echo "Done"
