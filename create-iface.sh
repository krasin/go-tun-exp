#!/bin/sh

set -ue

sudo openvpn --mktun --dev tun2
sudo ip link set tun2 up
sudo ip addr add 10.0.0.1/24 dev tun2

