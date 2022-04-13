FROM golang:1.18.1 as builder

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o confluent-demo

FROM scratch
COPY --from=builder /app/confluent-demo .
CMD ["./confluent-demo"]