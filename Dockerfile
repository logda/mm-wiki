
FROM golang:1.14.1-alpine

WORKDIR /app/mm-wiki

COPY . /app/mm-wiki

# 如果国内网络不好，可添加以下环境
# RUN go env -w GO111MODULE=on
# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN export GO111MODULE=on
# RUN export GOPROXY=https://goproxy.cn

# Set environment variables
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies
RUN go mod download

RUN mkdir -p /opt/mm-wiki && ls /app/mm-wiki

RUN go build -o /opt/mm-wiki/mm-wiki ./ \
    && cp -r ./conf/ /opt/mm-wiki \
    && cp -r ./install/ /opt/mm-wiki \
    && cp ./scripts/run.sh /opt/mm-wiki \
    && cp -r ./static/ /opt/mm-wiki \
    && cp -r ./views/ /opt/mm-wiki \
    && cp -r ./logs/ /opt/mm-wiki \
    && cp -r ./docs/ /opt/mm-wiki

CMD ["/opt/mm-wiki/mm-wiki", "--conf", "/opt/mm-wiki/conf/mm-wiki.conf"]
