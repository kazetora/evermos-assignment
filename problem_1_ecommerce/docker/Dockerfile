FROM golang:1.15-alpine AS builder

RUN mkdir /app

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/exe

FROM scratch 
COPY --from=builder /go/bin/exe /
ENTRYPOINT ["/exe"]