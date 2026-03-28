package util

import (
	"log"
)

// type JsonError struct {
// 	Error string `json:"error"`
// }

// func NewJsonError(err error) JsonError {
// 	jsonError := JsonError{"generic error"}
// 	if err != nil {
// 		jsonError.Error = err.Error()
// 	}

// 	return jsonError
// }

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
