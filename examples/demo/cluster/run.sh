#!/bin/bash

# Pitaya 服务控制脚本
# 用法: ./run.sh [1|2]
# 1 - 启动前端服务
# 2 - 启动后端服务

# 参数校验
if [ $# -ne 1 ]; then
  echo "错误: 需要 1 个参数！"
  echo "用法: $0 [1|2]"
  echo "  1 - 前端服务"
  echo "  2 - 后端服务"
  exit 1
fi

# 定义服务命令
FRONTEND_CMD="go run main.go -frontend=true -type=frontend -port=3250"
BACKEND_CMD="go run main.go -frontend=false -type=backend -port=3351"


# 根据参数执行对应服务
case $1 in
  1)
    echo "🚀 启动前端服务..."
    eval $FRONTEND_CMD
    ;;
  2)
    echo "🔧 启动后端服务..."
    eval $BACKEND_CMD
    ;;
  *)
    echo "错误: 无效参数 '$1'"
    echo "可用参数: 1 (前端) 或 2 (后端)"
    exit 2
    ;;
esac