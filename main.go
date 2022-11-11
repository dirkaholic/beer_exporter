// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "beer"

var (
	listenAddress = flag.String("web.listen-address", ":9141",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")
	configPath = flag.String("config.file-path", "",
		"Path to environment file")

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last beer query successful.",
		nil, nil,
	)

	beersConsumed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "beers_consumed_total"),
		"How many beers have been consumed (per person).",
		[]string{"type", "person"}, nil,
	)
)

type BeerExporter struct {
}

func NewExporter() *BeerExporter {
	return &BeerExporter{}
}

func (e *BeerExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- beersConsumed
}

func (e *BeerExporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	e.UpdateMetrics(ch)
}

func (e *BeerExporter) UpdateMetrics(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		beersConsumed, prometheus.GaugeValue, float64(2), "schwarzbier", "mike",
	)

	log.Println("Endpoint scraped")
}

func main() {
	flag.Parse()

	configFile := *configPath
	if configFile != "" {
		log.Printf("Loading %s env file.\n", configFile)
		err := godotenv.Load(configFile)
		if err != nil {
			log.Printf("Error loading %s env file.\n", configFile)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file, assume env variables are set.")
		}
	}

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Beer Exporter</title></head>
             <body>
             <h1>Beer Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
