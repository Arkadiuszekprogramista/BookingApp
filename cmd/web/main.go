package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/arkadiuszekprogramista/bookingapp/internal/config"
	"github.com/arkadiuszekprogramista/bookingapp/internal/handlers"
	"github.com/arkadiuszekprogramista/bookingapp/internal/helpers"
	"github.com/arkadiuszekprogramista/bookingapp/internal/models"
	"github.com/arkadiuszekprogramista/bookingapp/internal/render"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger


// main is the main application function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	fmt.Printf("Starting application on port %s \n", portNumber)


	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}
	err =  srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
		//what am i going to put in the session
		gob.Register(models.Reservation{})
		// change this to true when in production
		app.InProduction = false
	
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.InfoLog = infoLog

		errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog = errorLog

		// Set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
	
		app.Session = session
	
		tc, err := render.CreateTemplateCache()
		if err != nil {
			log.Fatal("cannot create template cache")
			return err
		}
	
		app.TemplateCache = tc
		app.UseCache = false
	
		repo := handlers.NewRepo(&app)
		handlers.NewHandlers(repo)
		render.NewTemplates(&app)
		helpers.NewHelpers(&app)

	return nil
}