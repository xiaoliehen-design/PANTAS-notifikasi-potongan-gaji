package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/config"
)

type Client struct {
	baseURL string
	key     string
	bucket  string
	http    *http.Client
}

func New(cfg config.Config) *Client {
	return &Client{
		baseURL: cfg.SupabaseURL,
		key:     cfg.SupabaseServiceKey,
		bucket:  cfg.SupabaseStorageBucket,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Upload(ctx context.Context, path, contentType string, data []byte) error {
	endpoint := c.objectURL(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}
	c.authorize(req)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "false")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("upload storage: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return fmt.Errorf("upload storage status %d: %s", resp.StatusCode, strings.TrimSpace(string(message)))
	}
	return nil
}

func (c *Client) Download(ctx context.Context, path string) (io.ReadCloser, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.objectURL(path), nil)
	if err != nil {
		return nil, "", err
	}
	c.authorize(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download storage: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return nil, "", fmt.Errorf("download storage status %d: %s", resp.StatusCode, strings.TrimSpace(string(message)))
	}
	return resp.Body, resp.Header.Get("Content-Type"), nil
}

func (c *Client) Delete(ctx context.Context, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	for _, path := range paths {
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.objectURL(path), nil)
		if err != nil {
			return err
		}
		c.authorize(req)
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("delete storage status %d", resp.StatusCode)
		}
	}
	return nil
}

func (c *Client) objectURL(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for index := range parts {
		parts[index] = url.PathEscape(parts[index])
	}
	return c.baseURL + "/storage/v1/object/" + url.PathEscape(c.bucket) + "/" + strings.Join(parts, "/")
}

func (c *Client) authorize(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("apikey", c.key)
}
