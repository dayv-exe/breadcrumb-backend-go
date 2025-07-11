FROM public.ecr.aws/amazonlinux/amazonlinux:2023 AS build
RUN yum install -y golang git zip && yum clean all
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64 \
    BIN_ROOT=/app/bin

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p ${BIN_ROOT}
ARG BUILD_DIRS
RUN set -e; \
  for dir in $(echo ${BUILD_DIRS} | tr ',' '\n'); do \
    echo "==> Building $dir"; \
    mkdir -p ${BIN_ROOT}/$dir; \
    echo "Running: go build -o ${BIN_ROOT}/$dir/bootstrap ./cmd/$dir/main.go"; \
    go build -o ${BIN_ROOT}/$dir/bootstrap ./cmd/$dir/main.go; \
    echo "Zipping binary..."; \
    cd ${BIN_ROOT}/$dir; \
    zip $dir.zip bootstrap; \
    cd -; \
  done

RUN ls -R ${BIN_ROOT}
