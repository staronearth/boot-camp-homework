# Build the manager binary
FROM golang:1.21 as builder

WORKDIR /webook


ENV GOPROXY https://goproxy.cn
COPY go.mod go.mod
COPY webook/ webook/

RUN go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -tags=k8s -o binwebook webook/main.go
RUN pwd&&ls -al

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:3.9
WORKDIR /webook
RUN ls -al
COPY --from=builder /webook/binwebook .
EXPOSE 8080
ENTRYPOINT ["./binwebook"]