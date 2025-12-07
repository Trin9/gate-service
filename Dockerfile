# Stage 1: The Builder Stage
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum，并下载依赖，这是为了利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制整个项目代码
COPY . .

RUN ls -l /app
# 编译 Go 应用
# CGO_ENABLED=0：禁用 CGO，确保生成的是静态链接的二进制文件
# -a：强制重新构建依赖（确保干净构建）
# -ldflags "-s -w": 移除符号表和调试信息，进一步减小二进制文件体积
# -o proxy：将输出文件命名为 'proxy'
# ./cmd/server：假设这是你的 main 包路径
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o gate-service ./cmd

# Stage 2: The Final (Runtime) Stage
# 推荐使用 scratch 基础镜像，它几乎是空的，体积最小（约 5MB）
# FROM scratch
# 如果你的应用需要 SSL 证书或其他系统依赖，请使用 alpine
FROM alpine:latest

# 创建一个非 root 用户以增强安全性（可选但推荐）
RUN adduser -D proxyuser
USER proxyuser

# 设置工作目录
WORKDIR /home/proxyuser/app

# 从 builder 阶段复制编译好的二进制文件
# 假设你在 builder 阶段将二进制文件命名为 'proxy'
COPY --from=builder /app/gate-service .

# 如果你的应用需要配置文件，也请复制过来
# COPY config/config.yaml .

# 暴露服务端口 (根据你的 Go 应用实际监听的端口调整)
EXPOSE 8080

# 定义容器启动时执行的命令
CMD ["./gate-service"]
