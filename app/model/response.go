package model

// WebResponse adalah envelope yang sama untuk semua respons API.
type WebResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Meta    *Meta             `json:"meta,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Meta berisi informasi paginasi untuk respons daftar.
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListQuery adalah hasil pembacaan query string untuk endpoint daftar.
type ListQuery struct {
	Page   int
	Limit  int
	Search string
	Status string
	Sort   string
	Order  string
}
