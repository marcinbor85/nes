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
	Username			string `json:"username"`
	PublicKeyMessage	string `json:"public_key_message"`
	PublicKeySign		string `json:"public_key_sign"`
	Email				string `json:"email"`
}

func (client *Client) RegisterNewUsername(username string, email string, publicKeyMessage string, publicKeySign string) error {
	url := strings.Join([]string{client.Address, "users"}, "/")
	req := &RegisterUsernameRequest{
		Username:  username,
		PublicKeyMessage: publicKeyMessage,
		PublicKeySign: publicKeySign,
		Email:     email,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}

	reqBody := bytes.NewBuffer(reqJson)
	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(strings.Trim(string(bodyBytes), "\r\n"))
	}
	return nil
}
