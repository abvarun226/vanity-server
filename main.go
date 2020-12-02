package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/abvarun226/vanity-server/handler"
)

func main() {

	var listenAddr, rulesFile, godocURL string
	var tlsCertFile, tlsKeyFile string

	// Flags/command-line arguments
	flag.StringVar(&listenAddr, "listen", ":8080", "Address where this server listens")
	flag.StringVar(&godocURL, "godocurl", "", "The godoc URL")
	flag.StringVar(&rulesFile, "config", "/etc/go-vanity/config.json", "Contains go-import mapping rules")
	flag.StringVar(&tlsCertFile, "config", "/etc/go-vanity/config.json", "Contains go-import mapping rules")
	flag.StringVar(&tlsKeyFile, "config", "/etc/go-vanity/config.json", "Contains go-import mapping rules")
	flag.Parse()

	// Make sure config file exists, or exit
	if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
		flag.PrintDefaults()
		log.Fatalf("failed to find mapping rules config file")
	}

	h := handler.New(
		handler.WithGodocURL(godocURL),
		handler.WithConfigFile(rulesFile),
	)

	// load the import path rules.
	h.GetImportRules()

	// HTTP routes.
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/reload", h.ReloadRules)
	http.HandleFunc("/", h.VanityServer)

	log.Printf("server listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
