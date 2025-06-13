FROM golang:1.24 AS builder

WORKDIR /src
COPY go.mod .
COPY *.go .

RUN CGO_ENABLED=0 go build -o app -ldflags "-w -s" -trimpath .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /src/app /usr/local/bin/app

USER 65534

EXPOSE 8080/tcp

ENTRYPOINT ["/usr/local/bin/app"]
