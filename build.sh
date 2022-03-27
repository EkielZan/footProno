cd src
CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -v -ldflags="-X main.Version=Test" -o ../bin/footProno .
cd ../healthcheck
CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -v -o ../bin/healthcheck
cd ..
