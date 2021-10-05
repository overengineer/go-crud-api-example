#!/bin/bash

pkill go
netstat -tulpn 2>/dev/null | grep 8088 | awk '{print $NF}' | grep -Eo '[0-9]+' | xargs kill -9 2>/dev/null

export DBUSER="admin"
export DBPASS="admin"

sleep 1
./test_app.py