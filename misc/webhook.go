package misc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DiscordWH struct {
	Contents      string `json:"content"`
	Username      string `json:"username"`
	AvatarURL     string `json:"avatarurl"`
	AllowMentions bool   `json:"allowed_mentions"`
}

func DiscordWebHook(id string, token string, contents string, username string, avatarurl string, allowmentions bool) bool {
	var data DiscordWH

	// Build out JSON/form data for Discord web hook.
	data.Contents = contents
	data.Username = username
	data.AvatarURL = avatarurl
	data.AllowMentions = allowmentions

	// Convert interface to JSON data string.
	datajson, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Build endpoint.
	urlstr := "https://discordapp.com/api/webhooks/" + id + "/" + token

	// Setup HTTP POST request with form data.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(datajson))

	// Set content type to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	resp.Body.Close()

	return true
}