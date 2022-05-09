package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"go.opentelemetry.io/collector/pdata/plog"
)

// github.com/open-telemetry/opentelemetry-collector/blob/main/pdata/plog

var valuesMap map[string][]int

const dis string = "Cilium with policies - 10Hz 32000 bytes"
const fileName string = "data-wacmfozmsttzzz"
const fileExtension string = ".json"
const path string = "./data/" + fileName

func main() {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	valuesMap = make(map[string][]int)

	file, err := os.Open("./data/" + fileName + fileExtension)
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

	makePlots()

	elapsed = time.Since(start)
	log.Printf("Took %s", elapsed)
}

func makePlots() {

	var values plotter.Values
	var valuesFloat []float64
	for _, data := range valuesMap {
		s := float64(data[2]) / math.Pow(10, 6)
		values = append(values, s)
		valuesFloat = append(valuesFloat, s)
	}

	file, err := os.OpenFile(path+"/"+fileName+".txt", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	median, _ := stats.Median(valuesFloat)
	std, _ := stats.StandardDeviation(valuesFloat)
	_, _ = datawriter.WriteString("median: " + fmt.Sprintf("%f", median) + " std: " + fmt.Sprintf("%f", std) + "\n")

	datawriter.Flush()
	file.Close()

	histPlot(values)
	boxPlot(values)
}

func histPlot(values plotter.Values) {
	p := plot.New()

	p.Title.Text = dis

	hist, err := plotter.NewHist(values, 20)
	if err != nil {
		panic(err)
	}
	p.Add(hist)

	if err := p.Save(3*vg.Inch, 3*vg.Inch, path+"/"+fileName+"-hist.png"); err != nil {
		panic(err)
	}
}

func boxPlot(values plotter.Values) {
	p := plot.New()

	p.Title.Text = dis

	box, err := plotter.NewBoxPlot(vg.Length(15), 0.0, values)
	if err != nil {
		panic(err)
	}
	p.Add(box)

	if err := p.Save(3*vg.Inch, 3*vg.Inch, path+"/"+fileName+"-box.png"); err != nil {
		panic(err)
	}
}

func writeFile() {
	file, err := os.OpenFile(path+"/"+fileName+".csv", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range valuesMap {
		_, _ = datawriter.WriteString(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data)), ","), "[]") + "\n")
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

					intVar, err := strconv.Atoi(words[0])
					if err != nil {
						fmt.Println("ERROR: ", err)
					}

					if len(valuesMap[words[5]]) == 0 {
						valuesMap[words[5]] = []int{-1, intVar, -1}
					} else {
						valuesMap[words[5]] = []int{valuesMap[words[5]][0], intVar, intVar - valuesMap[words[5]][0]}
					}

				} else if strings.Contains(val, "producer.produce") {
					words := strings.Fields(val)

					intVar, err := strconv.Atoi(words[0])
					if err != nil {
						fmt.Println("ERROR: ", err)
					}

					if len(valuesMap[words[5]]) == 0 {
						valuesMap[words[5]] = []int{intVar, -1, -1}
					} else {
						valuesMap[words[5]] = []int{intVar, valuesMap[words[5]][1], valuesMap[words[5]][1] - intVar}
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
