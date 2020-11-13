FROM golang:1.15

WORKDIR /urlprober
COPY . /urlprober
RUN CGO_ENABLED=0 go build

FROM alpine:3.12.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /urlprober/urlprober .
CMD ["./urlprober"]
