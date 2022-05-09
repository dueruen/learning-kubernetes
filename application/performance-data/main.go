package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/collector/pdata/plog"
)

// github.com/open-telemetry/opentelemetry-collector/blob/main/pdata/plog

func main() {
	file, err := os.Open("./data/data-first-02.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	start := time.Now()
	const maxCapacity int = 1000000000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	count := 30
	var wg sync.WaitGroup

	for count > 0 && scanner.Scan() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readLine(scanner.Bytes())
		}()

		count--
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func readLine(jsonBuf []byte) {
	decoder := plog.NewJSONUnmarshaler()
	var got plog.Logs
	got, err := decoder.UnmarshalLogs(jsonBuf)
	check(err)

	// fmt.Println(got.LogRecordCount())
	fmt.Println("Count: ", got.ResourceLogs().Len())

	// // var ss plog.internal.ResourceLogsSlice
	// // var s = got.ResourceLogs()

	for i := 0; i < got.ResourceLogs().Len(); i++ {
		resourceLog := got.ResourceLogs().At(i)

		for i := 0; i < resourceLog.ScopeLogs().Len(); i++ {
			scopeLog := resourceLog.ScopeLogs().At(i)

			for i := 0; i < scopeLog.LogRecords().Len(); i++ {
				logRecord := scopeLog.LogRecords().At(i)
				val := logRecord.Body().AsString()
				if strings.Contains(val, "consumer.consumed") {
					fmt.Println(val)
				}

			}
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
