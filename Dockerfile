FROM golang:1.13 AS builder
ARG VERSION=0.0.1
WORKDIR /go/src/gitlab.com/EkielZan/footProno
COPY src/*.go .
COPY go.mod .
RUN go install gitlab.com/EkielZan/footProno
RUN echo $VERSION && \
    CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags="-X main.Version=$VERSION" -o footProno .

FROM scratch
COPY --from=builder /go/src/gitlab.com/EkielZan/footProno/footProno /    
ADD .env /
ADD static static
ADD ressources ressources
ADD certs certs
CMD ["/footProno"]