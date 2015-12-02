#!/bin/bash -e
# Usage: Setup a Pixiecore Server

# ./pxe.sh IP SSH_KEY
# ./pixiecore.sh "192.168.33.10" "AABC.... name"

PIXIECORE_SERVER_IP=$1
SSH_AUTHORIZED_KEYS=$2

# Pixiecore kernel and init RAM disk
dnf install -yq wget
mkdir -p /var/lib/image
wget -q -O /var/lib/image/coreos_production_pxe.vmlinuz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz
wget -q -O /var/lib/image/coreos_production_pxe_image.cpio.gz http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz
chcon -Rt svirt_sandbox_file_t /var/lib/image

# Docker
dnf install -yq docker
systemctl enable docker.service
systemctl start docker.service

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

# Pixiecore
docker pull danderson/pixiecore
#docker run -v /var/lib/image:/image --net=host danderson/pixiecore -kernel /image/coreos_production_pxe.vmlinuz -initrd /image/coreos_production_pxe_image.cpio.gz --cmdline cloud-config-url=http://$PIXIECORE_SERVER_IP/cloud-config.yml

cat << EOF > /etc/systemd/system/pixiecore.service
[Unit]
Description=Pixicore Service

[Service]
Type=simple
ExecStart=/usr/bin/docker run -v /var/lib/image:/image --net=host danderson/pixiecore -kernel /image/coreos_production_pxe.vmlinuz -initrd /image/coreos_production_pxe_image.cpio.gz --cmdline cloud-config-url=http://$PIXIECORE_SERVER_IP/cloud-config.yml

[Install]
WantedBy=multi-user.target
EOF

systemctl enable pixiecore.service
systemctl start pixiecore.service

echo "Done"
