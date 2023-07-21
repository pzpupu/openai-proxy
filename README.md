# Open AI 网关代理（用于记录更详细的数据使用情况）

## 编译

### 代理程序

在项目目录下执行 `make build-(linux|darwin|windows)` 即可编译对应平台的代理程序。
在 build 目录中查看已编译的程序。

### token 生成程序

在项目目录下执行 `make build-jwt-(linux|darwin|windows)` 即可编译 token 生成程序。
在 build 目录中查看已编译的程序。

## 运行

### 代理程序

复制 .env.proxy 为 .evn 并修改其中的配置。与docker-compose.yml 放置同级目标下。

> OPEN_API_KEY 为 Open AI 的 API KEY
> SECRET_KEY 为 JWT 的密钥，与`token 生成程序`中的 SECRET_KEY 字段保持一致。

docker compose up -d 即可启动代理程序。

默认端口为 8080

### token 生成程序

复制 .env.jwt 为 .evn 并修改其中的配置。与编译的 jwt 程序放置同级目标下。

> SECRET_KEY 字段与`代理程序`中的 SECRET_KEY 字段保持一致。
> 否则使用代理程序会出现401响应码，认证未通过。

通过 `./jwt-xxx` 执行程序，并输入用户名，即可生成 token。

> 用户名为记录用户使用情况的唯一标识，可自行定义。
