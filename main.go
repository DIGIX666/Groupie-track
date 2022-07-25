package main

import (
	"groupie-tracker/handlers"
	"log"
	"net/http"
)

func main() {
	go handlers.SyncData("https://groupietrackers.herokuapp.com/api/artists", &handlers.Artists)
	go handlers.SyncData("https://groupietrackers.herokuapp.com/api/relation", handlers.Relation)

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/artists", handlers.ArtistIndexHandler)

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	var port = ":8888"
	//log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Println("🖥  Server launched at the adress localhost:8888")
	http.ListenAndServe(port, nil)
}
