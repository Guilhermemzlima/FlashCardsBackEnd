package main

import (
	log2 "FlashCardsBackEnd/internal/config/log"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

//func main() {
//	// Hello world, the web server
//
//	helloHandler := func(w http.ResponseWriter, req *http.Request) {
//		io.WriteString(w, "Hello, world!\n")
//	}
//
//	http.HandleFunc("/hello", helloHandler)
//	log.Println("Listing for requests at http://localhost:8000/hello")
//	log.Fatal(http.ListenAndServe(":8000", nil))
//
//}

func init() {
	var err = godotenv.Load()
	if err != nil {
		panic("Error to start application: cannot load environment variables: " + err.Error())
	}
	log2.SetupLogger()
}

func main() {
	application, err := SetupApplication()
	if err != nil {
		log2.Logger.Fatal("Setup Application Error", err)
		panic("Error to start application")
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		port = "8080"
	}

	handler := application.SystemRoutes.SetupHandler()
	server := &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//log.Logger.Infof("Application listening on port %s", port)
	log2.Logger.Fatal(server.ListenAndServe())
}
