package main

type Config struct {
	AppListenAddress string `envconfig:"BEER_APPLISTENADDRESS" default:":9141"`
	AppMetricsPath   string `envconfig:"BEER_APPMETRICSPATH" default:"/metrics"`
	// Database
	DbHost     string `envconfig:"BEER_DBHOST"`
	DbPort     int    `envconfig:"BEER_DBPORT" default:"5432"`
	DbUser     string `envconfig:"BEER_DBUSER" default:"api"`
	DbPassword string `envconfig:"BEER_DBPASSWORD"`
	DbDatabase string `envconfig:"BEER_DBDATABASE" default:"beer_exporter"`
}
