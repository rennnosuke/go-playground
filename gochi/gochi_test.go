package gochi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestGoChi_Routing(t *testing.T) {
	// init router
	r := chi.NewRouter()

	// add router to middleware: Logger
	r.Use(middleware.Logger)

	// define routing (pass http.HandleFunc)
	r.Get("/hoge", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("hello, hoge")); err != nil {
			log.Fatalln(err)
		}
	})

	go func() {
		time.Sleep(time.Second)

		resp, err := http.Get("http://localhost:3000/hoge")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(bytes))
	}()

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}

}

func TestGoChi_RESTRouting(t *testing.T) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)                 // set unique requestID to each requests. we can trace each request more easily.
	r.Use(middleware.Logger)                    // add logger.
	r.Use(middleware.RealIP)                    // set the result of forwarding to X-Forwarded-For or X-Real-IP. we should only use when we put http proxy forward this app or trust client.
	r.Use(middleware.Recoverer)                 // recover when panic.
	r.Use(middleware.Timeout(60 * time.Second)) // set timeout.

	// unit route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello."))
	})

	// routing
	r.Route("/hoge", func(r chi.Router) {
		// add middleware to inner handler.
		r.With(middleware.Logger).Get("/piyo", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("get piyo")
		})

		r.Route("/fuga", func(r chi.Router) {
			// add middleware to inner handler
			r.Use(middleware.Logger)
			r.Get("/{fugaID}", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println()
			})
		})
	})

	// mount sub router
	anotherRouter := chi.NewRouter()
	r.Mount("/another", anotherRouter)

	http.ListenAndServe(":3333", r)
}
