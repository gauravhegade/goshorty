package models

import "time"

type URLData struct {
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
	// title is optional
	// empty value exists therefore omitempty will not include the key at all
	Title     string    `json:"title,omitempty"`
	CreatedOn time.Time `json:"created_on"`
	// expires on is also optional
	// empty/nil value does not exist for time type,
	//	therefore I am using struct pointers here which will have a nil value
	ExpiresOn *time.Time `json:"expires_on"`
}
