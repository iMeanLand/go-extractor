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
	"net/url"
	"path"
	"strconv"
	"time"
)

var (
	retries = 5
)

const (
	limit = 100
)

type Response struct {
	Data     json.RawMessage `json:"data"`
	NextPage struct {
		Offset string `json:"offset"`
		Path   string `json:"path"`
		Uri    string `json:"uri"`
	} `json:"next_page"`
}

type Request struct {
	endpoint string
	method   string
	params   url.Values
	page     string
	offset   string
}

func Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	extract(ctx)

	for {
		select {
		case <-ticker.C:
			// TODO: handle if did not finish previous yet
			log.Println("extracting..")
			extract(ctx)
		case <-ctx.Done():
			ticker.Stop()
			log.Println("Stopping extractor..")
			return
		}
	}
}

func extract(ctx context.Context) {
	data, err := fetch(ctx, &Request{
		endpoint: "/workspaces",
		method:   "GET",
	})
	if err != nil {
		log.Println("Error fetching workspaces:", err)
		return
	}

	var workspaces []Workspace
	err = json.Unmarshal(*data, &workspaces)
	if err != nil {
		log.Println("Failed to decode workspaces:", err)
		return
	}

	if len(workspaces) == 0 {
		log.Println("No workspaces found")
	}

	for _, workspace := range workspaces {
		params := url.Values{}
		params.Add("workspace", workspace.Gid)
		data, err = fetch(ctx, &Request{
			endpoint: "/users",
			method:   "GET",
			params:   params,
		})

		if err != nil {
			log.Println("Error fetching users:", err)
			continue
		}

		var users []User
		err = json.Unmarshal(*data, &users)
		if err != nil {
			log.Println("Failed to decode users:", err)
			continue
		}

		for _, user := range users {
			log.Print("saving file user")

			err = helper.SaveToJsonFile(fmt.Sprintf("user-%s", user.Gid), user)
		}

		log.Print("saving file workspace")

		err = helper.SaveToJsonFile(fmt.Sprintf("workspace-%s", workspace.Gid), workspace)

	}
}

func fetch(ctx context.Context, request *Request) (*json.RawMessage, error) {
	var err error

	var allData json.RawMessage
	for i := 0; i < retries; i++ {

		for {
			response, err := fetchUrl(ctx, request)

			if err != nil {
				switch e := err.(type) {
				case *ErrorTechnical:
					log.Printf("[HTTP ERROR] %v", e)
				case *ErrorHTTP:
					log.Printf("[HTTP ERROR] %v", e)
				default:
					log.Printf("[UNKNOWN] %v", err)
				}
				break
			}

			log.Printf("log, d %v", response.Data)

			allData = append(allData, response.Data...)

			if response.Data == nil || response.NextPage.Path == "" {
				break
			}

			request.page = response.NextPage.Path
			request.offset = response.NextPage.Offset
		}

		if err == nil {
			break
		}

		log.Printf("Retrying fetch (%d/%d):", i+1, retries)
		time.Sleep(1 * time.Second) // retry interval
	}

	log.Println("Successfully fetched")

	return &allData, err
}

func fetchUrl(ctx context.Context, request *Request) (*Response, error) {
	client := &http.Client{}

	if request.page != "" {
		request.params.Add("page", request.page)
		request.params.Add("offset", request.offset)
		request.params.Add("limit", strconv.Itoa(limit))
	}

	baseUrl, _ := url.Parse(config.ApiUrl)
	baseUrl.Path = path.Join(baseUrl.Path, request.endpoint)

	baseUrl.RawQuery = request.params.Encode()

	req, _ := http.NewRequest(request.method, baseUrl.String(), nil)
	req.WithContext(ctx)

	log.Println(baseUrl.String())

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+config.ApiKey)

	res, err := client.Do(req)

	log.Printf("Requesting .. %s", baseUrl.String())

	if err != nil {
		return nil, &ErrorTechnical{"Failed to fetch URL: " + baseUrl.String(), err}
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter := res.Header.Get("Retry-After")
		if retryAfter != "" {
			return nil, &ErrorHTTP{
				res.StatusCode,
				fmt.Sprintf("Rate limit exceeded, retry after %s", retryAfter),
			}
		} else {
			return nil, &ErrorHTTP{res.StatusCode, "Too many requests"}
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, &ErrorHTTP{res.StatusCode, "Unknown HTTP error"}
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, &ErrorTechnical{"Failed to read body data: " + baseUrl.String(), err}
	}

	var data *Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, &ErrorTechnical{"Failed to unmarshal data: " + baseUrl.String(), err}
	}

	return data, nil
}
