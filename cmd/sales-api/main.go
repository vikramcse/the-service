package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/vikramcse/the-service/cmd/sales-api/internal/handlers"
	"github.com/vikramcse/the-service/internal/platform/database"
)

func main() {

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:0.0.0.0:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatal("error: generating conf usage: %v", err)
			}
			fmt.Println(usage)
			return
		}

		log.Fatalf("error: parsing config: %s", err)
	}

	log.Println("main: Started")
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("error: generating conf for output: %v", err)
	}
	log.Printf("main: Config: \n%v\n", out)

	// Start Database
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatalf("error: connection to db: %s", err)
	}
	defer db.Close()

	productsHandler := handlers.Products{DB: db}

	// Api service configuration

	// ReadTimeout: It defines how long you allow a connection to be open
	// during a client sends data or max time required for reading the entire
	// request, includig body

	// WriteTimeout: It is maximum duration before timing out writes of the
	// response.
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(productsHandler.List),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listend for errors coming from the listener. Use a
	// Buffered channel so the goroutine can exit if we do not collect error
	serverErrors := make(chan error, 1)

	// Start the service listening for requests
	// the reason for adding ListenAndServe in goroutine because ListenAndServe
	// is a blocking call and if we want to go one doing more work, like
	// making a second instance if ListenAndServe, then we need a separate
	// goroutine
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminal signal from OS
	// Use a buffered channel, as signal package requires it
	// Note: SIGTERM signal is used by kubernetes instead of os.Interrupt
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)
	case <-shutdown:
		log.Println("main: Starting shutdown")

		// Added a deadline for request completion
		// we can perfrom any chores in this time e.g clearing memory,
		// resources etc.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// SetKeepAlivesEnabled will inform the webserver to not keep any
		// existing connections alive which basically gives us the gracefull
		// shutdown behavior
		api.SetKeepAlivesEnabled(false)

		// Asking a listener to shutdown without interrupting any active
		// connections. Shutdown works by first closing all open listeners,
		// then closing idle connections and then waiting indefinitely for
		// connections to return to idle and then shutdown. If the provided
		// context expires before Shutdown is complete, shutdown returns
		// the context error, otherwise it returns any error returned from
		// closign the servers listeners
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main: Graceful shutdown did not complete in %v: %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main: could not stop server gracefully: %v", err)
		}
	}
}

type Product struct {
	Name     string `json:"name"`
	Cost     int    `json: "cost"`
	Quantity int    `json: "quantity"`
}

// ListProducts is an HTTP Handler for returning a list of Products.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := []Product{
		{Name: "Comic Books", Cost: 50, Quantity: 42},
		{Name: "McDonalds Toys", Cost: 75, Quantity: 120},
	}

	data, err := json.Marshal(list)
	if err != nil {
		log.Println("error marshaling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing resutl", err)
	}
}
