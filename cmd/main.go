package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// "otel-poc/common"
	"otel-poc/utils"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(30 * time.Millisecond)
	w.Write([]byte("Hello, world!"))
}
func main() {
	ctx := context.Background()
	gcpProject := "tilda-trial-dev"
	// fmt.Println(common.Version)
	fmt.Println(utils.ServiceName)
	tp, err := utils.InitTracing(ctx, utils.Config{
		ProjectID:      gcpProject,
		ServiceName:    utils.ServiceName,
		ServiceVersion: utils.Version,
	})

	if err != nil {
		fmt.Println(err)
	}
	defer tp.Shutdown(ctx)

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "/hello")
	http.Handle("/hello", otelHandler)
	http.ListenAndServe(":9000", nil)
}
