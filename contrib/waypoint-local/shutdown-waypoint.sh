#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


echo "==> Attempting to gracefully shutdown Waypoint server and runner..."
echo
echo "==> Shutting down waypoint server"
echo

pkill --signal SIGINT -f "waypoint server run"

echo
echo "==> Shutting down waypoint runner"
echo

pkill --signal SIGINT -f "waypoint runner agent"

echo
echo "Finished shutting down local waypoint server and runner!"
echo
