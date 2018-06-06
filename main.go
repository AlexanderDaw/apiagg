package main

import (
	"github.com/AlexanderDaw/apiagg/services"
	"log"
	"github.com/gorilla/mux"
	"github.com/AlexanderDaw/apiagg/transport"
	"syscall"
	"fmt"
	"os/signal"
	"os"
	"net/http"
	"flag"
)

/**
Main entry point into the REST Agg program.
The goal of this program is to demonstrate golangs ability
to quickly enable concurrent execution of REST requests
 */
func main(){
	//Properties Parsing.

	var (
		httpAddr     = flag.String("http.addr", ":8000", "Address for HTTP (JSON) server")
		concurrencyCount  = flag.Int("concurrency", 3, "Default concurrency number")
		responseDelay = flag.Int("delay", 30, "The number of ms to wait for default response times")
	)

	AggService, err := services.InitializeAggregationService(*concurrencyCount)
	if err != nil {
		log.Fatal("Error creating aggregation service", err.Error())
	}

	//Test response service
	TestService, terr := services.InitializeVariableResponseTestService(*responseDelay)
	if terr != nil {
		log.Fatal("Error creating the rest response service", terr.Error())
	}

	//Router for handlers and server.
	r := mux.NewRouter()

	transport.RegisterRoutes(r, AggService, TestService)

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		log.Println("transport", "HTTP", "addr", *httpAddr)
		errc <- http.ListenAndServe(*httpAddr, r)
	}()

	// Run!
	log.Println("exit", <-errc)
}