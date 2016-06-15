package api

import "fmt"

type OnlineError struct {
	Why        string `json:"error"`
	Code       int    `json:"code"`
	StatusCode int    `json:"-"`
}

func (o OnlineError) Error() string {
	return fmt.Sprintf("[%v] : %v", o.StatusCode, o.Why)
}
