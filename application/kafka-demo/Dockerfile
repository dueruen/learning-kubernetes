FROM golang:1.18.0 as builder

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kafka-demos

FROM scratch
COPY --from=builder /app/kafka-demos .
CMD ["./kafka-demos"]