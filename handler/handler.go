package handler

import (
	"fmt"

	"github.com/qasim-sajid/clockify-api/dbhandler"
)

//Handler defines the handler struct for APIs
type Handler struct {
	DB dbhandler.DbHandler
}

//NewHandler implements constructor for Handler
func NewHandler() (*Handler, error) {
	dbC, err := dbhandler.NewDBClient("DBClient")
	if err != nil {
		return nil, fmt.Errorf("NewHandler: %v", err)
	}

	return &Handler{
		DB: dbC,
	}, nil
}
