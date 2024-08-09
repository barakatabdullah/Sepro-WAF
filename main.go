package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/corazawaf/coraza/v3"
	txhttp "github.com/corazawaf/coraza/v3/http"
	"github.com/corazawaf/coraza/v3/types"
)

func initProxy() *httputil.ReverseProxy {
	// Create a reverse proxy
	targetURL, err := url.Parse("http://localhost:4000") // Replace with your target URL
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return proxy
}

func handler(w http.ResponseWriter, r *http.Request) {
	proxy := initProxy()
	proxy.ServeHTTP(w, r)
}

func createWAF() coraza.WAF {
	directivesFile := "./default.conf"
	if s := os.Getenv("DIRECTIVES_FILE"); s != "" {
		directivesFile = s
	}

	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithErrorCallback(logError).
			WithDirectivesFromFile(directivesFile),
	)
	if err != nil {
		log.Fatal(err)
	}
	return waf
}

func logError(error types.MatchedRule) {
	msg := error.ErrorLog()
	fmt.Printf("[logError][%s] %s\n", error.Rule().Severity(), msg)
}

func main() {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Create a new WAF instance
	waf := createWAF()

	// Start the server

	// Wrap the handler with the WAF with the default http package
	// http.Handle("/", txhttp.WrapHandler(waf, http.HandlerFunc(handler)))

	// Wrap the handler with the WAF with the gin package
	router.Any("/*path", func(c *gin.Context) {
		txhttp.WrapHandler(waf, http.HandlerFunc(handler)).ServeHTTP(c.Writer, c.Request)
	})

	// Start the server with the gin package
	router.Run()

	log.Println("Reverse proxy server started on port 8080")

	// Start the server with the default http package
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
