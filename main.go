package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Filmi struct {
	Id          string
	Emri        string
	Studio      string
	Kohezgjatja string
	Kategoria   string
	Viti        int
	Vleresimi   float32
}

var db, _ = sql.Open("mysql", "root:urankajtazaj@/testdb")

var t = template.Must(template.ParseGlob("tmpl/*.html"))

func getData(w http.ResponseWriter) {
	rs, _ := db.Query("select * from filmat")

	var fslice = make([]Filmi, 0)
	var flm Filmi
	for rs.Next() {
		rs.Scan(&flm.Id, &flm.Emri, &flm.Studio, &flm.Kohezgjatja, &flm.Vleresimi, &flm.Viti, &flm.Kategoria)
		fslice = append(fslice, flm)
	}
	t.ExecuteTemplate(w, "index.html", fslice)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	getData(w)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "shto.html", nil)
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/"):]
	rs, _ := db.Query("select * from filmat where id = " + id)

	var film Filmi
	for rs.Next() {
		rs.Scan(&film.Id, &film.Emri, &film.Studio, &film.Kohezgjatja, &film.Vleresimi, &film.Viti, &film.Kategoria)
	}
	t.ExecuteTemplate(w, "edit.html", film)
}

func handleDb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	emri := r.FormValue("emri")
	studio := r.FormValue("studio")
	kohezgjatja := r.FormValue("kohezgjatja")
	kategoria := r.FormValue("kategoria")
	viti := r.FormValue("viti")
	vleresimi, _ := strconv.ParseFloat(r.FormValue("vleresimi"), 32)

	stm, _ := db.Prepare("insert into filmat values (null, ?, ?, ?, ?, ?, ?)")
	_, err := stm.Exec(emri, studio, kohezgjatja, vleresimi, viti, kategoria)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/delete/"):]
	stm, _ := db.Prepare("delete from filmat where id = ?")
	_, err := stm.Exec(id)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

func handleEditUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/updating/"):]
	stm, _ := db.Prepare("update filmat set kategoria = ?, viti = ?, emri = ?, studio = ?, viti = ?, kohezgjatja = ?, vleresimi = ? where id = ?")
	r.ParseForm()
	vleresimi, _ := strconv.ParseFloat(r.FormValue("vleresimi"), 32)
	_, err := stm.Exec(r.FormValue("kategoria"), r.FormValue("viti"), r.FormValue("emri"), r.FormValue("studio"), r.FormValue("viti"), r.FormValue("kohezgjatja"), vleresimi, id)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("tmpl/css"))))

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/delete/", handleDelete)
	http.HandleFunc("/add/", handleAdd)
	http.HandleFunc("/add/adding/", handleDb)
	http.HandleFunc("/edit/", handleEdit)
	http.HandleFunc("/edit/updating/", handleEditUpdate)

	http.ListenAndServe(":8080", nil)
}
