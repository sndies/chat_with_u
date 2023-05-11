#!/usr/bin/env bash
sysctl net.ipv4.tcp_keepalive_time
sysctl -w net.ipv4.tcp_keepalive_time=200
sysctl -w net.ipv4.tcp_keepalive_intvl=200
sysctl -w net.ipv4.tcp_keepalive_probes=5
/bin/bash ./main