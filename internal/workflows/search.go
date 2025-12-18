package workflows

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/subspace/linkedin-automation/internal/stealth"
)

// SearchConfig holds configuration for search operations
type SearchConfig struct {
	SearchURL    string
	MaxProfiles  int
	ScrollPauses int
}

// ProfileResult represents a scraped LinkedIn profile
type ProfileResult struct {
	URL      string
	Name     string
	Headline string
}

// SearchAndScrape searches LinkedIn and extracts profile URLs
func SearchAndScrape(page *rod.Page, config SearchConfig) ([]ProfileResult, error) {
	// Navigate to search URL
	err := page.Navigate(config.SearchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to search URL: %w", err)
	}
	
	stealth.RandomSleep(3000, 5000)
	
	// Wait for search results to load
	// NOTE: User may need to update this selector
	resultsSelector := ".search-results-container, .reusable-search__result-container, [data-test-id='search-result']"
	_, err = page.Timeout(10 * time.Second).Element(resultsSelector)
	if err != nil {
		return nil, fmt.Errorf("search results not found (selector may have changed): %w", err)
	}
	
	var profiles []ProfileResult
	
	// Scroll and collect profiles
	for i := 0; i < config.ScrollPauses; i++ {
		// Scroll down
		scrollAmount := 300 + rand.Intn(400)
		err = stealth.HumanScroll(page, scrollAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to scroll: %w", err)
		}
		
		// Pause between scrolls
		stealth.RandomSleep(1000, 2500)
		
		// Extract profile links from current view
		// NOTE: User may need to update these selectors
		profileLinks, err := page.Elements("a[href*='/in/']:not([href*='#']):not([href*='?'])")
		if err == nil {
			for _, link := range profileLinks {
				href, err := link.Attribute("href")
				if err != nil || href == nil {
					continue
				}
				
				// Clean URL
				fullURL := *href
				if fullURL[0] == '/' {
					fullURL = "https://www.linkedin.com" + fullURL
				}
				
				// Check if we already have this profile
				exists := false
				for _, p := range profiles {
					if p.URL == fullURL {
						exists = true
						break
					}
				}
				
				if !exists {
					// Try to get name and headline
					name := ""
					headline := ""
					
					// Try to find name in parent container
					parent, err := link.Parent()
					if err == nil {
						nameElem, err := parent.Timeout(1 * time.Second).Element(".entity-result__title-text a, .search-result__title a")
						if err == nil {
							name, _ = nameElem.Text()
						}
						
						headlineElem, err := parent.Timeout(1 * time.Second).Element(".entity-result__primary-subtitle, .search-result__snippets")
						if err == nil {
							headline, _ = headlineElem.Text()
						}
					}
					
					profiles = append(profiles, ProfileResult{
						URL:      fullURL,
						Name:     name,
						Headline: headline,
					})
					
					// Stop if we've reached max profiles
					if len(profiles) >= config.MaxProfiles {
						break
					}
				}
			}
		}
		
		if len(profiles) >= config.MaxProfiles {
			break
		}
	}
	
	return profiles, nil
}

