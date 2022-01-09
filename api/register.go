package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"encoding/json"
)

type RegisterUsernameRequest struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
	Email     string `json:"email"`
}

func (client *Client) RegisterNewUsername(username string, email string, publicKey string) *RequestError {
	url := strings.Join([]string{client.Address, "user"}, "/")
	req := &RegisterUsernameRequest{
		Username:  username,
		PublicKey: publicKey,
		Email:     email,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return &RequestError{500, err}
	}

	reqBody := bytes.NewBuffer(reqJson)
	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		return &RequestError{500, err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return &RequestError{resp.StatusCode, err}
		}

		return &RequestError{resp.StatusCode, errors.New(string(bodyBytes))}
	}
	return nil
}
