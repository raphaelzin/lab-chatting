FROM golang:1.19-alpine AS dev

RUN apk update && apk add --no-cache git 

# RUN mkdir /build

WORKDIR /home/go/app

# COPY . .

# RUN ls
# RUN go mod download

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/build
# FROM scratch
# COPY --from=builder /go/bin/build /go/bin/build
# ENTRYPOINT ["/go/bin/build"]

CMD ["go", "run", "main.go"]

EXPOSE 8082