package main

import (
	"fmt"
	"net/http"

	"github.com/pp861/httptrace"

	"github.com/opentracing/opentracing-go"
	"github.com/pp861/go-stdlib/nethttp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func main() {
	// 1. init trace
	tracer, closer, _ := httptrace.InitTrace("gin-trace", "localhost:6831")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	// 2. use trace middleware
	mw := nethttp.Middleware(tracer, mux)

	http.ListenAndServe(":8000", mw)
}
