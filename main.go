package main

import (
	"flag"
	"fmt"
	. "github.com/semka/hysteria/models"
	. "github.com/semka/hysteria/outputs"
	"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <OPTIONS> <PROJECT-PATH>\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	riemannFlag := flag.String("riemann", "localhost:5555", "Address of the Riemann instance")
	timeoutFlag := flag.Int("timeout", 30, "Timeout after which test will be considered failed")
	serviceName := flag.String("service", "hysteria/tests", "Service name for tests")

	flag.Parse()

	projectPath := flag.Arg(0)
	if projectPath == "" {
		usage()
		os.Exit(1)
	}

	proj := Project{Path: projectPath, TestTimeout: *timeoutFlag}
	suites, err := proj.GetTestSuites()
	if err != nil {
		panic(err)
	}

	for _, suite := range suites {
		results, _ := suite.Run()
		for _, res := range results {
			log.Printf(res.LogLine())
		}

		if riemannFlag != nil {
			r := Riemann{URL: *riemannFlag, ServiceName: *serviceName}
			err := r.Connect()
			if err != nil {
				log.Fatal(err)
			}
			r.SendBulk(results)
		}
	}
}
