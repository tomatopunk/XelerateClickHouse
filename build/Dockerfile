# 使用 golang 作为基础镜像
FROM golang:1.19 AS builder

# 设置工作目录
WORKDIR /app

# 复制 Go 源代码到容器中
COPY . .

# 编译 Go 应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# 使用 scratch 作为最终的基础镜像
FROM scratch

# 复制编译好的应用程序到最终镜像中
COPY --from=builder /app/app /app

# 设置容器启动时执行的命令
ENTRYPOINT ["/app"]
