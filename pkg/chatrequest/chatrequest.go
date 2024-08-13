package chatrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
)

func Request(token string, msg string) (string, error) {
	if msg != "" {
		client := &http.Client{}
		u := url.URL{Scheme: "https", Host: "api.openai.com", Path: "/v1/chat/completions"}
		chatOptions := []string{
			"Respond with bardic insults",
			"Respond like a court jester",
			"Respond with sardonic wit",
			"Respond like a mystic",
		}
		i := rand.Intn(len(chatOptions) - 1)
		fmt.Println("Chat option", i)
		option := chatOptions[i]
		payload := map[string]interface{}{
			"model": "gpt-4o-mini",
			"messages": []map[string]interface{}{
				{"role": "system",
					"content": option},
				{"role": "user",
					"content": msg},
			},
		}
		pBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(pBytes))
		fmt.Println(string(pBytes))

		if err != nil {
			fmt.Println("Request failed", err)
			return "", err
		}

		req.Header = http.Header{}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		fmt.Println(req.Header)
		res, err := client.Do(req)

		if err != nil {
			fmt.Println("Do req failed", err)
			return "", err
		}
		defer res.Body.Close()

		bBytes, _ := io.ReadAll(res.Body)
		jsonBody := map[string]interface{}{}
		if err = json.Unmarshal(bBytes, &jsonBody); err != nil {
			return "", nil
		}
		fmt.Println(jsonBody)

		choices := jsonBody["choices"].([]interface{})
		fmt.Println("Response option", 0)
		choice := choices[0].(map[string]interface{})
		msg := choice["message"].(map[string]interface{})
		content := msg["content"].(string)
		fmt.Println(content)
		return content, nil

	}
	return "", nil

}
