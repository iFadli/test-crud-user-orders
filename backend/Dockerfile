FROM golang:1.18 AS build
WORKDIR /go/src/test-crud-user-orders
COPY . .

ENV CGO_ENABLED=0

RUN go build -o server ./cmd

FROM build AS development
RUN apt-get update \
    && apt-get install -y git
CMD ["go", "run", "cmd/main.go"]

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/test-crud-user-orders/server /server
CMD ["/server"]