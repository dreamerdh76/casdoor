package object

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type EnzonaSMS struct {
	token          string
	consumerKey    string
	consumerSecret string
	auth           string
	endpoint       string
}

type NotificationsSms struct {
	Message string   `json:"message"`
	Phone   []string `json:"phone"`
}

func newEnzonaSmsClient(consumerKey, consumerSecret, auth, endpoint string) (*EnzonaSMS, error) {
	token, err := getTokenEnzonaSms(consumerKey, consumerSecret, auth)
	if err != nil {
		return nil, err
	}
	return &EnzonaSMS{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		auth:           auth,
		endpoint:       endpoint,
		token:          token,
	}, nil
}

// :return: Returns an API access token It depends on the input of the public and private keys.
func getTokenEnzonaSms(consumerKey, consumerSecret, endpoint string) (string, error) {
	auth := consumerKey + ":" + consumerSecret
	bs4 := base64.StdEncoding.EncodeToString([]byte(auth))
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	data := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", endpoint, data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+bs4)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("ERROR: Get Token failed with HTTP status code: " + strconv.Itoa(resp.StatusCode))
	}
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}
	return jsonResp["access_token"].(string), nil
}

func (e *EnzonaSMS) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	objectN := NotificationsSms{
		Message: param["code"],
		Phone:   targetPhoneNumber,
	}
	data, _ := json.Marshal(objectN)

	req, err := http.NewRequest("POST", e.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+e.token)
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("EnzonaSmsClient's SendMessage() error, Enzona SMS request failed with status: %s", resp.Status)
	}

	return nil
}
