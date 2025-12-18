package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const cookiesFile = "data/cookies.json"

// CookieData represents saved cookie information
type CookieData struct {
	Cookies []*proto.NetworkCookie `json:"cookies"`
	SavedAt time.Time              `json:"saved_at"`
}

// SaveCookies saves browser cookies to a JSON file
func SaveCookies(page *rod.Page) error {
	cookies, err := page.Cookies([]string{})
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}
	
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	
	cookieData := CookieData{
		Cookies: cookies,
		SavedAt: time.Now(),
	}
	
	data, err := json.MarshalIndent(cookieData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}
	
	if err := os.WriteFile(cookiesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write cookies file: %w", err)
	}
	
	return nil
}

// LoadCookies loads cookies from a JSON file and applies them to the page
func LoadCookies(page *rod.Page) (bool, error) {
	// Check if cookies file exists
	if _, err := os.Stat(cookiesFile); os.IsNotExist(err) {
		return false, nil
	}
	
	data, err := os.ReadFile(cookiesFile)
	if err != nil {
		return false, fmt.Errorf("failed to read cookies file: %w", err)
	}
	
	var cookieData CookieData
	if err := json.Unmarshal(data, &cookieData); err != nil {
		return false, fmt.Errorf("failed to unmarshal cookies: %w", err)
	}
	
	// Check if cookies are too old (older than 7 days)
	if time.Since(cookieData.SavedAt) > 7*24*time.Hour {
		return false, nil
	}
	
		// Apply cookies to the page - convert NetworkCookie to NetworkCookieParam
		if len(cookieData.Cookies) > 0 {
			cookieParams := make([]*proto.NetworkCookieParam, 0, len(cookieData.Cookies))
			for _, cookie := range cookieData.Cookies {
				// Build URL from domain
				url := ""
				if cookie.Domain != "" {
					// Remove leading dot if present
					domain := cookie.Domain
					if len(domain) > 0 && domain[0] == '.' {
						domain = domain[1:]
					}
					if cookie.Secure {
						url = "https://" + domain
					} else {
						url = "http://" + domain
					}
				}
				
				cookieParams = append(cookieParams, &proto.NetworkCookieParam{
					Name:     cookie.Name,
					Value:    cookie.Value,
					URL:      url,
					Domain:   cookie.Domain,
					Path:     cookie.Path,
					Secure:   cookie.Secure,
					HTTPOnly: cookie.HTTPOnly,
					SameSite: cookie.SameSite,
					Expires:  cookie.Expires,
				})
			}
			
			err = page.SetCookies(cookieParams)
			if err != nil {
				return false, fmt.Errorf("failed to set cookies: %w", err)
			}
		}
		
		return true, nil
	}

// ClearCookies removes the saved cookies file
func ClearCookies() error {
	if _, err := os.Stat(cookiesFile); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(cookiesFile)
}