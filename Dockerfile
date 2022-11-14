FROM golang:1.18-alpine AS builder
ARG VERSION=0.0.1
WORKDIR /go/src/gitlab.com/EkielZan/footProno
COPY build.sh build.sh
COPY src src
COPY go.mod go.mod
COPY go.sum go.sum
COPY healthcheck healthcheck
#RUN go install gitlab.com/EkielZan/footProno
RUN chmod +x ./build.sh && \
    apk add gcc musl-dev && \
    ls -l && \
    cd src && \
    echo "Build Main Binaries" && \
    CGO_ENABLED=1 GOOS=linux go build -installsuffix cgo -v -ldflags="-X main.Version=Test" -o ../bin/footProno . && \
    cd ../healthcheck && \
    echo "Build Healthcheck Binaries" && \
    CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -v -o ../bin/healthcheck && \
    cd .. && \
    ls -l bin && \
    ls -l /go/src/gitlab.com/EkielZan/footProno/bin/* 

FROM alpine
WORKDIR /
COPY --from=builder /go/src/gitlab.com/EkielZan/footProno/bin/footProno /footProno    
COPY --from=builder /go/src/gitlab.com/EkielZan/footProno/bin/healthcheck /healthcheck
COPY .env .env
COPY static static
COPY templates templates    
COPY DB DB
CMD ["/footProno"]
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 CMD [ "/healthcheck" ]