package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Define a command-line flag "addr" with default value ":6969" and a help message.
	// It returns a pointer to a string variable that will hold the value.
	addr := flag.String("addr", ":6969", "HTTP network address")

	// Parse the command-line flags provided by the user (e.g. -addr=":8080")
	flag.Parse()

	// Creates a new logger that writes to standard output (os.Stdout).
	// Every log message will start with "INFO\t" prefix.
	// It will include the current date and time in each log (useful for timestamping).
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Creates a logger for errors, writing to standard error (os.Stderr).
	// Messages are prefixed with "ERROR\t" and also include date, time, and the file name + line number where the log was called.
	// log.Lshortfile helps in quickly locating the source of an error in your code.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves the files out of the "./ui/static/" directory.
	// Note that the path given to the http.Dir function is relative to the project's directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with /static/. For matching paths, we strip the "/static" prefix before the request reaches the server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// *addr dereferences the flag pointer to get the actual address string (e.g., ":6969")
	// It's passed to Printf to log the address the server is running on.
	infoLog.Printf("Server started running on port address: %s", *addr)

	// *addr is also passed here to start the HTTP server on the specified address.
	// The mux router is used to handle incoming HTTP requests.
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

	// start server: go run cmd/web/* -addr=":9999"
}
