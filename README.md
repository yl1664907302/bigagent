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

```powershell
go env -w GOOS=windows
go build -o "bigagent"
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

## Build
```shell
go mod tidy
```

```shell
go env -w GOOS=linux
go build -o "bigagent-server"
chmod 755 bigagent-server
```

```shell
go env -w GOOS=windows
go build -o "bigagent-server"
```

## Run
```shell
nohup ./bigagent-server > /dev/null 2>&1  &
```



# web端
> bigagent-server的前端部署
>

## Build
```shell
#前提，完成前端prod配置文件配置
#dist文件就编译后的产物
npm install pnpm -g
pnpm i
cd bigagent-server/web/
pnpm build:prod
```

## Apply
```nginx
# /etc/nginx/conf.d/bigagent-server.conf
server {
  listen 80;
  server_name your_domain.com;  # 替换为你的域名

  location / {
    root /path/to/your/dist;  # 替换为前端文件的路径
    index index.html index.htm;
    try_files $uri $uri/ /index.html;  # 支持前端路由
  }

  location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg)$ {
    expires 30d;  # 缓存静态文件
    add_header Cache-Control "public";
  }

  location = /favicon.ico {
    log_not_found off;  # 不记录 favicon 的 404 错误
  }

  error_page 404 /404.html;  # 自定义 404 页面
  location = /404.html {
    internal;
  }
}
```

## Contributing
<font style="color:rgb(31, 35, 40);">如果你有好的意见或建议，欢迎给我们提 Issues 或 Pull Requests</font>

## Partners
[@李泽建](undefined/lizejian)[@叶凌](undefined/yeling-cqybb)

