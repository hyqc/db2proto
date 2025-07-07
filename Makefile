# 版本信息
VERSION:=$(shell git describe --tags --always)

# proto文件
PROTO_FILES:=$(shell find proto/$(dir) -name *.proto)

# proto文件所在目录列表
PROTO_DIRS:=$(shell find proto/* -type d)

#根目录
ROOT=$(shell pwd)
#ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

## 可选 os，默认值 linux，可能值：linux,win
os ?= linux

# 根据os设置输出路径和可执行文件名称
ifeq ($(os),windows)
	DEPLOY_OUTPUT_DIR := .
	BINARY_FILE := $(DEPLOY_OUTPUT_DIR)/db2proto.exe
	os = windows
else
	DEPLOY_OUTPUT_DIR := .
    BINARY_FILE := $(DEPLOY_OUTPUT_DIR)/db2proto
endif

#列出命令并加以解释说明
# 颜色定义（可选）
GREEN  := $(shell echo -e "\033[32m")
RESET  := $(shell echo -e "\033[0m")


# 默认目标
.DEFAULT_GOAL := help

# 列出你的各个目标，并为其加上说明
.PHONY: help
help: ## 显示此帮助信息
	@echo "$(HELP_MSG)"
	@cat $(MAKEFILE_LIST) | sort | \
 	grep -E '^[a-zA-Z_-]+:.*?##.*$$' | \
	awk '{ gsub("##", "\n");print $0}' | \
	awk '{ gsub(/\\n/, "\n");print $0}'

.PHONY: init
# init env
init: ## 下载最新包并执行go mod tidy
	go mod tidy

## 构建可执行文件
.PHONY: build
build: ## 构建可执行文件 \n 	os: 系统变量值 \n	- windows \n 	- linux
	@mkdir -p $(DEPLOY_OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build -o $(BINARY_FILE) *.go






