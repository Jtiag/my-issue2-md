#!/bin/bash
# 仅仅负责把过去 10 条 commit 提取出来，不让大模型去猜 Git 命令
git log -n 10 --pretty=format:"%h - %s (%an)"