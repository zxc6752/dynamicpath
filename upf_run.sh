#!/bin/bash

exe() { echo "\$ $@" ; "$@" ; }

function show_usage {
    echo
    echo "RUN UPF"
    echo "Usage: $0 host eth [-v|--virtual=] [-r=] [-t=]"
    echo
    echo "Arguments:"
    echo "  -v|--virtual=<viface>  : the virtual interface for incoming (default: no virtual)"
    echo "  -r=<bandwidth>       : define the Tx total bandwidth (default: 100(mbps))"
    echo "  -t=<bandwidth>       : define the Rx total bandwidth (default: 100(mbps))"
}

if [ -z "$2" ]
then
    show_usage
    exit 1
fi

HOST=$1
INTERFACE=$2
TX="100"
RX="100"
PARAM=

shift 2

for i in "$@"
do
case $i in
    -v=*|--virtual=*)
    VIRTUAL="v=${i#*=},"
    shift
    ;;
    -r=*)
    RX="${i#*=}"
    shift
    ;;
    -t=*)
    TX="${i#*=}"
    shift
    ;;
esac
done

exe ./bin/monitor -h ${HOST} "-"${INTERFACE} ${VIRTUAL}rx=${RX},tx=${TX}