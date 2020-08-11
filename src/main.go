package main

import (
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/go-echarts/go-echarts/charts"
)

const (
	stockjson = "./assets/stock.json"
	confjson  = "./assets/conf.json"

	maxNum = 50
)

type router struct {
	name string
	charts.RouterOpts
}

var (
	host    = "http://127.0.0.1:8080"
	routers = []router{
		{"bar", charts.RouterOpts{URL: host + "/bar", Text: "Bar"}},
	}
)

var seed = rand.NewSource(time.Now().UnixNano())

func orderRouters(chartType string) []charts.RouterOpts {
	for i := 0; i < len(routers); i++ {
		if routers[i].name == chartType {
			routers[i], routers[0] = routers[0], routers[i]
			break
		}
	}

	rs := make([]charts.RouterOpts, 0)
	for i := 0; i < len(routers); i++ {
		rs = append(rs, routers[i].RouterOpts)
	}
	return rs
}

func getRenderPath(f string) string {
	return path.Join("./", f)
}

func logTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s\n", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func main() {

	loadData()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// Avoid "404 page not found".
	http.HandleFunc("/", logTracing(barHandler))
	http.HandleFunc("/bar", logTracing(barHandler))

	log.Println("Run server at " + host)
	http.ListenAndServe(":8080", nil)
}
