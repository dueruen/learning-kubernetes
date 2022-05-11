package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartProducer(brokerList []string, kafkaTopic string) chan bool {
	shutdown := make(chan bool)

	if *messageSize == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *messageFrequency == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *producerRunTime == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	messageSizeInt, err := strconv.Atoi(*messageSize)
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("messageSizeInt: ", messageSizeInt)

	messageFrequency, err := strconv.Atoi(*messageFrequency)
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("messageFrequency: ", messageFrequency)

	producerRunTime, err := strconv.Atoi(*producerRunTime)
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("producerRunTime: ", producerRunTime)

	var tr trace.Tracer
	if IsInstrumented() {
		tr = otel.Tracer("producer")
	}

	rand.Seed(42)

	produceSignal := make(chan string, 10)
	producesigchan := make(chan os.Signal, 1)
	signal.Notify(producesigchan, syscall.SIGINT, syscall.SIGTERM)

	// create kafka client
	producer := newAccessLogProducer(brokerList)

	stopProducerChan := make(chan bool)
	var timerStop []chan bool

	waitTime := 10000 / messageFrequency
	for i := 0; i < 10; i++ {
		stop := make(chan bool)
		timerStop = append(timerStop, stop)
		go func(index int) {
			ticker := time.NewTicker(time.Duration(waitTime) * time.Millisecond)

			var count int64 = 0
			run := true
			for run {
				select {
				case _ = <-stop:
					ticker.Stop()
					run = false
					close(stop)
				case _ = <-ticker.C:
					count++
					time.Sleep(time.Duration(rand.Intn((waitTime * 100000))) * time.Nanosecond)
					produceSignal <- "" //fmt.Sprintf("%d-%d", index, s.UnixMicro())
				}
			}
		}(i)
	}

	count := 0
	go func() {
		if IsInDebug() {
			log.Println("Ready to start producer")
		}

		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		run := true
		for run == true {
			select {
			case _ = <-sigchan:
				producer.Close()
				run = false
				shutdown <- true
			case _ = <-stopProducerChan:
				fmt.Println("Producer have stopped scheduling new messages. Total scheduled messages is: ", count)
			case _ = <-produceSignal:
				count++
				go func() {
					//id := roundId //.Nanosecond()
					randBytes := RandASCIIBytes(messageSizeInt)

					if len(sigchan) > 3 {
						fmt.Println("!!!!!!!!!!!!!! Producer cannot produce fast enough !!!!!!!!!!!!!")
					}
					t := time.Now().UnixMicro()
					fmt.Println(t, "producer.produce", "app=", *appName, "id=", count)

					var ctx context.Context
					var span trace.Span
					if IsInstrumented() {
						ctx, span = tr.Start(context.Background(), "produce message")
						defer span.End()
						fmt.Println("traceId: ", span.SpanContext().TraceID())
					}

					var keyValue string
					if messageSizeInt == 0 {
						keyValue = ""
					} else {
						keyValue = "random__bytes"
					}

					msg := sarama.ProducerMessage{
						Topic: kafkaTopic,
						Key:   sarama.StringEncoder(keyValue),
						Value: sarama.ByteEncoder(randBytes),
						Headers: []sarama.RecordHeader{
							{Key: []byte("id"), Value: []byte(fmt.Sprintf("%d", count))}},
					}

					if IsInstrumented() {
						otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(&msg))
					}

					// size := byteSize(&msg)
					// fmt.Println("Message size: ", size)

					producer.Input() <- &msg
					successMsg := <-producer.Successes()
					if IsInDebug() {
						log.Println("Successful to write message, offset:", successMsg.Offset)
					}
				}()
			}
		}
	}()

	go func() {
		time.Sleep((time.Duration(producerRunTime) * time.Minute) + time.Duration(waitTime*100000)*time.Nanosecond)
		// ticker.Stop()
		stopProducerChan <- true
		for _, c := range timerStop {
			c <- true
		}
	}()

	return shutdown
}

func newAccessLogProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Producer.Return.Successes = true

	var producer sarama.AsyncProducer
	retries := 5
	retry := true
	for retry {
		pro, err := sarama.NewAsyncProducer(brokerList, config)
		if err != nil {
			if retries != 0 {
				retries--
				log.Println("RETRY: Failed to start Sarama producer:", err)
				time.Sleep(2 * time.Second)
				continue
			}
			log.Fatalln("Failed to start Sarama producer:", err)
		}
		log.Println("Started producer group")
		producer = pro
		retry = false
	}

	if IsInstrumented() {
		// Wrap instrumentation
		producer = otelsarama.WrapAsyncProducer(config, producer)
	}

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write message:", err)
		}
	}()

	if IsInDebug() {
		log.Println("Producer created")
	}
	return producer
}

func byteSize(m *sarama.ProducerMessage) int {
	maximumRecordOverhead := 5*binary.MaxVarintLen32 + binary.MaxVarintLen64 + 1
	size := maximumRecordOverhead
	for _, h := range m.Headers {
		fmt.Println("Header key: ", string(h.Key), "  - Header value: ", string(h.Value))
		size += len(h.Key) + len(h.Value) + 2*binary.MaxVarintLen32
	}

	if m.Key != nil {
		size += m.Key.Length()
	}
	if m.Value != nil {
		fmt.Println("Value size: ", m.Value.Length())
		size += m.Value.Length()
	}
	return size
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)
	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)
	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])
		// random % 64
		randomPos := random % uint8(l)
		// put into output
		output[pos] = letterBytes[randomPos]
	}

	return output
}
