cd src
CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags="-X main.Version=$VERSION" -o ../bin/footProno .
cd ../healthcheck
CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o ../bin/healthcheck
cd ..
