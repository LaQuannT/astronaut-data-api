package model

type JSONResponse struct {
	Astronaut  *Astronaut   `json:"astronaut,omitempty"`
	Astronauts []*Astronaut `json:"astronauts,omitempty"`
	User       *User        `json:"user,omitempty"`
	Users      []*User      `json:"users,omitempty"`
	Message    string       `json:"message,omitempty"`
	Error      string       `json:"error,omitempty"`
	Errors     []string     `json:"errors,omitempty"`
}

type ApiError struct{}

func (e *ApiError) Error() string {
	return ""
}
