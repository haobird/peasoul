FROM registry.cn-hangzhou.aliyuncs.com/city-cloud/golang:1.15-1.0.1
LABEL maintainer="xyqbgn@126.com"

ENV HOME  "/www"
RUN mkdir ${HOME}
WORKDIR ${HOME}/

COPY . .

# RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server main.go &&\
#   upx --best server -o _upx_server && \
#   mv -f _upx_server server
RUN go build -ldflags "-s -w" -o server main.go

CMD ["./server"]