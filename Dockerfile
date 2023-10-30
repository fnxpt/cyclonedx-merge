FROM golang:latest AS builder
WORKDIR /go/src/
ENV CGO_ENABLED=0

ADD . .

RUN go build -ldflags="-s -w" -o release/cyclonedx-merge .

FROM scratch AS runtime
WORKDIR /

COPY --from=builder /go/src/release/cyclonedx-merge .

ENTRYPOINT [ "/cyclonedx-merge" ]