package main

import (
	_ "github.com/snowflakedb/gosnowflake"
	"log"
	"snowlastic-cli/demo"
)

func main() {
	err := demo.IndexDemos(
		"./demo/demos.json",           //demosPath
		"./demo/demo-settings.json",   //demoSettings
		"./ignore/local_elastic.json", //credPath
		"./ignore/local_http_ca.crt",  //caCertPath
	)
	if err != nil {
		log.Fatal(err)
	}
}
