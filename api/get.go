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
	PublicKey string `json:"public_key"`
}

func (client *Client) GetPublicKeyByUsername(username string) (*rsa.PublicKey, error) {

	cachedKey := (*client.PublicKeyCache)[username]
	if cachedKey != nil {
		return cachedKey, nil
	}

	url := strings.Join([]string{client.Address, "users", username}, "/")
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

	key, err := r.DecodePublicKey(data.PublicKey)
	if err != nil {
		return nil, err
	}

	(*client.PublicKeyCache)[username] = key

	return key, nil
}
