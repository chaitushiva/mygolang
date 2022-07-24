FROM golang:1.18-alpine as builder

WORKDIR /go/src/healthchecker
COPY . .

RUN go mod init healthchecker && go build .

FROM alpine:3.15 as runtime

ENV PORT=8080
WORKDIR "/app"
RUN apk add wget
HEALTHCHECK --interval=10s --timeout=5s CMD wget --no-verbose --tries=1 --spider localhost:${PORT}/health

COPY --from=builder /go/src/healthchecker/healthchecker .
CMD ["./healthchecker"]
