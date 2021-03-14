FROM golang:1.15 as builder

WORKDIR /src/httpstat-monitor

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /src/httpstat-monitor/bin/httpstat-monitor /
ENTRYPOINT ["/httpstat-monitor"]