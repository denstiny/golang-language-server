######################################################################
# @author      : caohuaming.evpn (2254228017@qq.com)
# @file        : makefile
# @created     : Friday May 02, 2025 14:52:27 CST
######################################################################

build:
	go build github.com/denstiny/golang-language-server

help:
	go run github.com/denstiny/golang-language-server -help

run:
	go run github.com/denstiny/golang-language-server
