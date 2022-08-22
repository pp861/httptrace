package main

import (
	"io/ioutil"
	"net/http"

	"github.com/pp861/httptrace"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pp861/go-gin/ginhttp"
)

func main() {
	// 1. init trace
	tracer, closer, _ := httptrace.InitTrace("gin-trace", "localhost:6831")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	r := gin.Default()

	// 2. use trace middleware
	// spanFilterFn := func(r *http.Request) bool {
	// 	ctx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	// 	return true
	// }
	r.Use(ginhttp.Middleware(tracer /*, ginhttp.MWSpanFilter(spanFilterFn)*/))

	fn := func(c *gin.Context) {
		client := &http.Client{}
		url := "http://localhost:8000/hello"
		req, _ := http.NewRequest("GET", url, nil)

		// 3. use trace http client
		rsp, err := httptrace.DoHttpSend(c.Request.Context(), client, req)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"errno":  1,
				"errmsg": "do http call err",
				"data":   "",
			})

			return
		}

		data, _ := ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()

		c.JSON(http.StatusOK, gin.H{
			"errno":  0,
			"errmsg": "success",
			"data":   data,
		})
	}

	r.GET("/ping", fn)
	r.Run(":8001")
}
