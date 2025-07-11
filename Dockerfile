FROM public.ecr.aws/amazonlinux/amazonlinux:2023 AS build
RUN yum install -y golang git zip && yum clean all
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64 \
    BIN_ROOT=/app/bin

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN mkdir -p ${BIN_ROOT}
COPY . .
ARG BUILD_DIRS
RUN for dir in $(echo ${BUILD_DIRS} | tr ',' '\n'); do \
    mkdir -p ${BIN_ROOT}/$dir && \
    go build -o ${BIN_ROOT}/$dir/bootstrap ./cmd/$dir/main.go && \
    cd ${BIN_ROOT}/$dir && \
    zip $dir.zip bootstrap; \
  done

RUN ls -R ${BIN_ROOT}
