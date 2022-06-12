FROM golang:alpine as build-env
WORKDIR /src

ADD . /src
RUN go mod tidy

RUN go build -o rate-limit
FROM alpine
WORKDIR /app 
COPY --from=build-env /src/rate-limit /app/

CMD ./rate-limit