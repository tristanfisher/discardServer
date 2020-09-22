package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/tristanfisher/discardServer/logging"
)

type Feedback struct {
	log *zap.Logger
}

func returnListeningServer(addr string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	srv.SetKeepAlivesEnabled(true)
	return srv
}

func (f *Feedback) wildcardHandler(_ http.ResponseWriter, r *http.Request) {
	logFields := f.log.WithOptions(zap.Fields(zap.String("url", r.URL.String())))
	logFields.Info("received request")

	// just quickly discard and return OK
	_, _ = io.Copy(ioutil.Discard, r.Body)
	err := r.Body.Close()
	if err != nil {
		logFields.Error("failed to close body", zap.String("error", err.Error()))
	}
}

func main() {
	logLevelOptions := []string{"debug", "info", "warn", "error", "panic", "fatal"}
	logLevelPtr := flag.String("logLevel", "debug", "log level (options: debug, info, warn, error, panic, fatal)")
	listeningPortPtr := flag.String("listeningPort", "8080", "port to bind for the listening server")
	flag.Parse()

	// can we get our logger going?
	validLogLevel := false
	for _, lvl := range logLevelOptions {
		if lvl == *logLevelPtr {
			validLogLevel = true
			break
		}
	}
	if !validLogLevel {
		_, _ = fmt.Fprintf(os.Stderr, "invalid log level: %s", *logLevelPtr)
		flag.Usage()
		os.Exit(1)
	}
	log := logging.MustSetLevelLog(*logLevelPtr)

	// lazily format port request
	lp := fmt.Sprintf(":%s", *listeningPortPtr)

	f := &Feedback{
		log: log,
	}

	// http router that listens to anything and just returns a 200
	mux := http.NewServeMux()
	mux.HandleFunc("/", f.wildcardHandler)

	// setup server
	server := returnListeningServer(lp, mux)

	log.Debug("starting listening server", zap.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("http server failed", zap.String("error", err.Error()))
	}
}
