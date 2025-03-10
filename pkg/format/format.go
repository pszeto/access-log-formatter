package format

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func New(config *Config) *App {
	return &App{
		config: config,
	}
}

// declaring a struct
type LogMessage struct {

	// defining struct variables
	START_TIME                            string
	REQUEST_METHOD                        string
	REQUEST_PATH_OR_X_ENVOY_ORIGINAL_PATH string
	REQUEST_PROTOCOL                      string
	RESPONSE_CODE                         string
	RESPONSE_FLAGS                        string
	RESPONSE_CODE_DETAILS                 string
	CONNECTION_TERMINATION_DETAILS        string
	UPSTREAM_TRANSPORT_FAILURE_REASON     string
	BYTES_RECEIVED                        string
	BYTES_SENT                            string
	DURATION                              string
	REQUEST_X_ENVOY_UPSTREAM_SERVICE_TIME string
	REQUEST_X_FORWARDED_FOR               string
	REQUEST_USER_AGENT                    string
	REQUEST_X_REQUEST_ID                  string
	REQUEST_AUTHORITY                     string
	UPSTREAM_HOST                         string
	UPSTREAM_CLUSTER_RAW                  string
	UPSTREAM_LOCAL_ADDRESS                string
	DOWNSTREAM_LOCAL_ADDRESS              string
	DOWNSTREAM_REMOTE_ADDRESS             string
	REQUESTED_SERVER_NAME                 string
	ROUTE_NAME                            string
}

func removeQuotes(text string) string {
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		return text[1 : len(text)-1]
	}

	return text
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
	//fmt.Println("====== Formatted Log Line ======")
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

	// START_TIME :  [2025-03-05T17:25:17.486Z]
	// %REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL% :  "- - -"
	// RESPONSE_CODE :  0
	// RESPONSE_FLAGS :  UH
	// RESPONSE_CODE_DETAILS :  -
	// CONNECTION_TERMINATION_DETAILS :  -
	// UPSTREAM_TRANSPORT_FAILURE_REASON :  "-"
	// BYTES_RECEIVED :  0
	// BYTES_SENT :  0
	// DURATION :  0
	// RESP(X-ENVOY-UPSTREAM-SERVICE-TIME) :  -
	// REQ(X-FORWARDED-FOR) :  "-"
	// REQ(USER-AGENT) :  "-"
	// REQ(X-REQUEST-ID) :  "-"
	// REQ(:AUTHORITY) :  "-"
	// UPSTREAM_HOST :  "-"
	// UPSTREAM_CLUSTER_RAW :  BlackHoleCluster
	// UPSTREAM_LOCAL_ADDRESS :  -
	// DOWNSTREAM_LOCAL_ADDRESS :  192.168.184.80:4317
	// DOWNSTREAM_REMOTE_ADDRESS :  192.168.8.121:50090
	// REQUESTED_SERVER_NAME :  -
	// ROUTE_NAME :  -

	var logMessage LogMessage
	logSplit := strings.Split(newLog, " ")
	if len(logSplit) == 22 {
		for i := 0; i < len(logSplit); i++ {
			//fmt.Println(logSplit[i])
			switch i {
			case 0:
				//fmt.Println("START_TIME : ", logSplit[i])
				logMessage.START_TIME = logSplit[i]
			case 1:
				//fmt.Println("%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL% : ", replacementStrings[0])
				tmp := strings.Fields(removeQuotes(replacementStrings[0]))
				logMessage.REQUEST_METHOD = tmp[0]
				logMessage.REQUEST_PATH_OR_X_ENVOY_ORIGINAL_PATH = tmp[1]
				logMessage.REQUEST_PROTOCOL = tmp[2]
			case 2:
				//fmt.Println("RESPONSE_CODE : ", logSplit[i])
				logMessage.RESPONSE_CODE = logSplit[i]
			case 3:
				//fmt.Println("RESPONSE_FLAGS : ", logSplit[i])
				logMessage.RESPONSE_FLAGS = logSplit[i]
			case 4:
				//fmt.Println("RESPONSE_CODE_DETAILS : ", logSplit[i])
				logMessage.RESPONSE_CODE_DETAILS = logSplit[i]
			case 5:
				//fmt.Println("CONNECTION_TERMINATION_DETAILS : ", logSplit[i])
				logMessage.CONNECTION_TERMINATION_DETAILS = logSplit[i]
			case 6:
				//fmt.Println("UPSTREAM_TRANSPORT_FAILURE_REASON : ", replacementStrings[1])
				logMessage.UPSTREAM_TRANSPORT_FAILURE_REASON = removeQuotes(replacementStrings[1])
			case 7:
				//fmt.Println("BYTES_RECEIVED : ", logSplit[i])
				logMessage.BYTES_RECEIVED = logSplit[i]
			case 8:
				//fmt.Println("BYTES_SENT : ", logSplit[i])
				logMessage.BYTES_SENT = logSplit[i]
			case 9:
				//fmt.Println("DURATION : ", logSplit[i])
				logMessage.DURATION = logSplit[i]
			case 10:
				//fmt.Println("RESP(X-ENVOY-UPSTREAM-SERVICE-TIME) : ", logSplit[i])
				logMessage.REQUEST_X_ENVOY_UPSTREAM_SERVICE_TIME = logSplit[i]
			case 11:
				//fmt.Println("REQ(X-FORWARDED-FOR) : ", replacementStrings[2])
				logMessage.REQUEST_X_FORWARDED_FOR = removeQuotes(replacementStrings[2])
			case 12:
				//fmt.Println("REQ(USER-AGENT) : ", replacementStrings[3])
				logMessage.REQUEST_USER_AGENT = removeQuotes(replacementStrings[3])
			case 13:
				//fmt.Println("REQ(X-REQUEST-ID) : ", replacementStrings[4])
				logMessage.REQUEST_X_REQUEST_ID = removeQuotes(replacementStrings[4])
			case 14:
				//fmt.Println("REQ(:AUTHORITY) : ", replacementStrings[5])
				logMessage.REQUEST_AUTHORITY = removeQuotes(replacementStrings[5])
			case 15:
				//fmt.Println("UPSTREAM_HOST : ", replacementStrings[6])
				logMessage.UPSTREAM_HOST = removeQuotes(replacementStrings[6])
			case 16:
				//fmt.Println("UPSTREAM_CLUSTER_RAW : ", logSplit[i])
				logMessage.UPSTREAM_CLUSTER_RAW = logSplit[i]
			case 17:
				//fmt.Println("UPSTREAM_LOCAL_ADDRESS : ", logSplit[i])
				logMessage.UPSTREAM_LOCAL_ADDRESS = logSplit[i]
			case 18:
				//fmt.Println("DOWNSTREAM_LOCAL_ADDRESS : ", logSplit[i])
				logMessage.DOWNSTREAM_LOCAL_ADDRESS = logSplit[i]
			case 19:
				//fmt.Println("DOWNSTREAM_REMOTE_ADDRESS : ", logSplit[i])
				logMessage.DOWNSTREAM_REMOTE_ADDRESS = logSplit[i]
			case 20:
				//fmt.Println("REQUESTED_SERVER_NAME : ", logSplit[i])
				logMessage.REQUESTED_SERVER_NAME = logSplit[i]
			case 21:
				//fmt.Println("ROUTE_NAME : ", logSplit[i])
				logMessage.ROUTE_NAME = logSplit[i]
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
				//fmt.Println("START_TIME : ", logSplit[i])
				logMessage.START_TIME = logSplit[i]
			case 1:
				//fmt.Println("%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL% : ", replacementStrings[0])
				tmp := strings.Fields(removeQuotes(replacementStrings[0]))
				logMessage.REQUEST_METHOD = tmp[0]
				logMessage.REQUEST_PATH_OR_X_ENVOY_ORIGINAL_PATH = tmp[1]
				logMessage.REQUEST_PROTOCOL = tmp[2]
			case 2:
				//fmt.Println("RESPONSE_CODE : ", logSplit[i])
				logMessage.RESPONSE_CODE = logSplit[i]
			case 3:
				//fmt.Println("RESPONSE_FLAGS : ", logSplit[i])
				logMessage.RESPONSE_FLAGS = logSplit[i]
			case 4:
				//fmt.Println("BYTES_RECEIVED : ", logSplit[i])
				logMessage.BYTES_RECEIVED = logSplit[i]
			case 5:
				//fmt.Println("BYTES_SENT : ", logSplit[i])
				logMessage.BYTES_SENT = logSplit[i]
			case 6:
				//fmt.Println("DURATION : ", logSplit[i])
				logMessage.DURATION = logSplit[i]
			case 7:
				//fmt.Println("%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% : ", logSplit[i])
				logMessage.REQUEST_X_ENVOY_UPSTREAM_SERVICE_TIME = logSplit[i]
			case 8:
				//fmt.Println("REQ(X-FORWARDED-FOR) : ", replacementStrings[1])
				logMessage.REQUEST_X_FORWARDED_FOR = removeQuotes(replacementStrings[1])
			case 9:
				//fmt.Println("REQ(USER-AGENT) : ", replacementStrings[2])
				logMessage.REQUEST_USER_AGENT = removeQuotes(replacementStrings[2])
			case 10:
				//fmt.Println("REQ(X-REQUEST-ID) : ", replacementStrings[3])
				logMessage.REQUEST_X_REQUEST_ID = removeQuotes(replacementStrings[3])
			case 11:
				//fmt.Println("REQ(:AUTHORITY) : ", replacementStrings[4])
				logMessage.REQUEST_AUTHORITY = removeQuotes(replacementStrings[4])
			case 12:
				//fmt.Println("UPSTREAM_HOST : ", replacementStrings[5])
				logMessage.UPSTREAM_HOST = removeQuotes(replacementStrings[5])
			}
		}
	}
	log_message, err := json.MarshalIndent(logMessage, "", "    ")

	if err != nil {

		// if error is not nil
		// print error
		fmt.Println(err)
	}
	// fmt.Println("json")
	// as human_enc is in a byte array
	// format, it needs to be
	// converted into a string
	fmt.Println(string(log_message))
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
