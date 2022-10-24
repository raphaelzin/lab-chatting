FROM golang:1.19-alpine AS builder

RUN apk update && apk add --no-cache git 

RUN mkdir /build

WORKDIR /build

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -v -a -installsuffix cgo -o /go/bin/build

FROM alpine:3.16

COPY --from=builder /go/bin/build /go/bin/build

EXPOSE 8082

ENTRYPOINT ["/go/bin/build"]