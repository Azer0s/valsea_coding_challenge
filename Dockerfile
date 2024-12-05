FROM golang:1.23.3 AS build
WORKDIR /go/src/app
COPY . .
RUN go get ./...
RUN go build -o /go/bin/app

FROM scratch
COPY --from=build /go/bin/app /app
CMD ["/app"]