# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.23

# 设置工作目录
WORKDIR /app

# 设置 GOPROXY 为国内的代理，加速依赖下载
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制二进制文件
COPY  bigagent .
COPY  config.yml .

RUN chmod +x bigagent

# 设置时区
ENV TZ=Asia/Shanghai

#备注端口暴露
EXPOSE 8010
EXPOSE 5678

# 设置容器启动时执行的命令
CMD ["./bigagent -s start"]