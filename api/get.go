package api

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"crypto/rsa"
	"encoding/json"

	r "github.com/marcinbor85/nes/crypto/rsa"
)

type GetPublicKeyResponse struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

func (client *Client) GetPublicKeyByUsername(username string) (*rsa.PublicKey, error) {
	url := strings.Join([]string{client.Address, "user", username}, "/")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("user not found")
	}

	var data GetPublicKeyResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if username != data.Username {
		return nil, errors.New("username mismatch")
	}

	key, err := r.DecodePublicKey(data.PublicKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}