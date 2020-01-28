#!/bin/bash
set -e

echo "== Install dependecies =="
go mod download
echo ""

echo "== Build =="
make -B
echo ""
