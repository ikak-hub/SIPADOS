package model

const (
	DocumentCV         = "cv"
	DocumentTranscript = "transkrip"
	DocumentOther      = "lainnya"
)

type ApplicationDocument struct {
	ID            int    `json:"id"`
	ApplicationID int    `json:"application_id"`
	Type          string `json:"type"`
	FileURL       string `json:"file_url"`
}
