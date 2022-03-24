FROM golang:1.13 AS builder
ARG VERSION=0.0.1
WORKDIR /go/src/gitlab.com/EkielZan/footProno
COPY src src
COPY healthcheck healthcheck
COPY go.mod src/
COPY go.mod healthcheck/
#RUN go install gitlab.com/EkielZan/footProno
RUN echo $VERSION && \
    cd src && \
    go install gitlab.com/EkielZan/footProno && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.Version=$VERSION" -o ../bin/footProno .
RUN cd healthcheck && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/healthcheck

FROM scratch
COPY --from=builder /go/src/gitlab.com/EkielZan/footProno/bin/* /    
ADD .env /
ADD static static
ADD ressources ressources
ADD templates templates
ADD certs certs
CMD ["/footProno"]
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 CMD [ "/healthcheck" ]