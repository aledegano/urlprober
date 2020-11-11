FROM golang:1.15

WORKDIR /dyndns-updater
COPY . /dyndns-updater
RUN CGO_ENABLED=0 go build

FROM alpine:3.12.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /dyndns-updater/dyndnsupdater .
CMD ["./dyndnsupdater"]
