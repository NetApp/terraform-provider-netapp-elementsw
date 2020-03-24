package elementsw

import (
	"crypto/tls"
	"net/http"
)

// Config is a struct for user input
type configStuct struct {
	User            string
	Password        string
	ElementSwServer string
	APIVersion      string
}

// Client contain the api endpoint
type clientStuct struct {
	Endpoint string
}

// APIError is any error the api gives
type APIError struct {
	ID    int `json:"id"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"error"`
}

// Client is the main function to connect to the APi
func (c *configStuct) clientFun() (*Client, error) {
	client := &Client{
		Host:     "https://" + c.ElementSwServer,
		Username: c.User,
		Password: c.Password,
		HTTPTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true},
		},
	}

	client.SetAPIVersion(c.APIVersion)

	return client, nil
}
