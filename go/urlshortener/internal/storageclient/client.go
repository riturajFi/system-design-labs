package storageclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"urlshortener/internal/domain"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Mapping struct {
	ID       domain.ID       `json:"id"`
	ShortURL domain.ShortURL `json:"short_url"`
	LongURL  domain.LongURL  `json:"long_url"`
}

type idResponse struct {
	ID domain.ID `json:"id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetByShortURL(ctx context.Context, shortURL domain.ShortURL) (*Mapping, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/mappings/short/"+url.PathEscape(string(shortURL)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("build get by short request: %w", err)
	}

	var mapping Mapping
	found, err := c.doJSON(req, &mapping)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return &mapping, nil
}

func (c *Client) GetByLongURL(ctx context.Context, longURL domain.LongURL) (*Mapping, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/mappings/long?url="+url.QueryEscape(string(longURL)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("build get by long request: %w", err)
	}

	var mapping Mapping
	found, err := c.doJSON(req, &mapping)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return &mapping, nil
}

func (c *Client) Save(ctx context.Context, mapping Mapping) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("marshal mapping: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/mappings",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("build save request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("save mapping request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	return c.decodeError(resp)
}

func (c *Client) NextID(ctx context.Context) (domain.ID, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/ids/next", nil)
	if err != nil {
		return 0, fmt.Errorf("build next id request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("next id request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, c.decodeError(resp)
	}

	var payload idResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("decode next id response: %w", err)
	}

	return payload.ID, nil
}

func (c *Client) doJSON(req *http.Request, out any) (bool, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("%s %s: %w", req.Method, req.URL.String(), err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return false, fmt.Errorf("decode response body: %w", err)
		}
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, c.decodeError(resp)
	}
}

func (c *Client) decodeError(resp *http.Response) error {
	var payload errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err == nil && payload.Error != "" {
		return fmt.Errorf("storage service returned %d: %s", resp.StatusCode, payload.Error)
	}

	return fmt.Errorf("storage service returned %d", resp.StatusCode)
}
