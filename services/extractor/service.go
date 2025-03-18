package extractor

import (
	"context"
	"encoding/json"
	"fmt"
	"go-extractor/config"
	"go-extractor/helper"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	retries = 5
)

type Response struct {
	Data     []interface{} `json:"data"`
	NextPage struct {
		Offset string `json:"offset"`
		Path   string `json:"path"`
		Uri    string `json:"uri"`
	} `json:"next_page"`
}

type ResponseProject struct {
	Data     []Project `json:"data"`
	NextPage struct {
		Offset string `json:"offset"`
		Path   string `json:"path"`
		Uri    string `json:"uri"`
	} `json:"next_page"`
}
type ResponseUser struct {
	Data     []User `json:"data"`
	NextPage struct {
		Offset string `json:"offset"`
		Path   string `json:"path"`
		Uri    string `json:"uri"`
	} `json:"next_page"`
}

func Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	extract(ctx)

	for {
		select {
		case <-ticker.C:
			extract(ctx)
		case <-ctx.Done():
			ticker.Stop()
			log.Println("Stopping extractor..")
			return
		}
	}
}

func extract(ctx context.Context) {
	var responseProject ResponseProject
	response, err := fetch(ctx, "/projects")
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &responseProject)
	if err != nil {
		log.Println(err)
	}

	for _, project := range responseProject.Data {
		var responseUser ResponseUser

		response, err := fetch(ctx, "/users")

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(response, &responseUser)
		if err != nil {
			log.Println(err)
		}

		for _, user := range responseUser.Data {
			err = helper.SaveToJsonFile("users", user)
		}

		err = helper.SaveToJsonFile("project", project)

	}
}

func fetch(ctx context.Context, endpoint string) ([]byte, error) {
	var err error

	for i := 0; i < retries; i++ {
		data, err := fetchUrl(ctx, endpoint)

		if err == nil {
			return data, nil
		}

		//switch e := err.(type) {
		//case nil:
		//}

		log.Print(err)
		log.Printf("Retrying fetch (%d/%d):", i+1, retries)
		time.Sleep(1 * time.Second) // retry interval
	}

	log.Printf("Failed to fetch URL: %s data", endpoint)
	return nil, err
}

func fetchUrl(ctx context.Context, endpoint string) ([]byte, error) {
	client := &http.Client{}
	log.Println(config.ApiUrl + endpoint)
	req, _ := http.NewRequest("GET", config.ApiUrl+endpoint, nil)
	req.WithContext(ctx)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+config.ApiKey)

	res, err := client.Do(req)

	if err != nil {
		return nil, &ErrorTechnical{"Failed to fetch URL: " + endpoint, err}
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter := res.Header.Get("Retry-After")
		if retryAfter != "" {
			// TODO:
			return nil, &ErrorHTTP{
				res.StatusCode,
				fmt.Sprintf("Rate limit exceeded, retry after %s", retryAfter),
			}
		} else {
			return nil, &ErrorHTTP{res.StatusCode, "Too many requests"}
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, &ErrorHTTP{res.StatusCode, res.Status}
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, &ErrorTechnical{"Failed to read body data: " + endpoint, err}
	}

	return body, nil
}
