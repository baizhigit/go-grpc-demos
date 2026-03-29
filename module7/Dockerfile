FROM golang:1.26-alpine as build

ENV GONOSUMDB "github.com/baizhigit"
ENV GOPRIVATE "github.com/baizhigit"

WORKDIR /build

COPY . .

RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -gcflags='-N -l' -ldflags="-s -w" -o /server ./cmd/server

FROM alpine as release

COPY --from=build /server /server

ENTRYPOINT ["/server"]
