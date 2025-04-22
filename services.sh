#!/bin/bash
# services.sh

# 定义颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 定义服务目录
services=(
    "ucenter-api"
    "ucenter"
    "market-api"
    "market"
    "jobcenter"
)

# 启动所有服务
start_all() {
    echo -e "${GREEN}Starting all services...${NC}"
    for service in "${services[@]}"; do
        echo -e "${GREEN}Starting $service...${NC}"
        cd $service
        go run . &
        cd ..
    done
    echo -e "${GREEN}All services started${NC}"
}

# 停止所有服务
stop_all() {
    echo -e "${YELLOW}Stopping all services...${NC}"
    pkill -f "go run ."
    echo -e "${YELLOW}All services stopped${NC}"
}

# 命令处理
case "$1" in
    start)
        start_all
        ;;
    stop)
        stop_all
        ;;
    restart)
        stop_all
        sleep 2
        start_all
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        ;;
esac