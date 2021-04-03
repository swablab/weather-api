############################
# STEP 1 build executable binary
############################
FROM golang:1.16 AS builder 
RUN mkdir -p /weather-api
WORKDIR /weather-api
ADD . /weather-api
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "$(govvv -flags)" -o app

# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /weather-api/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 10000
ENTRYPOINT ["/app"]