package model

import "github.com/supasiti/prac-go-htmx-tailwind/internal/pkg/json"

type Contact struct {
	ContactID int    `json:"contact_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

func (c Contact) String() string {
	return json.ToJSONString(c)
}
