package format

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func New(config *Config) *App {
	return &App{
		config: config,
	}
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

func printFormattedLog(log string) error {
	// find all indices of the quote character
	allInstanceOfQuote := findAllInstances(log, `"`)
	// check for proper number of quotes.  Istio should have 14 and envoy should have 12
	fmt.Println("len(allInstanceOfQuote) = ", len(allInstanceOfQuote))
	if len(allInstanceOfQuote) != 14 && len(allInstanceOfQuote) != 12 {
		return errors.New("invalid log message")
	}
	replacementStrings := []string{}
	// iterate through the log message to pull out quoted strings, these are special entries
	for i := 0; i < len(allInstanceOfQuote); i = i + 2 {
		substring := log[allInstanceOfQuote[i] : allInstanceOfQuote[i+1]+1]
		replacementStrings = append(replacementStrings, substring)
	}
	newLog := log
	// iterate through the log message to replace quote string with temp string
	for i := 0; i < len(replacementStrings); i++ {
		newLog = strings.Replace(newLog, replacementStrings[i], "=======", 1)
	}

	var logMessage LogMessage
	logSplit := strings.Split(newLog, " ")
	fmt.Println("len(logSplit) = ", len(logSplit))
	if len(logSplit) == 22 {
		for i := 0; i < len(logSplit); i++ {
			switch i {
			case 0:
				logMessage.START_TIME = logSplit[i]
			case 1:
				tmp := strings.Fields(removeQuotes(replacementStrings[0]))
				logMessage.REQUEST_METHOD = tmp[0]
				logMessage.REQUEST_PATH_OR_X_ENVOY_ORIGINAL_PATH = tmp[1]
				logMessage.REQUEST_PROTOCOL = tmp[2]
			case 2:
				logMessage.RESPONSE_CODE = logSplit[i]
			case 3:
				logMessage.RESPONSE_FLAGS = logSplit[i]
			case 4:
				logMessage.RESPONSE_CODE_DETAILS = logSplit[i]
			case 5:
				logMessage.CONNECTION_TERMINATION_DETAILS = logSplit[i]
			case 6:
				logMessage.UPSTREAM_TRANSPORT_FAILURE_REASON = removeQuotes(replacementStrings[1])
			case 7:
				logMessage.BYTES_RECEIVED = logSplit[i]
			case 8:
				logMessage.BYTES_SENT = logSplit[i]
			case 9:
				logMessage.DURATION = logSplit[i]
			case 10:
				logMessage.REQUEST_X_ENVOY_UPSTREAM_SERVICE_TIME = logSplit[i]
			case 11:
				logMessage.REQUEST_X_FORWARDED_FOR = removeQuotes(replacementStrings[2])
			case 12:
				logMessage.REQUEST_USER_AGENT = removeQuotes(replacementStrings[3])
			case 13:
				logMessage.REQUEST_X_REQUEST_ID = removeQuotes(replacementStrings[4])
			case 14:
				logMessage.REQUEST_AUTHORITY = removeQuotes(replacementStrings[5])
			case 15:
				logMessage.UPSTREAM_HOST = removeQuotes(replacementStrings[6])
			case 16:
				logMessage.UPSTREAM_CLUSTER_RAW = logSplit[i]
			case 17:
				logMessage.UPSTREAM_LOCAL_ADDRESS = logSplit[i]
			case 18:
				logMessage.DOWNSTREAM_LOCAL_ADDRESS = logSplit[i]
			case 19:
				logMessage.DOWNSTREAM_REMOTE_ADDRESS = logSplit[i]
			case 20:
				logMessage.REQUESTED_SERVER_NAME = logSplit[i]
			case 21:
				logMessage.ROUTE_NAME = logSplit[i]
			}
		}
	} else if len(logSplit) == 13 {
		// [%START_TIME%] "%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL%"
		// %RESPONSE_CODE% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION%
		// %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% "%REQ(X-FORWARDED-FOR)%" "%REQ(USER-AGENT)%"
		// "%REQ(X-REQUEST-ID)%" "%REQ(:AUTHORITY)%" "%UPSTREAM_HOST%"\n
		for i := 0; i < len(logSplit); i++ {
			switch i {
			case 0:
				logMessage.START_TIME = logSplit[i]
			case 1:
				tmp := strings.Fields(removeQuotes(replacementStrings[0]))
				logMessage.REQUEST_METHOD = tmp[0]
				logMessage.REQUEST_PATH_OR_X_ENVOY_ORIGINAL_PATH = tmp[1]
				logMessage.REQUEST_PROTOCOL = tmp[2]
			case 2:
				logMessage.RESPONSE_CODE = logSplit[i]
			case 3:
				logMessage.RESPONSE_FLAGS = logSplit[i]
			case 4:
				logMessage.BYTES_RECEIVED = logSplit[i]
			case 5:
				logMessage.BYTES_SENT = logSplit[i]
			case 6:
				logMessage.DURATION = logSplit[i]
			case 7:
				logMessage.REQUEST_X_ENVOY_UPSTREAM_SERVICE_TIME = logSplit[i]
			case 8:
				logMessage.REQUEST_X_FORWARDED_FOR = removeQuotes(replacementStrings[1])
			case 9:
				logMessage.REQUEST_USER_AGENT = removeQuotes(replacementStrings[2])
			case 10:
				logMessage.REQUEST_X_REQUEST_ID = removeQuotes(replacementStrings[3])
			case 11:
				logMessage.REQUEST_AUTHORITY = removeQuotes(replacementStrings[4])
			case 12:
				logMessage.UPSTREAM_HOST = removeQuotes(replacementStrings[5])
			}
		}
	} else {
		return errors.New("invalid log message")
	}

	log_message, err := json.MarshalIndent(logMessage, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(log_message))
	return nil
}

func (app *App) Entry() error {
	if app.config.File == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Log Line : ")
		log, _ := reader.ReadString('\n')
		//removes the newline character
		log = strings.TrimSuffix(log, "\n")
		err := printFormattedLog(log)
		if err != nil {
			fmt.Println(err)
		}
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
			line := scanner.Text()
			// check if it's an empty line
			if utf8.RuneCountInString(line) > 0 {
				// Envoy and Istio default access log will always start with [%START_TIME%]
				if line[0:1] == "[" {
					err := printFormattedLog(line)
					fmt.Println(err)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
	return nil
}
