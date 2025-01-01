# 使用官方Go镜像作为构建环境
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 设置构建参数
ARG LDFLAGS="-s -w -extldflags '-fno-PIC -static' -buildid="
ARG GCFLAGS="-trimpath=/go -m -l=4"
ARG ASMFLAGS="-trimpath=/go"

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app \
    -ldflags="${LDFLAGS}" \
    -gcflags="${GCFLAGS}" \
    -asmflags="${ASMFLAGS}" \
    ./cmd

# 使用轻量级的alpine镜像作为运行环境
FROM alpine:latest  

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/app .

# 暴露50050端口
EXPOSE 50050
EXPOSE 50080

# 设置 ENTRYPOINT
ENTRYPOINT ["./app"]

# 设置默认 CMD
CMD ["server"]