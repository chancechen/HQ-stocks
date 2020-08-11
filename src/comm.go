package main

import "time"

type roeData struct {
	StatDate   []string            `json:"statDate"`
	Data       map[string][]string `json:"data"`
	Name       map[string]string   `json:"name"`
	Updatetime time.Time           `json:"updatetime"`
	PeData     map[string][]string `json:"pedata"`
	PbData     map[string][]string `json:"pbdata"`
}

type config struct {
	MinRoe     float32 `json:"minRoe"`
	MaxRoe     float32 `json:"maxRoe"`
	ForceLoad  bool    `json:"forceLoad"`
	ShowMinRoe float32 `json:"showMinRoe"`
	ShowMaxRoe float32 `json:"showMaxRoe"`
	Host       string  `json:"host"`
	APIAccount string  `json:"apiAccount"`
	APIPwd     string  `json:"apiPwd"`
}
