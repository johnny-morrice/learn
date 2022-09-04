FROM golang:1.19-alpine3.16 as builder
WORKDIR /source
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN go install .

FROM alpine:3.16 as app_runtime
WORKDIR /app
COPY --from=builder /go/bin/learnblog /app
COPY --from=builder /source/entrypoint.sh /app
COPY --from=builder /source/migrations/ /app
EXPOSE 8080
CMD ["/app/entrypoint.sh"]