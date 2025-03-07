package format

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func New(config *Config) *App {
	return &App{
		config: config,
	}
}

func findAllInstances(text string, letter string) []int {
	var indices []int
	for i, char := range text {
		if string(char) == letter {
			indices = append(indices, i)
		}
	}
	return indices
}

func printFormattedLog(log string) {
	fmt.Println("====== Formatted Log Line ======")
	allInstanceOfQuote := findAllInstances(log, `"`)
	// fmt.Println(allInstanceOfQuote)
	replacementStrings := []string{}
	for i := 0; i < len(allInstanceOfQuote); i = i + 2 {
		//fmt.Printf("allInstanceOfQuote[%d] = %d\n", i, allInstanceOfQuote[i])
		substring := log[allInstanceOfQuote[i] : allInstanceOfQuote[i+1]+1]
		// fmt.Println(substring)
		replacementStrings = append(replacementStrings, substring)
	}
	newLog := log
	for i := 0; i < len(replacementStrings); i++ {
		newLog = strings.Replace(newLog, replacementStrings[i], "=======", 1)

	}
	// fmt.Println(log)
	// fmt.Println(newLog)
	logSplit := strings.Split(newLog, " ")
	if len(logSplit) == 22 {
		for i := 0; i < len(logSplit); i++ {
			//fmt.Println(logSplit[i])
			switch i {
			case 0:
				fmt.Println("START_TIME : ", logSplit[i])
			case 1:
				fmt.Println("%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL% : ", replacementStrings[0])
			case 2:
				fmt.Println("RESPONSE_CODE : ", logSplit[i])
			case 3:
				fmt.Println("RESPONSE_FLAGS : ", logSplit[i])
			case 4:
				fmt.Println("RESPONSE_CODE_DETAILS : ", logSplit[i])
			case 5:
				fmt.Println("CONNECTION_TERMINATION_DETAILS : ", logSplit[i])
			case 6:
				fmt.Println("UPSTREAM_TRANSPORT_FAILURE_REASON : ", replacementStrings[1])
			case 7:
				fmt.Println("BYTES_RECEIVED : ", logSplit[i])
			case 8:
				fmt.Println("BYTES_SENT : ", logSplit[i])
			case 9:
				fmt.Println("DURATION : ", logSplit[i])
			case 10:
				fmt.Println("RESP(X-ENVOY-UPSTREAM-SERVICE-TIME) : ", logSplit[i])
			case 11:
				fmt.Println("REQ(X-FORWARDED-FOR) : ", replacementStrings[2])
			case 12:
				fmt.Println("REQ(USER-AGENT) : ", replacementStrings[3])
			case 13:
				fmt.Println("REQ(X-REQUEST-ID) : ", replacementStrings[4])
			case 14:
				fmt.Println("REQ(:AUTHORITY) : ", replacementStrings[5])
			case 15:
				fmt.Println("UPSTREAM_HOST : ", replacementStrings[6])
			case 16:
				fmt.Println("UPSTREAM_CLUSTER_RAW : ", logSplit[i])
			case 17:
				fmt.Println("UPSTREAM_LOCAL_ADDRESS : ", logSplit[i])
			case 18:
				fmt.Println("DOWNSTREAM_LOCAL_ADDRESS : ", logSplit[i])
			case 19:
				fmt.Println("DOWNSTREAM_REMOTE_ADDRESS : ", logSplit[i])
			case 20:
				fmt.Println("REQUESTED_SERVER_NAME : ", logSplit[i])
			case 21:
				fmt.Println("ROUTE_NAME : ", logSplit[i])
			}
		}
	} else {

		// [%START_TIME%] "%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL%"
		// %RESPONSE_CODE% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION%
		// %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% "%REQ(X-FORWARDED-FOR)%" "%REQ(USER-AGENT)%"
		// "%REQ(X-REQUEST-ID)%" "%REQ(:AUTHORITY)%" "%UPSTREAM_HOST%"\n
		for i := 0; i < len(logSplit); i++ {
			//fmt.Println(logSplit[i])
			switch i {
			case 0:
				fmt.Println("START_TIME : ", logSplit[i])
			case 1:
				fmt.Println("%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL% : ", replacementStrings[0])
			case 2:
				fmt.Println("RESPONSE_CODE : ", logSplit[i])
			case 3:
				fmt.Println("RESPONSE_FLAGS : ", logSplit[i])
			case 4:
				fmt.Println("BYTES_RECEIVED : ", logSplit[i])
			case 5:
				fmt.Println("BYTES_SENT : ", logSplit[i])
			case 6:
				fmt.Println("DURATION : ", logSplit[i])
			case 7:
				fmt.Println("%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% : ", logSplit[i])
			case 8:
				fmt.Println("REQ(X-FORWARDED-FOR) : ", replacementStrings[1])
			case 9:
				fmt.Println("REQ(USER-AGENT) : ", replacementStrings[2])
			case 10:
				fmt.Println("REQ(X-REQUEST-ID) : ", replacementStrings[3])
			case 11:
				fmt.Println("REQ(:AUTHORITY) : ", replacementStrings[4])
			case 12:
				fmt.Println("UPSTREAM_HOST : ", replacementStrings[5])
			}
		}
	}
}

func (app *App) Entry() error {
	if app.config.File == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Log Line : ")
		log, _ := reader.ReadString('\n')
		log = strings.TrimSuffix(log, "\n") //removes the newline character
		printFormattedLog(log)
	} else {
		file, err := os.Open(app.config.File)
		if err != nil {
			fmt.Print(err)
		}
		defer file.Close()

		// Create a new scanner to read the file line by line.
		scanner := bufio.NewScanner(file)

		// Iterate over the scanner, printing each line.
		for scanner.Scan() {
			printFormattedLog(scanner.Text())
		}
	}
	return nil
}
