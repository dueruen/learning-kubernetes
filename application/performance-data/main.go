package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"

	"go.opentelemetry.io/collector/pdata/plog"
)

// github.com/open-telemetry/opentelemetry-collector/blob/main/pdata/plog

var valuesMap map[int][]int64

var (
	inputPath     = flag.String("inputpath", os.Getenv("INPUT_PATH"), "INPUT_PATH")
	fileName      = flag.String("filename", os.Getenv("FILE_NAME"), "FILE_NAME")
	fileExtension = flag.String("fileextension", os.Getenv("FILE_EXTENSION"), "FILE_EXTENSION eg: .json")
	outputPath    = flag.String("outputpath", os.Getenv("OUTPUT_PATH"), "OUTPUT_PATH")
)

// const dis string = "Cilium with policies - 1000Hz 320 bytes"
// const fileName string = "data-wacmfozzzmsttz"
// const fileExtension string = ".json"
// const path string = "./data/data-00/" + fileName

func main() {
	flag.Parse()

	if *fileName == "" || *fileExtension == "" || *outputPath == "" || *inputPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// fmt.Println(time.Now())
	// t := time.Now().UnixMicro()
	// fmt.Println(t)
	// unixTimeUTC := time.UnixMicro(1652207057620399) //gives unix time stamp in utc

	// unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339) // converts utc time to RFC3339 format

	// fmt.Println("unix time stamp in UTC :--->", unixTimeUTC)
	// fmt.Println("unix time stamp in unitTimeInRFC3339 format :->", unitTimeInRFC3339)
	// frequencyInMicro := 1000000000000 / 1000

	if _, err := os.Stat(*outputPath); os.IsNotExist(err) {
		err = os.Mkdir(*outputPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	valuesMap = make(map[int][]int64)

	file, err := os.Open(*inputPath + "/" + *fileName + *fileExtension)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	start := time.Now()
	const maxCapacity int = 1000000000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// count := 30
	// var wg sync.WaitGroup

	for scanner.Scan() {
		// wg.Add(1)
		// go func() {
		// defer wg.Done()
		readLine(scanner.Bytes())
		// }()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)

	fmt.Println("Writing to file")
	writeFile()
	fmt.Println("Writing to file Done")

	// makePlots()

	elapsed = time.Since(start)
	log.Printf("Took %s", elapsed)
}

func writeFile() {
	file, err := os.OpenFile(*outputPath+"/"+*fileName+".csv", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	headers := []string{"id", "produce", "consume", "diff"}
	_, _ = datawriter.WriteString(strings.Join(headers, ",") + "\n")

	keys := maps.Keys(valuesMap)
	sort.Ints(keys)

	for _, key := range keys {
		vals := []string{fmt.Sprint(key)}
		vals = append(vals, strings.Fields(strings.Trim(fmt.Sprint(valuesMap[key]), "[]"))...)
		_, _ = datawriter.WriteString(strings.Join(vals, ",") + "\n")
	}

	datawriter.Flush()
	file.Close()
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
					words := strings.Fields(val)
					intVar, err := strconv.ParseInt(words[0], 10, 64)
					if err != nil {
						fmt.Println("ERROR: ", err)
					}
					id, err := strconv.Atoi(words[5])
					if err != nil {
						continue
					}

					if len(valuesMap[id]) == 0 {
						valuesMap[id] = []int64{-1, intVar, -1}
					} else {
						valuesMap[id] = []int64{valuesMap[id][0], intVar, intVar - valuesMap[id][0]}
					}

				} else if strings.Contains(val, "producer.produce") {
					words := strings.Fields(val)

					intVar, err := strconv.ParseInt(words[0], 10, 64)
					if err != nil {
						fmt.Println("ERROR: ", err)
					}
					id, err := strconv.Atoi(words[5])
					if err != nil {
						continue
					}

					if len(valuesMap[id]) == 0 {
						valuesMap[id] = []int64{intVar, -1, -1}
					} else {
						valuesMap[id] = []int64{intVar, valuesMap[id][1], valuesMap[id][1] - intVar}
					}
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

// func makePlots() {

// 	var values plotter.Values
// 	var valuesFloat []float64
// 	incomplete := 0
// 	for _, data := range valuesMap {
// 		if data[2] != -1 {
// 			s := float64(data[2]) / math.Pow(10, 6)
// 			values = append(values, s)
// 			valuesFloat = append(valuesFloat, s)
// 		} else {
// 			incomplete++
// 		}
// 	}

// 	file, err := os.OpenFile(*outputPath+"/"+*fileName+".txt", os.O_CREATE|os.O_WRONLY, 0644)

// 	if err != nil {
// 		log.Fatalf("failed creating file: %s", err)
// 	}

// 	datawriter := bufio.NewWriter(file)

// 	median, _ := stats.Median(valuesFloat)
// 	std, _ := stats.StandardDeviation(valuesFloat)
// 	_, _ = datawriter.WriteString("median: " + fmt.Sprintf("%f", median) + " std: " + fmt.Sprintf("%f", std) + " incomplete: " + fmt.Sprintf("%d", incomplete) + "\n")

// 	datawriter.Flush()
// 	file.Close()

// 	histPlot(values)
// 	boxPlot(values)
// }

// func histPlot(values plotter.Values) {
// 	p := plot.New()

// 	p.Title.Text = dis

// 	hist, err := plotter.NewHist(values, 20)
// 	if err != nil {
// 		panic(err)
// 	}
// 	p.Add(hist)

// 	if err := p.Save(3*vg.Inch, 3*vg.Inch, path+"/"+fileName+"-hist.png"); err != nil {
// 		panic(err)
// 	}
// }

// func boxPlot(values plotter.Values) {
// 	p := plot.New()

// 	p.Title.Text = dis

// 	box, err := plotter.NewBoxPlot(vg.Length(15), 0.0, values)
// 	if err != nil {
// 		panic(err)
// 	}
// 	p.Add(box)

// 	if err := p.Save(3*vg.Inch, 3*vg.Inch, path+"/"+fileName+"-box.png"); err != nil {
// 		panic(err)
// 	}
// }
