# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.23

# 设置工作目录
WORKDIR /app

# 设置 GOPROXY 为国内的代理，加速依赖下载
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./
# 复制整个项目代码
COPY  . .

# 下载依赖
RUN go mod download

# 构建 Go 应用程序
RUN go build -o bigagnt  cmd/agent/main.go

# 安装必要的工具和依赖
RUN apt-get update && apt-get install -y ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
EXPOSE 8010
EXPOSE 5678
# 设置容器启动时执行的命令
CMD ["./bigagnt"]