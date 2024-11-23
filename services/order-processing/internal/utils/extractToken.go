package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractToken retrieves the JWT token from the request context.
func ExtractToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return "", fmt.Errorf("invalid token format")
	}

	return tokenString, nil
}

func MakeRequestWithToken(ctx context.Context, method, url string, payload interface{}, token string) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("could not marshal payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %v", err)
	}

	return resp, nil
}
