
# Vagrant Network CIDR
$network_range="192.168.32.0/24"

# PXE Server IP, must be from the network_range
$pxe_server_ip="192.168.32.10"

# DHCP range dnsmasq should serve, must be a subset of network_range
$dhcp_range="192.168.32.2,192.168.32.254,12h"

# SSH Authorized Key for client CoreOS instances
$ssh_authorized_key="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC9oRIjXKgC1It3U22INv9sDbQjzZNbY6fdzN28hl2gnWf7b4/KJjbCE8cldAiV6qiLwnaqnINgoAy8JN718qos8VsLRdB/GvhlVOQvjJf6gSI9WcG1kVbbYuZ7WV1cxnxjE21+oHHz4IZyGKP6rEv0ODcFWokJt13zpK9isG7iQyBi51KNFPgox/jfM0uDCf+yzSsCX2HUUxmqKDUXD9XDihrGRpbqL6gH5VDYzDmVAHq5e3er1Sz2n+Gx/wUSXzNk9TdCY/cS6k2C6H3+dwA45HFADjmeK+k3dE+cDrXkLsB9GTXnvcmtdoVAFoHBZo8GqRKocaejVgDaRo+prQyJ dghubble"