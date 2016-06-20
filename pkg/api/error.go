package api

import "fmt"

// OnlineError represents the structure returned by the Online API when an error occurred
type OnlineError struct {
	Why        string `json:"error"`
	Code       int    `json:"code"`
	StatusCode int    `json:"-"`
}

func (o OnlineError) Error() string {
	return fmt.Sprintf("[%v] : %v", o.StatusCode, o.Why)
}
