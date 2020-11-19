package load

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	url           = flag.String("url", "", "(Required)Defines url of the service")
	showResponses = flag.Bool("show-response", false, "Defines whether show response body or not.")
	method        = flag.String("method", "GET", "Defines method of the request")
	runTime       = flag.Int("runtime", 60, "Defines running time in second")
	requestCount  = flag.Int("request-count", 100, "Defines count of requests")
	maxConcurrecy = flag.Int("max-concurrency", 1, "Defines how many process will be used")
	parameters    = flag.String("parameters", "", "Defines parameter for the request")
	data          = flag.String("data", "", "Defines data will be sent")
	getHeaders    = flag.String("headers", "", "Defines headers of the request. You should seperate headers with comma(,)")
	timeout       = flag.Int("request-timeout", 10, "Defines timeout time of the per requests")
	sem           = make(chan int, *maxConcurrecy)
)

func main() {
	flag.Parse()
	var responseCodes [6]int
	responseBodies := []string{""}
	var statusCode int
	var responseBody string
	var j int = 0
	if *url == "" {
		log.Println("Required flags are not set")
		flag.PrintDefaults()
		os.Exit(1)
	}
	startTime := time.Now()
	for i := 0; i < *requestCount+1; i++ {
		sem <- 0
		go func() {

			statusCode, responseBody = LoadGenerator()
			if *showResponses == true {
				if responseBodies[j] != responseBody {
					responseBodies = append(responseBodies, responseBody)
					j++
				}
			}
			switch statusCode / 100 {
			case 1:
				responseCodes[1]++
				break
			case 2:
				responseCodes[2]++
				break
			case 3:
				responseCodes[3]++
				break
			case 4:
				responseCodes[4]++
				break
			case 5:
				responseCodes[5]++
				break
			default:
				responseCodes[0]++
				break
			}
			<-sem
		}()

	}
	endTime := time.Now()
	if *showResponses == true {
		fmt.Println("Responses:")
	}
	fmt.Println(responseBodies)
	fmt.Println("It is done in:", endTime.Sub(startTime))
	for index, _ := range responseCodes {
		fmt.Printf("Count of %dxx : %d \n", index, responseCodes[index])
	}

}

func LoadGenerator() (int, string) {
	headers := strings.Split(*getHeaders, ",")
	return MakeRequest(*url, *method, headers, *parameters, *timeout, *data)
}
