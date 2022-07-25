package models

type ArtistInfo struct {
	ID           int
	Name         string
	CreationDate int
	FirstAlbum   string
	Image        string
	Members      []string
}

type Relation struct {
	Index []ArtistRelation
}

type ArtistRelation struct {
	ID             int
	DatesLocations map[string]interface{}
}

type ArtistData struct {
	Artist         ArtistInfo
	DatesLocations map[string]interface{}
}
