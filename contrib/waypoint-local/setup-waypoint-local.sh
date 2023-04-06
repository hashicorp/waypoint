#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


LOGDIR="${WP_LOG_DIR:-"/tmp"}"
DBDIR="${WP_DB_DIR:-"."}"

echo
echo "==> Starting waypoint server"
echo
echo "database dir: ${DBDIR}"

waypoint server run -accept-tos -advertise-addr=127.0.0.1:9701 \
  -listen-grpc=0.0.0.0:9701 -listen-http=0.0.0.0:9702 -db=$DBDIR/data.db \
  -advertise-tls-skip-verify -url-enabled -vvv > $LOGDIR/wp-server-logs.txt 2>&1 &

echo
echo "==> Bootstrapping waypoint server"
echo
echo "Server bootstrap token will print to STDOUT"

waypoint server bootstrap -server-addr=127.0.0.1:9701 -server-tls-skip-verify

echo
echo "=>> Starting a waypoint runner"
echo

waypoint runner agent -vvv > $LOGDIR/wp-runner-logs.txt 2>&1 &

echo
echo "Finished setting up a local waypoint server and runner!"
echo 
echo "Database file saved at: ${DBDIR}/data.db"
echo
echo "Logs can be found at:"
echo "waypoint server: ${WP_LOG_DIR}/wp-server-logs.txt"
echo "waypoint runner: ${WP_LOG_DIR}/wp-runner-logs.txt"
