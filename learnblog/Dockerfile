FROM golang:1.18-alpine3.14 as builder
WORKDIR /source
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN go install .

FROM alpine:3.14 as app_runtime
WORKDIR /app
COPY --from=builder /go/bin/learnblog /app
EXPOSE 8080
CMD ["/app/learnblog"]