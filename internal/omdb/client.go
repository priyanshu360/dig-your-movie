package omdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "http://www.omdbapi.com/"
)

type Movie struct {
	Title       string `json:"Title"`
	Description string `json:"Plot"`
	Year        string `json:"Year"`
	Director    string `json:"Director"`
	Error       string `json:"Error"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetMovieDescription(ctx context.Context, movieName string) (Movie, error) {
	queryParams := url.Values{
		"apikey": {c.apiKey},
		"t":      {movieName},
	}
	u := baseURL + "?" + queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Movie{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Movie{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Movie{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Movie{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var movie Movie
	if err := json.Unmarshal(body, &movie); err != nil {
		return Movie{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if movie.Error != "" {
		return Movie{}, fmt.Errorf("omdb api error: %s", movie.Error)
	}

	return movie, nil
}
