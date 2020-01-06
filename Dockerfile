FROM golang:latest as builder
LABEL maintainer="Tyr Mactire <tyr@pettingzoo.co>"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/gobuffalo/packr/v2/packr2
COPY . .
RUN packr2
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"] 