package main

// Config variables
const (
	apmServerIPAddress = "http://localhost"
	apmServerPort      = "8200"
)

// DB variable
const (
	TRANSACTION = "transactions"
	METADATA    = "metadata"
)

// Metadata Structure
type Metadata struct {
	Process struct {
		Ppid  int         `json:"ppid"`
		Pid   int         `json:"pid"`
		Argv  []string    `json:"argv"`
		Title interface{} `json:"title"`
	} `json:"process"`
	System struct {
		Platform     string `json:"platform"`
		Hostname     string `json:"hostname"`
		Architecture string `json:"architecture"`
	} `json:"system"`
	Service struct {
		Name     string `json:"name"`
		Language struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"language"`
		Agent struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"agent"`
		Environment interface{} `json:"environment"`
		Framework   struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"framework"`
		Version string `json:"version"`
		Runtime struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"runtime"`
	} `json:"service"`
}

// Transaction Structure
type Transaction struct {
	TraceID string `json:"trace_id"`
	Result  string `json:"result"`
	Sampled bool   `json:"sampled"`
	Name    string `json:"name"`
	Context struct {
		Request struct {
			Body    string        `json:"body"`
			Cookies interface{} `json:"cookies"`
			Socket  struct {
				Encrypted     bool   `json:"encrypted"`
				RemoteAddress string `json:"remote_address"`
			} `json:"socket"`
			URL struct {
				Pathname string `json:"pathname"`
				Full     string `json:"full"`
				Protocol string `json:"protocol"`
				Hostname string `json:"hostname"`
				Port     string `json:"port"`
			} `json:"url"`
			Headers interface{} `json:"headers"`
			Env     struct {
				SERVERNAME string `json:"SERVER_NAME"`
				SERVERPORT string `json:"SERVER_PORT"`
				REMOTEADDR string `json:"REMOTE_ADDR"`
			} `json:"env"`
			Method string `json:"method"`
		} `json:"request"`
		Response struct {
			StatusCode int           `json:"status_code"`
			Headers interface{} `json:"headers"`
		} `json:"response"`
		Tags struct {
		} `json:"tags"`
	} `json:"context"`
	Duration  float64 `json:"duration"`
	Timestamp int64   `json:"timestamp"`
	Types     string  `json:"type"`
	ID        string  `json:"id"`
	SpanCount struct {
		Started int `json:"started"`
		Dropped int `json:"dropped"`
	} `json:"span_count"`
}

// MongoObject Structure
type MongoObject struct {
	TraceID   string `json:"trace_id"`
	Timestamp int64  `json:"timestamp"`
	Result    string `json:"result"`
	Metadata  struct {
		Service struct {
			Name string `json:"name"`
		} `json:"service"`
		Version  string `json:"version"`
		Language struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"language"`
		Agent struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"agent"`
		Framework struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"framework"`
	} `json:"metadata"`
	Request struct {
		URL     string        `json:"url"`
		Body    string        `json:"body"`
		Headers interface{} `json:"headers"`
		Method  string        `json:"method"`
	} `json:"request"`
	Response struct {
		StatusCode int           `json:"status_code"`
		Headers interface{} `json:"headers"`
	} `json:"response"`
	Duration float64 `json:"duration"`
}
