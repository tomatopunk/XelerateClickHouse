#!/bin/bash

# 获取当前时间戳
CURRENT_TIME=$(date +'%Y%m%d%H%M%S')

# 获取 Git 最后一个提交的哈希值
GIT_COMMIT=$(git rev-parse --short HEAD)

# 将当前时间和 Git 提交组合作为版本号
VERSION="${CURRENT_TIME}-${GIT_COMMIT}"

echo $VERSION
