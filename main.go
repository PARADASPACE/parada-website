package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const ver = 2

type member struct {
	Name  string
	Pfp   string
	Roles string
	Bio   string
	socials
}
type socials struct {
	Linkedin string
}

func loadMembers() []member {
	_members, err := os.ReadFile("data/members.json")
	if err != nil {
		log.Fatal("error trying to read members.json: ", err)
	}
	var members []member
	err = json.Unmarshal(_members, &members)
	if err != nil {
		log.Fatal("func error json.Unmarshal(_members,&members): ", err)
	}
	return members
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := "web/routes/index.html"
	index := template.Must(template.ParseFiles(indexPath))

	err := index.Execute(w, loadMembers())
	if err != nil {
		log.Fatal("error executing the index template: ", err)
	}
}
func livePreviewHandler(w http.ResponseWriter, r *http.Request) {
	livePath := "web/routes/live.html"
	live := template.Must(template.ParseFiles(livePath))

	err := live.Execute(w, nil)
	if err != nil {
		log.Fatal("error executing the index template: ", err)
	}
}
func notFound(w http.ResponseWriter, r *http.Request) {
	notFoundPath := "web/routes/404.html"
	notFound := template.Must(template.ParseFiles(notFoundPath))
	err := notFound.Execute(w, nil)
	if err != nil {
		log.Fatal("error executing the notFound template: ", err)
	}
}
func main() {
	logFile, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening the log file: ", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Printf("PARADA site v%d starting...", ver)
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	r.PathPrefix("/modules/").Handler(http.StripPrefix("/modules/", http.FileServer(http.Dir("web/modules"))))
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/live", livePreviewHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	server := &http.Server{
		Handler: r,
		Addr:    ":80",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
