#!/bin/bash -e
# Usage: Setup an iPXE server


IPXE_SERVER_IP=$1
DHCP_RANGE=$2
SSH_AUTHORIZED_KEYS=$3

# Sanity
dnf install -yq vim

# dnsmasq - your all in one TFTP
dnf install -yq dnsmasq

cp /etc/dnsmasq.conf /etc/dnsmasq.old
cat << EOF > "/etc/dnsmasq.conf"
dhcp-range=$DHCP_RANGE
dhcp-authoritative
enable-tftp
tftp-root=/var/lib/tftpboot
# set tag "ipxe" if request comes from iPXE ("iPXE" user class)
dhcp-userclass=set:ipxe,iPXE
# if PXE request came from regular firmware, TFTP serve iPXE firmware
dhcp-boot=tag:!ipxe,undionly.kpxe
# if PXE request comes from iPXE, HTTP serve an iPXE boot script
dhcp-boot=tag:ipxe,http://$IPXE_SERVER_IP/boot.ipxe
log-queries
log-dhcp
conf-dir=/etc/dnsmasq.d,.rpmnew,.rpmsave,.rpmorig
EOF

# Create TFTP root directory
if [ ! -d "/var/lib/tftpboot" ]; then
  mkdir -p "/var/lib/tftpboot"
fi

# TFTP undionly.kpxe
dnf install -yq wget
wget -q -O /var/lib/tftpboot/undionly.kpxe http://boot.ipxe.org/undionly.kpxe
restorecon -R /var/lib/tftpboot

systemctl enable dnsmasq
systemctl start dnsmasq

# HTTP hosted kernel, initramfs, cloud-config
dnf install -yq httpd

cat << EOF > "/var/www/html/boot.ipxe"
#!ipxe
chain http://$IPXE_SERVER_IP/ipxe?mac=\${net0/mac}&uuid=\${uuid}
EOF

cat << EOF > "/var/www/html/ipxe"
#!ipxe
set base-url http://stable.release.core-os.net/amd64-usr/current
kernel http://$IPXE_SERVER_IP/coreos_production_pxe.vmlinuz cloud-config-url=http://$IPXE_SERVER_IP/cloud-config.yml
initrd http://$IPXE_SERVER_IP/coreos_production_pxe_image.cpio.gz
boot
EOF

# Kernel image and initramfs over HTTP
wget -q -O /var/www/html/coreos_production_pxe.vmlinuz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz
wget -q -O /var/www/html/coreos_production_pxe_image.cpio.gz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz

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
