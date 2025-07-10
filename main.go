package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
type sponsor struct{
	Company string
	Logo string
	Link string
	Invert int	
}
type language struct {
    NavAbout          string 
    NavMembers        string
    NavSponsors       string 
    NavLive           string 

    HeaderParagraph   string 
    About             string
    AboutParagraph    string
    AboutNameOrigin   string

    Members           string
    Contakt           string
    Footer            string
}
var db *sql.DB

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
func loadSponsors() []sponsor{
	_sponsors, err := os.ReadFile("data/sponsors.json")
	if err != nil{
		log.Fatal("error trying to read sponsors.json: ", err)
	}
	var sponsors []sponsor
	err = json.Unmarshal(_sponsors, &sponsors)
	if err != nil{
		log.Fatal("func error json.Unmarshal(_sponsors, &sponsors): ", err)
	}	
	return sponsors
}
func loadLanguage() []language{
	_language, err := os.ReadFile("data/language.json")
	if err != nil {
		log.Fatal("error trying to read languages.json")
	}
	var languages []language
	err = json.Unmarshal(_language, &languages)
	if err != nil{
		log.Fatal("func error json.Unmarshal(_language, &language)", err)
	}
	return languages
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
		log.Fatal("error executing the live template: ", err)
	}
}
func sponsorsHandler(w http.ResponseWriter, r *http.Request){
  sponsorsPath := "web/routes/sponsors.html"
  sponsors := template.Must(template.ParseFiles(sponsorsPath))

  err := sponsors.Execute(w, loadSponsors())
  if err != nil{
    log.Fatal("error executing the sponsors template: ", err)
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

func loadDB() string{
	file, err := os.Open("dbpwd")
	if err != nil{
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan(){
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil{
		fmt.Println("Error reading from dbpwd file: ", err)
	}
	return "" 
}

func dataHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	sensorType := vars["type"]
	scalarSensors := map[string]bool{
		"temperature": true,
		"humidity": true,
		"pressure": true,
	}
	vectorSensors := map[string]bool{
		"acceleration": true,
		"gyroscope": true,
		"gps": true,
	}
	w.Header().Set("Content-Type", "application/json")
	if scalarSensors[sensorType]{
		query := `
		SELECT s.value, s.timestamp FROM sensor_data s
		JOIN sensor_types t ON s.sensor_type_id = t.id
		WHERE t.name = ?
		ORDER BY s.timestamp DESC
		LIMIT 10
		`
		if db == nil {
			fmt.Println("Database connection is nil!")
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(query, sensorType)
		if err != nil{
			http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
			fmt.Println("query error"+err.Error())
			return
		}
		defer rows.Close()

		var data []map[string]interface{}
		for rows.Next(){
			var value float64
			var timestamp time.Time
			rows.Scan(&value, &timestamp)
			data = append(data, map[string]interface{}{
				"value": value,
				"timestamp": timestamp.Format("15:04:05"),
			})
		}
		json.NewEncoder(w).Encode(data)
	}else if vectorSensors[sensorType]{
		query := `
		SELECT v.x, v.y, v.z, v.timestamp
		FROM vector_data v
		JOIN sensor_types t ON v.sensor_type_id = t.id
		WHERE t.name = ?
		ORDER BY v.timestamp DESC
		LIMIT 10;
		`
		rows, err := db.Query(query, sensorType)
		if err != nil {
			http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []map[string]interface{}
		for rows.Next() {
			var x, y, z float64
			var timestamp time.Time
			rows.Scan(&x, &y, &z, &timestamp)
			data = append(data, map[string]interface{}{
				"x":         x,
				"y":         y,
				"z":         z,
				"timestamp": timestamp,
			})
		}
		json.NewEncoder(w).Encode(data)
	} else{
		fmt.Println("unknown type: ", sensorType)
		http.Error(w, "Unknown sensor type: ", http.StatusBadRequest)
	}
}

func main() {
	logFile, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening the log file: ", err)
	}

	defer logFile.Close()
	log.SetOutput(logFile)
	dsn := fmt.Sprintf("wtf:%s@tcp(localhost:3306)/PARADA?parseTime=true", loadDB())
	db, err = sql.Open("mysql",dsn)

	if err != nil{
		log.Fatal("error opening the db: ", err)
	}

	log.Printf("PARADA site v%d starting...", ver)


	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	r.PathPrefix("/modules/").Handler(http.StripPrefix("/modules/", http.FileServer(http.Dir("web/modules"))))
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/live", livePreviewHandler)
	r.HandleFunc("/sponsors", sponsorsHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.HandleFunc("/api/data/{type}", dataHandler);
	server := &http.Server{
		Handler: r,
		Addr:    ":80",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
