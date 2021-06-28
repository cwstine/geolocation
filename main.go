package main

/*
import "fmt"

func main() {
	fmt.Println("Hello bob")
}
*/

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//Add in CORS headers for Cross-Origin Requests
	addCorsHeader(w)
	//Need to respond to an OPTIONS request with 200 status and CORS headers
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Fprintf(os.Stdout, "%s %s %s \n", r.Method, r.URL, r.Proto)
	//Iterate over all header fields
	for k, v := range r.Header {
		fmt.Fprintf(os.Stdout, "Header field %q, Value %q\n", k, v)
	}

	fmt.Fprintf(os.Stdout, "Host = %q\n", r.Host)
	fmt.Fprintf(os.Stdout, "RemoteAddr= %q\n", r.RemoteAddr)
	renderTemplate(w, title, nil)
}

var tempPath = "./html/"

//var templates = template.Must(template.ParseFiles(tempPath+"contact.html", tempPath+"footer.html", tempPath+"resume.html",
//	tempPath+"home.html", tempPath+"about.html", tempPath+"posts.html", tempPath+"donate.html", tempPath+"shop.html", tempPath+"volunteer.html", tempPath+"whatwedo_events.html", tempPath+"whatwedo_mission.html",
//	tempPath+"whatwedo_programs.html", tempPath+"whatwedo.html"))

var templates = template.Must(template.ParseFiles(tempPath + "geolocation.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
   Add CORS headers to an http response object to allow for Cross-Origin requests made to this server
   Dont remember if using CORs relevent for geolocation
*/
func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

var validPath = regexp.MustCompile("^/(geolocation)")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func main() {
	http.HandleFunc("/geolocation", makeHandler(viewHandler))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
