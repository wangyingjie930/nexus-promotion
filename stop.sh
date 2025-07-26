#!/bin/bash

# Promotion Service 停止脚本
# 作者: AI Assistant
# 描述: 停止 promotion-service 微服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 服务配置
SERVICE_NAME="promotion-service"
SERVICE_PORT="8087"
BINARY_PATH="$SCRIPT_DIR/deploy/${SERVICE_NAME}"
PID_FILE="$SCRIPT_DIR/${SERVICE_NAME}.pid"

echo -e "${BLUE}🛑 开始停止 $SERVICE_NAME...${NC}"

# 检查PID文件是否存在
if [ ! -f "$PID_FILE" ]; then
    echo -e "${YELLOW}⚠️  没有找到PID文件，尝试通过其他方式查找进程...${NC}"
else
    # 读取PID并停止进程
    while IFS= read -r pid; do
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            echo -e "${BLUE}🔧 停止进程 (PID: $pid)...${NC}"
            kill -TERM "$pid"

            # 等待进程优雅关闭
            for i in {1..10}; do
                if ! kill -0 "$pid" 2>/dev/null; then
                    echo -e "${GREEN}✅ 进程已优雅停止 (PID: $pid)${NC}"
                    break
                fi
                sleep 1
            done

            # 如果进程仍然存在，强制杀死
            if kill -0 "$pid" 2>/dev/null; then
                echo -e "${YELLOW}⚠️  强制停止进程 (PID: $pid)...${NC}"
                kill -9 "$pid"
                sleep 1
                if ! kill -0 "$pid" 2>/dev/null; then
                    echo -e "${GREEN}✅ 进程已强制停止 (PID: $pid)${NC}"
                else
                    echo -e "${RED}❌ 无法停止进程 (PID: $pid)${NC}"
                fi
            fi
        else
            echo -e "${YELLOW}⚠️  进程不存在或已停止 (PID: $pid)${NC}"
        fi
    done < "$PID_FILE"

    # 删除PID文件
    rm -f "$PID_FILE"
fi

# 通过二进制文件路径查找并停止进程
if pgrep -f "$BINARY_PATH" > /dev/null; then
    echo -e "${YELLOW}🔍 发现通过二进制路径运行的进程，正在停止...${NC}"
    pkill -TERM -f "$BINARY_PATH"
    sleep 2

    # 检查是否还有残留进程
    if pgrep -f "$BINARY_PATH" > /dev/null; then
        echo -e "${YELLOW}⚠️  强制停止残留进程...${NC}"
        pkill -9 -f "$BINARY_PATH"
    fi
    echo -e "${GREEN}✅ 已停止通过二进制路径运行的 $SERVICE_NAME 进程${NC}"
fi

# 通过端口查找并停止进程
port_pid=$(lsof -ti tcp:$SERVICE_PORT 2>/dev/null)
if [ -n "$port_pid" ]; then
    echo -e "${YELLOW}🔍 发现占用端口 $SERVICE_PORT 的进程 (PID: $port_pid)，正在停止...${NC}"
    kill -TERM $port_pid
    sleep 2

    # 检查进程是否还在运行
    if kill -0 $port_pid 2>/dev/null; then
        echo -e "${YELLOW}⚠️  强制停止端口进程 (PID: $port_pid)...${NC}"
        kill -9 $port_pid
    fi
    echo -e "${GREEN}✅ 端口 $SERVICE_PORT 已释放${NC}"
fi

# 最终检查
sleep 1
if pgrep -f "$SERVICE_NAME" > /dev/null || lsof -ti tcp:$SERVICE_PORT >/dev/null 2>&1; then
    echo -e "${RED}❌ 仍有 $SERVICE_NAME 相关进程在运行${NC}"
    echo -e "${YELLOW}📋 相关进程:${NC}"
    pgrep -af "$SERVICE_NAME" || true
    lsof -i tcp:$SERVICE_PORT || true
else
    echo -e "${GREEN}🎉 $SERVICE_NAME 已完全停止！${NC}"
fi

# 清理日志文件（可选）
rm -rf "$SCRIPT_DIR/logs"
echo -e "${GREEN}✅ 日志文件已删除${NC}"

echo -e "${BLUE}✨ $SERVICE_NAME 停止脚本执行完成${NC}"