#!/bin/bash

# Promotion Service 启动脚本
# 作者: AI Assistant
# 描述: 启动 promotion-service 微服务

set -e  # 遇到错误时退出

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
LOG_DIR="$SCRIPT_DIR/logs"
PID_FILE="$SCRIPT_DIR/${SERVICE_NAME}.pid"

# 设置环境变量 (从你的主服务配置中提取)
export JAEGER_ENDPOINT="http://jaeger.infra:14268/api/traces"
export KAFKA_BROKERS="kafka-service.infra:9092"
export NACOS_SERVER_ADDRS="nacos.infra:8848"
export DB_SOURCE="root:root@tcp(mysql.infra:3306)/test"
export REDIS_ADDR="redis.infra:6379"
export ZK_SERVERS="zookeeper-headless.infra:2181"
export REDIS_ADDRS="redis-cluster-0.redis-cluster-headless.infra:6379,redis-cluster-1.redis-cluster-headless.infra:6379,redis-cluster-2.redis-cluster-headless.infra:6379"
export NACOS_NAMESPACE="d586122c-170f-40e9-9d17-5cede728cd7e"
export NACOS_GROUP="nexus-group"

# 创建必要的目录
mkdir -p "$LOG_DIR"
mkdir -p "$SCRIPT_DIR/deploy"

echo -e "${BLUE}🚀 开始启动 $SERVICE_NAME...${NC}"

# 检查端口是否被占用
if lsof -Pi :$SERVICE_PORT -sTCP:LISTEN -t >/dev/null ; then
    echo -e "${YELLOW}⚠️  端口 $SERVICE_PORT 已被占用，尝试停止现有服务...${NC}"
    pid=$(lsof -Pi :$SERVICE_PORT -sTCP:LISTEN -t)
    kill -9 $pid
    echo -e "${GREEN}✅ 已停止占用端口的进程 (PID: $pid)${NC}"
    sleep 2
fi

# 清理旧的PID文件
rm -f "$PID_FILE"

# 检查并停止可能残留的旧进程
if [ -f "$PID_FILE" ]; then
    while IFS= read -r pid; do
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            echo -e "${YELLOW}🔧 停止残留的服务进程 (PID: $pid)...${NC}"
            kill -9 "$pid"
            sleep 1
        fi
    done < "$PID_FILE"
    rm -f "$PID_FILE"
fi

# 检查并杀死可能残留的旧进程 (按二进制文件路径)
if pgrep -f "$BINARY_PATH" > /dev/null; then
    pkill -f "$BINARY_PATH"
    echo -e "${GREEN}✅ 已停止残留的 $SERVICE_NAME 服务${NC}"
    sleep 1
fi

# 编译服务
echo -e "${BLUE}🔧 编译 $SERVICE_NAME...${NC}"
SERVICE_PATH="$SCRIPT_DIR/cmd"

if [ ! -d "$SERVICE_PATH" ]; then
    echo -e "${RED}❌ 服务目录不存在: $SERVICE_PATH${NC}"
    echo -e "${YELLOW}💡 请确保你在 promotion-service 项目根目录下运行此脚本${NC}"
    exit 1
fi

# 编译
(cd "$SERVICE_PATH" && go build -o "$BINARY_PATH")
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 编译失败: $SERVICE_NAME${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 编译成功${NC}"

# 启动服务
echo -e "${BLUE}🔧 启动 $SERVICE_NAME (端口: $SERVICE_PORT)...${NC}"
"$BINARY_PATH" > "$LOG_DIR/$SERVICE_NAME.log" 2>&1 &
SERVICE_PID=$!
echo $SERVICE_PID > "$PID_FILE"

# 等待服务启动
echo -e "${YELLOW}⏳ 等待服务启动...${NC}"
sleep 3

# 检查服务是否成功启动
if kill -0 $SERVICE_PID 2>/dev/null; then
    echo -e "${GREEN}✅ $SERVICE_NAME 已成功启动 (PID: $SERVICE_PID)${NC}"

    # 健康检查
    echo -e "${BLUE}🩺 执行健康检查...${NC}"
    for i in {1..10}; do
        if curl -s http://localhost:$SERVICE_PORT/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ 健康检查通过${NC}"
            break
        elif [ $i -eq 10 ]; then
            echo -e "${YELLOW}⚠️  健康检查失败，但服务进程正在运行${NC}"
        else
            echo -e "${YELLOW}⏳ 健康检查中... ($i/10)${NC}"
            sleep 2
        fi
    done
else
    echo -e "${RED}❌ $SERVICE_NAME 启动失败${NC}"
    echo -e "${YELLOW}📋 查看日志: tail -f $LOG_DIR/$SERVICE_NAME.log${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 $SERVICE_NAME 启动完成！${NC}"
echo -e "${BLUE}📋 服务信息:${NC}"
echo -e "  - 服务名称: $SERVICE_NAME"
echo -e "  - 服务端口: $SERVICE_PORT"
echo -e "  - 进程ID: $SERVICE_PID"
echo -e "  - 日志文件: $LOG_DIR/$SERVICE_NAME.log"
echo ""
echo -e "${YELLOW}💡 API 测试示例:${NC}"
echo -e "  ${GREEN}健康检查:${NC}"
echo -e "  curl http://localhost:$SERVICE_PORT/health"
echo -e "  ${GREEN}创建促销模板:${NC}"
echo -e "  curl -X POST http://localhost:$SERVICE_PORT/templates \\"
echo -e "    -H 'Content-Type: application/json' \\"
echo -e "    -d '{\"name\":\"测试满减券\", \"promotion_type\":\"STORE_COUPON\", \"discount_type\":\"FIXED_AMOUNT\"}'"
echo -e "  ${GREEN}计算最优优惠:${NC}"
echo -e "  curl -X POST http://localhost:$SERVICE_PORT/offers/calculate-best \\"
echo -e "    -H 'Content-Type: application/json' \\"
echo -e "    -d '{\"user\":{\"id\":123,\"is_vip\":true}, \"total_amount\":10000}'"
echo ""
echo -e "${BLUE}📁 日志文件位置: $LOG_DIR/$SERVICE_NAME.log${NC}"
echo -e "${BLUE}🛑 停止服务请运行: ./stop-promotion-service.sh${NC}"