

# 准备
## 依赖
> mysql8.0
>
> redis7.x
>
> golang1.23.4
>

## 建库与导库
```sql
#建库
CREATE BATABASE bigagent;
USE  bigagent;

#导入server仓库中的sql
source  bigagent.sql
```

# agent端
> agent使用教程，请先部署并运行agent端
>

## Build
```shell
go mod tidy
```

```shell
go env -w GOOS=linux
go build -o "bigagent"
chmod 755 bigagent
```

配置文件修改

```yaml
system:
  #http端口
  addr: :8010
  #grpc客户端host
  grpc: 127.0.0.1
  #grpc客户端端口
  grpc_port: 5678
  #grpc服务端套接字
  grpc_server: 0.0.0.0:8765
  #日志文件路径
  logfile: log.txt
  #认证密钥， 前端、agent端、server端、配置文件都需要一致
  serct: "123456"
#如下配置为自动生成，请勿修改!!！
grpc_cmdb1_stand1: 192.168.11.11:5555
grpc_cmdb1_stand1_token: "123456"
grpc_cmdb2_stand1:
grpc_cmdb2_stand1_token:
grpc_cmdb3_stand1:
grpc_cmdb3_stand1_token:
grpc_cmdb1_stand2:
grpc_cmdb1_stand2_token:
grpc_cmdb2_stand2:
grpc_cmdb2_stand2_token:
grpc_cmdb3_stand2:
grpc_cmdb3_stand2_token:
grpc_cmdb1_stand3:
grpc_cmdb1_stand3_token:
grpc_cmdb2_stand3:
grpc_cmdb2_stand3_token:
grpc_cmdb3_stand3:
grpc_cmdb3_stand3_token:
action_detail: "26"
collection_frequency: 2s

```

## Run
```shell
nohup ./bigagent -s start > /dev/null 2>&1  &
```

```shell
nohup ./bigagent -c /path/config.yml > /dev/null 2>&1  &
```

# server端
> 部署bigagent-server端
>
> 地址：[https://gitee.com/yl166490/bigagent-server.git](https://gitee.com/yl166490/bigagent-server.git)
>

## Build
```shell
git clone https://gitee.com/yl166490/bigagent-server.git
cd bigagent-server
mkdir conf
mv bigagent-server/conf/config.yml ./
```

```shell
go mod tidy
go env -w GOOS=linux
go build -o "bigagent-server" cmd/server/main.go
chmod 755 bigagent-server
```

配置文件修改

```yaml
system:
  # http端口
  addr: ":8080"
  # grpc服务端端口
  grpc: "0.0.0.0:8765"
  # 日志文件
  logfile: "log.txt"
  # 认证密钥， 前端、agent端、server端、配置文件都需要一致
  serct: 123456
  # grpc客户端端地址
  client_port: "0.0.0.0:5678"
  # agent端口
  agent_port: ":8010"
  # 查询agent离线时间间隔
  times: "10s"
  # 判断agent离线超时时间
  agent_outtime: 20
  # mysql数据库
  database:
    mysqlhost: "localhost"
    mysqlport: "3306"
    mysqluser: "root"
    mysqlpassword: "123456"
    mysqldatabasename: "bigagent"
  # redis数据库
  redisaddr: "localhost:6379"
  redispassword: "166490"
  redisdb: 0
```

## Run
```shell
nohup ./bigagent-server > /dev/null 2>&1  &
```

## Docker
> 完成源码编译后
>

Dockerfile

```yaml
# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.23

# 设置工作目录
WORKDIR /app

# 设置 GOPROXY 为国内的代理，加速依赖下载
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制整个项目代码
COPY cmd /app/cmd
COPY docs /app/docs
COPY inits /app/inits
COPY internel /app/internel
COPY config.yml /app/config.yml

# 构建 Go 应用程序
RUN go build -o bigagnt-server cmd/server/main.go

# 安装必要的工具和依赖
RUN apt-get update && apt-get install -y ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
EXPOSE 8080
# 设置容器启动时执行的命令
CMD ["./bigagnt-server"]
```

```shell
docker build  -t bigagent-server .
```

```yaml
docker run -d \
  --name bigagentserver \
  -p 8080:8080 \
  -p 5678:5678 \
  -p 8765:8765 \
  -v conf:/app/conf/ \
  bigagent-server
```

# web端
> bigagent-server的前端部署
>
> 地址：[https://gitee.com/yl166490/bigagent-server.git](https://gitee.com/yl166490/bigagent-server.git)
>

## Build
修改前端配置文件

```shell
git clone https://gitee.com/yl166490/bigagent-server.git
cd bigagent-server/web/
vim .env.pro
```

```shell
# 环境
VITE_NODE_ENV=production

# 接口前缀，后端server端接口地址
VITE_API_BASE_PATH='http://192.168.0.83:8080'

# nginx的路径
VITE_BASE_PATH=/agent/

#TOKEN 前端、agent端、server端、配置文件都需要一致
VITE_TOKEN=123456

# 是否删除debugger
VITE_DROP_DEBUGGER=true

# 是否删除console.log
VITE_DROP_CONSOLE=true

# 是否sourcemap
VITE_SOURCEMAP=false

# 输出路径
VITE_OUT_DIR=dist-pro

# 标题
VITE_APP_TITLE=BGG

# 是否包分析
VITE_USE_BUNDLE_ANALYZER=true

# 是否全量引入element-plus样式
VITE_USE_ALL_ELEMENT_PLUS_STYLE=true

# 是否开启mock
VITE_USE_MOCK=false

# 是否切割css
VITE_USE_CSS_SPLIT=true

# 是否使用在线图标
VITE_USE_ONLINE_ICON=true
```

编译生成dist-pro文件

```shell
npm config set registry https://registry.npmmirror.com
npm install pnpm -g
pnpm i
npm run build:pro
cd dist-pro
```

## Apply
```shell
vim bigagent.conf
```

```nginx
server {
  listen 3000;
  server_name localhost;
  location /agent/ {
    alias /usr/share/nginx/html/web/;
    index index.html;
    try_files $uri $uri/ /index.html;
  }
}
```

## Docker
> 在完成源码配置与编译后进行
>

Dockerfile

```dockerfile
# 基础镜像
FROM nginx:latest
# 镜像维护
LABEL maintainer=yeling
# 将 dist 文件夹拷贝到 Nginx 的静态资源目录
COPY dist-pro/. /usr/share/nginx/html/web/
# 将前端 Nginx 配置覆盖基础镜像的配置文件（最好在启动容器的时候挂载到宿主机由运维维护）
COPY bigagent.conf /etc/nginx/conf.d/bigagent.conf
```

执行镜像构建与容器运行（默认上述）

```shell
docker  build -t bigagent-server-web .
docker run -d --name bigagent-web -p 3000:3000 bigagent-server-web
```

访问 

```http
http://localhost:3000/agent/
user：admin
password：admin
```



# Docker-compose
> 快速体验模式，一键部署server前后端以及依赖组件
>
> [@李泽建](undefined/lizejian)
>

```nginx
...待补充
```

# Contributing
<font style="color:rgb(31, 35, 40);">如果你有好的意见或建议，欢迎给我们提 Issues 或 Pull Requests</font>

# Partners
[@李泽建](undefined/lizejian)[@叶凌](undefined/yeling-cqybb)

