package main

import (
	"net/http"
	"snacomds/SeproWAF/initializers"

	"fmt"
	"log"

	txhttp "github.com/corazawaf/coraza/v3/http"

	"strings"

	"os"
)

func exampleHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	resBody := "Hello world, transaction not disrupted."

	if body := os.Getenv("RESPONSE_BODY"); body != "" {
		resBody = body
	}

	if h := os.Getenv("RESPONSE_HEADERS"); h != "" {
		key, val, _ := strings.Cut(h, ":")
		w.Header().Set(key, val)
	}

	// The server generates the response
	w.Write([]byte(resBody))
}

func main() {

	waf, err := initializers.CreateWaf()
	if err != nil {
		fmt.Println(err)
	}

	tx := waf.NewTransaction()
	tx.ProcessConnection("127.0.0.1", 8066, "127.0.0.1", 8000)
	http.Handle("/", txhttp.WrapHandler(waf, http.HandlerFunc(exampleHandler)))
	fmt.Println("Server is running. Listening port: 8066")

	log.Fatal(http.ListenAndServe(":8066", nil))

}
