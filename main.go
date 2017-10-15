package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// estructura donde se guarda cada fila del archivo
type row struct {
	Catalogue, Area, Item string
}

func main() {

	ruta := mux.NewRouter()
	ruta.HandleFunc("/app/index", index)
	ruta.HandleFunc("/app/data", listData)
	ruta.HandleFunc("/app/updateForm/{row:[0-9]+}", updateForm)
	ruta.HandleFunc("/app/update/", update)
	ruta.HandleFunc("/app/delete/{row:[0-9]+}", delete)

	ruta.PathPrefix("/").Handler(http.FileServer(http.Dir("resources")))
	http.Handle("/", ruta)
	http.ListenAndServe(":8080", nil)

}

// funcion que lee archivo y regresa el Json
func fileJSON() ([]byte, error) {

	// se lee archivo
	stream, err := ioutil.ReadFile("Catalogo de Servicios.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(stream)))

	var rows []row

	for {
		record, err := r.Read()
		if err == io.EOF {
			break // detiene ciclo al final del archivo
		} else if err != nil {
			log.Fatalf("error al leer fila")
		}
		// guarda la fila (record) en una estructura (row) y la agrega al arreglo (rows)
		rows = append(rows, row{Catalogue: record[0], Area: record[1], Item: record[2]})
	}

	//se parsea el arreglo como Json
	return json.Marshal(rows) // se omite el primer elemento ya que es metadata

}

//funcion que escribe JSON al archivo CSV
func writeCSV(filedata []byte) {

	rows := make([]row, 0)
	json.Unmarshal(filedata, &rows)

	//TODO sobreescribir en vez de crear
	file, err := os.Create("Catalogo de Servicios.csv")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range rows {
		sval := r.Catalogue + "," + r.Area + "," + r.Item + "\n"
		file.WriteString(sval)
	}

}

//funcion que escribe JSON al archivo JSON
func writeJSON(filedata []byte, fileName string) {

	//TODO sobreescribir en vez de crear
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(string(filedata))

}

// Handlers
func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func listData(w http.ResponseWriter, r *http.Request) {

	json, err := fileJSON()

	if err != nil {
		log.Fatal(err)
	} else {
		writeJSON(json, "resources/JSON/outputJson.json")
		http.ServeFile(w, r, "views/data.html")
	}

}

func updateForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/update.html")
}

func update(w http.ResponseWriter, r *http.Request) {

	rowIndex, _ := strconv.Atoi(r.FormValue("index"))
	cata := r.FormValue("cata")
	area := r.FormValue("area")
	item := r.FormValue("item")

	dataArray := make([]row, 0)
	jsonData, err := fileJSON()
	json.Unmarshal(jsonData, &dataArray)

	dataArray[rowIndex].Catalogue = cata
	dataArray[rowIndex].Area = area
	dataArray[rowIndex].Item = item

	out, err := json.Marshal(dataArray)

	if err != nil {
		log.Fatal(err)
	} else {
		writeJSON(out, "resources/JSON/outputJson.json")
		writeCSV(out)
		http.ServeFile(w, r, "views/data.html")
	}

}

func delete(w http.ResponseWriter, r *http.Request) {

	jsonData, err := fileJSON()
	urlParams := mux.Vars(r)

	rowIndex, _ := strconv.Atoi(urlParams["row"])

	dataArray := make([]row, 0)
	json.Unmarshal(jsonData, &dataArray)

	//Se borra elemento para este caso no hay muchos datos por lo que no debe haber problema de rendimiento
	dataArray = append(dataArray[:rowIndex], dataArray[rowIndex+1:]...)
	out, err := json.Marshal(dataArray)

	if err != nil {
		log.Fatal(err)
	} else {
		writeJSON(out, "resources/JSON/outputJson.json")
		writeCSV(out)
		http.ServeFile(w, r, "views/data.html")

	}

}
