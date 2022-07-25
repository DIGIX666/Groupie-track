package handlers

import (
	"encoding/json"
	"errors"
	"groupie-tracker/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var Artists = []models.ArtistInfo{}
var Relation = &models.Relation{}

/*Rendu de la page d'accueil*/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != "GET" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	writeTemplate(w, "./templates/index.html", Artists)
}

/*Rendu de la page artist*/
func ArtistIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	id, err := extractQueryID(w, r) // Récuperation de L'ID
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	data, err := getArtistByID(id) // Récuperation de l'artiste en fonction de l'ID
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	writeTemplate(w, "./templates/artist.html", data) // Envoi des data sur la page artists
}

// Rendu page 404 not found
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		writeTemplate(w, "./templates/404.html", "404 Not Found")
	}
}

// Rendu general des template
func writeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}

/*Récuperation de l'ID pour l'url
 */
func extractQueryID(w http.ResponseWriter, r *http.Request) (int, error) {
	keys, ok := r.URL.Query()["ID"]
	if !ok || len(keys) != 1 {
		return 0, errors.New("❌ Url Param 'ID' is missing")
	}
	key := keys[0]
	id, err := strconv.Atoi(key)
	if err != nil {
		return 0, err
	}

	return id, nil
}

/*Récuperation data de l'API*/
func SyncData(api string, data interface{}) {
	log.Println("⌛️ Started syncronysation API" + api)
	res, err := http.Get(api)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = json.Unmarshal(bodyBytes, &data)

		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	log.Println("✅ Completed synchronization api" + api)
}

func filterArtistByID(artists []models.ArtistInfo, id int) *models.ArtistInfo {
	for _, item := range artists {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func filterRelationByID(relations []models.ArtistRelation, id int) *models.ArtistRelation {
	for _, item := range relations {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

/*Cherche un match entre ID pour recuperer les data des dates et concert*/
func getArtistByID(id int) (*models.ArtistData, error) {
	artist := filterArtistByID(Artists, id)
	if artist == nil {
		return nil, errors.New("Artist Not Found")
	}

	var data = models.ArtistData{Artist: *artist}
	var dates = filterRelationByID(Relation.Index, id)

	if dates != nil {
		data.DatesLocations = make(map[string]interface{})
		for key, value := range dates.DatesLocations {
			var locationName = strings.ReplaceAll(key, "_", " ")
			locationName = strings.ReplaceAll(locationName, "-", " - ")
			locationName = strings.ToUpper(locationName)
			data.DatesLocations[locationName] = value
		}
	}

	return &data, nil
}
