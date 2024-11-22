package models

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"lyrics"`
	Link        string `json:"link"`
}
