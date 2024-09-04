package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"tolling/aggregator/client"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", "6000", "the listen address of the HTTP server")
	flag.Parse()

	var (
		client     = client.NewHTTPClient("http://localhost:4000")
		invHandler = newInvoiceHandler(client)
	)

	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))

	logrus.Infof("HTTP Gateway listen on port :%s", *listenAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *listenAddr), nil))
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: c}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), 162444750)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("REQ :: ")
		}(time.Now())

		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
