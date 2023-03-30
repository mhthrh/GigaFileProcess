package entity

import "github.com/google/uuid"

type FileRequest struct {
	ID       uuid.UUID `json:"id"`
	Md5      string    `json:"Md5"`
	Name     string    `json:"name"`
	Count    int       `json:"count"`
	Sum      float64   `json:"sum"`
	Priority int       `json:"priority"`
}

type FileResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
}
