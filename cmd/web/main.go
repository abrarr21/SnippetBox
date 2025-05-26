package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abrarr21/snippet/pkg/mysql"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *scs.SessionManager
	users         *mysql.UserModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Define a command-line flag "addr" with default value ":6969" and a help message.
	// It returns a pointer to a string variable that will hold the value.
	addr := flag.String("addr", ":6969", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL DSN")

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

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new tamplate cache....
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a session manager.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	// Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db}, // (create a new SnippetModel struct and set its DB field to your DB connection and & ->returns a pointer to that struct )
		templateCache: templateCache,
		session:       sessionManager, // And add the session manager to our application dependencies.
		users:         &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// *addr dereferences the flag pointer to get the actual address string (e.g., ":6969")
	// It's passed to Printf to log the address the server is running on.
	infoLog.Printf("Server started running on port address: %s", *addr)

	// *addr is also passed here to start the HTTP server on the specified address.
	// The mux router is used to handle incoming HTTP requests.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	// start server: go run cmd/web/* -addr=":9999"
}
