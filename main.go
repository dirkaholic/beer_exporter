// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "beer"

var (
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

func NewBeerExporter() *BeerExporter {
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

const createPersons string = `
CREATE TABLE IF NOT EXISTS persons (
	username varchar(32) NOT NULL,
	fullname varchar(256) NOT NULL,
	PRIMARY KEY (username)
  );`

func main() {
	var c Config
	err := envconfig.Process("beer", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s",
		c.DbUser,
		c.DbPassword,
		c.DbHost,
		c.DbDatabase,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to database.")
	}

	_, execErr := db.Exec(createPersons)
	if execErr != nil {
		log.Fatalf("An error occured while executing query: %v", execErr)
	}

	exporter := NewBeerExporter()
	prometheus.MustRegister(exporter)

	http.Handle(c.AppMetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Beer Exporter</title></head>
             <body>
             <h1>Beer Exporter</h1>
             <p><a href='` + c.AppMetricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(c.AppListenAddress, nil))
}
