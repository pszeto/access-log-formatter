package format

type App struct {
	config *Config
}

type Config struct {
	File string
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
