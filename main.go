package main

import (
  "fmt"
  "net/http"

  "github.com/nathanielwheeler/go-fullstack/controllers"
  "github.com/nathanielwheeler/go-fullstack/middleware"
  "github.com/nathanielwheeler/go-fullstack/models"
  "github.com/nathanielwheeler/go-fullstack/rand"

  "github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	// Load Configuration
	cfg := LoadConfig()
	dbCfg := cfg.Database

  // Initialize Services
  services, err := models.NewServices(
    models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionString()),
    models.WithLogMode(!cfg.IsProd()),
    models.WithUser(cfg.Pepper, cfg.HMACKey)
  )
  defer services.Close()
  services.AutoMigrate()

  // Router initialization
  r := mux.NewRouter()

  // Initialize controllers
  staticC := controllers.NewStatic()
  usersC := controllers.NewUsers(services.User)
  valuesC := controllers.NewValues(services.Values)

  // Middleware
	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{}

	// CSRF Protection
	b, err := rand.Bytes(cfg.CSRFBytes)
	if err != nil {
		panic(err)
	}
  csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

  // FileServer
  publicHandler := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/assets/").
		Handler(publicHandler)
  
  // Static Routes
  r.handle("/", staticC.Home).
    Methods("GET")

  // User Routes
  r.HandleFunc("/register",
    usersC.Registration).
    Methods("GET")
  r.HandleFunc("/register",
    usersC.Register).
    Methods("POST")
  r.Handle("/login",
    usersC.LoginView).
    Methods("GET")
  r.HandleFunc("/login",
    usersC.Login).
    Methods("POST")
  r.HandleFunc("/cookietest",
    usersC.CookieTest).
    Methods("GET")

  // Values Routes
  r.HandleFunc("/values").
    Methods("GET")

  // Start server
  port := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Now listening on %s...\n", port)
	http.ListenAndServe(port, csrfMw(userMw.Apply(r)))
}

