#!/usr/bin/env bash

echo "Start build Load Balancer...."
go build -o bin/balancer -x src/load_balancer/balancer.go

echo "Start build UPF Monitor...."
go build -o bin/monitor -x src/monitor/monitor.go