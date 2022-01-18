package api

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"crypto/rsa"
	"encoding/json"
	"encoding/base64"

	r "github.com/marcinbor85/nes/crypto/rsa"
)

type GetPublicKeyResponse struct {
	PublicKeyMessage string `json:"public_key_message"`
	PublicKeySign 	 string `json:"public_key_sign"`
}

type GetPublicKeyData struct {
	Response 	json.RawMessage `json:"response"`
	Signature	string `json:"signature"`
}

func (client *Client) GetPublicKeyByUsername(username string) (*rsa.PublicKey, *rsa.PublicKey, error) {

	cachedKeys := (*client.PublicKeyCache)[username]
	if cachedKeys != nil {
		return cachedKeys.PublicKeyMessage, cachedKeys.PublicKeySign, nil
	}

	url := strings.Join([]string{client.Address, "users", username}, "/")
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}

		return nil, nil, errors.New("user not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	data := &GetPublicKeyData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, nil, err
	}

	response := &GetPublicKeyResponse{}
	err = json.Unmarshal([]byte(data.Response), &response)
	if err != nil {
		return nil, nil, err
	}

	keyMessage, err := r.DecodePublicKey(response.PublicKeyMessage)
	if err != nil {
		return nil, nil, err
	}

	keySign, err := r.DecodePublicKey(response.PublicKeySign)
	if err != nil {
		return nil, nil, err
	}

	responseBin := []byte(data.Response)
	signatureBin, err := base64.URLEncoding.DecodeString(data.Signature)
	if err != nil {
		return nil, nil, err
	}

	err = r.Verify(responseBin, signatureBin, client.ServerPublicKey)
	if err != nil {
		return nil, nil, err
	}

	(*client.PublicKeyCache)[username] = &KeyCache{
		PublicKeyMessage: keyMessage,
		PublicKeySign: keySign,
	}

	return keyMessage, keySign, nil
}
