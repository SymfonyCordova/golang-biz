#!/bin/bash
rm -rf golang-biz -f
rm -rf *.db -f
# rm -rf *.bat -f
go build
./golang-biz
