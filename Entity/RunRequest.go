package Entity

import "github.com/google/uuid"

type RunRequest struct {
	ID       uuid.UUID `json:"id"`
	Md5      string    `json:"Md5"`
	FileName string    `json:"file_name"`
	Count    int       `json:"count"`
	Priority int       `json:"priority"`
}
