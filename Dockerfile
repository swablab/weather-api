FROM golang:1.16
 
RUN mkdir -p /weather-api
 
WORKDIR /weather-api
 
ADD . /weather-api
 
RUN go build ./main.go

EXPOSE 10000

CMD ["./main"]