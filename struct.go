package main

import ("gopkg.in/mgo.v2/bson")

// Metadata Structure
type Metadata struct {
	Process struct {
		Ppid  int         `bson:"ppid" json:"ppid"`
		Pid   int         `bson:"pid" json:"pid"`
		Argv  []string    `bson:"argv" json:"argv"`
		Title interface{} `bson:"title" json:"title"`
	} `bson:"process" json:"process"`
	System struct {
		Platform     string `bson:"platform" json:"platform"`
		Hostname     string `bson:"hostname" json:"hostname"`
		Architecture string `bson:"architecture" json:"architecture"`
	} `bson:"system" json:"system"`
	Service struct {
		Name     string `bson:"name" json:"name"`
		Language struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"language" json:"language"`
		Agent struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"agent" json:"agent"`
		Environment interface{} `bson:"environment" json:"environment"`
		Framework   struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"framework" json:"framework"`
		Version string `bson:"version" json:"version"`
		Runtime struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"runtime" json:"runtime"`
	} `bson:"service" json:"service"`
}

// Transaction Structure
type Transaction struct {
	TraceID string `bson:"trace_id" json:"trace_id"`
	Result  string `bson:"result" json:"result"`
	Sampled bool   `bson:"sampled" json:"sampled"`
	Name    string `bson:"name" json:"name"`
	Context struct {
		Request struct {
			Body    string        `bson:"body" json:"body"`
			Cookies interface{} `bson:"cookies" json:"cookies"`
			Socket  struct {
				Encrypted     bool   `bson:"encrypted" json:"encrypted"`
				RemoteAddress string `bson:"remote_address" json:"remote_address"`
			} `bson:"socket" json:"socket"`
			URL struct {
				Pathname string `bson:"pathname" json:"pathname"`
				Full     string `bson:"full" json:"full"`
				Protocol string `bson:"protocol" json:"protocol"`
				Hostname string `bson:"hostname" json:"hostname"`
				Port     string `bson:"port" json:"port"`
			} `bson:"url" json:"url"`
			Headers interface{} `bson:"headers" json:"headers"`
			Env     struct {
				SERVERNAME string `bson:"SERVER_NAME" json:"SERVER_NAME"`
				SERVERPORT string `bson:"SERVER_PORT" json:"SERVER_PORT"`
				REMOTEADDR string `bson:"REMOTE_ADDR" json:"REMOTE_ADDR"`
			} `bson:"env" json:"env"`
			Method string `bson:"method" json:"method"`
		} `bson:"request" json:"request"`
		Response struct {
			StatusCode int           `bson:"status_code" json:"status_code"`
			Headers interface{} `bson:"headers" json:"headers"`
		} `bson:"response" json:"response"`
		Tags struct {
		} `bson:"tags" json:"tags"`
	} `bson:"context" json:"context"`
	Duration  float64 `bson:"duration" json:"duration"`
	Timestamp int64   `bson:"timestamp" json:"timestamp"`
	Types     string  `bson:"type" json:"type"`
	ID        string  `bson:"id" json:"id"`
	SpanCount struct {
		Started int `bson:"started" json:"started"`
		Dropped int `bson:"dropped" json:"dropped"`
	} `bson:"span_count" json:"span_count"`
}

// Data Access Object to manage database operations
type ApmDAO struct {
	Server   string
	Database string
}

// MongoObject Structure
type MongoObject struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	TraceID   string `bson:"trace_id" json:"trace_id"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
	Sampled bool   `bson:"sampled" json:"sampled"`
	Result    string `bson:"result" json:"result"`
	Metadata  struct {
		Service struct {
			Name string `bson:"name" json:"name"`
		} `bson:"service" json:"service"`
		Version  string `bson:"version" json:"version"`
		Language struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"language" json:"language"`
		Agent struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"agent" json:"agent"`
		Framework struct {
			Version string `bson:"version" json:"version"`
			Name    string `bson:"name" json:"name"`
		} `bson:"framework" json:"framework"`
	} `bson:"metadata" json:"metadata"`
	Request struct {
		URL     string        `bson:"url" json:"url"`
		Body    string        `bson:"body" json:"body"`
		Headers interface{} `bson:"headers" json:"headers"`
		Method  string        `bson:"method" json:"method"`
	} `bson:"request" json:"request"`
	Response struct {
		StatusCode int           `bson:"status_code" json:"status_code"`
		Headers interface{} `bson:"headers" json:"headers"`
	} `bson:"response" json:"response"`
	Duration float64 `bson:"duration" json:"duration"`
}