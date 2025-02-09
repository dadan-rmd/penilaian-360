FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache git && apk add gcc libc-dev
WORKDIR /app
ENV GOSUMDB=off
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -ldflags '-linkmode=external' -o /app/penilaian-360 main.go

FROM golang:1.18-alpine
COPY --from=builder /app/penilaian-360 /app/penilaian-360
RUN apk add --no-cache tzdata ca-certificates libc6-compat

WORKDIR /app
ENTRYPOINT ["/app/penilaian-360"]
