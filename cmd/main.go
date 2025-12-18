package main

import (
	"fmt"
	"log"
	"os"

	"github.com/subspace/linkedin-automation/config"
	"github.com/subspace/linkedin-automation/internal/auth"
	"github.com/subspace/linkedin-automation/internal/stealth"
	"github.com/subspace/linkedin-automation/internal/storage"
	"github.com/subspace/linkedin-automation/internal/workflows"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Initialize history database
	historyDB, err := storage.NewHistoryDB()
	if err != nil {
		log.Fatalf("Failed to initialize history database: %v", err)
	}
	defer historyDB.Close()
	
	// Create stealth browser
	browser, err := stealth.StealthBrowser(cfg.Headless)
	if err != nil {
		log.Fatalf("Failed to create browser: %v", err)
	}
	defer browser.Close()
	
	// Create stealth page
	page, err := stealth.CreateStealthPage(browser, "")
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()
	
	// Authenticate
	credentials := auth.LoginCredentials{
		Email:    cfg.LinkedInEmail,
		Password: cfg.LinkedInPassword,
	}
	
	fmt.Println("Authenticating with LinkedIn...")
	err = auth.Login(page, credentials)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	fmt.Println("✓ Successfully authenticated")
	
	// Check if we have a search URL
	if cfg.SearchURL == "" {
		fmt.Println("No SEARCH_URL provided. Exiting.")
		fmt.Println("To use search functionality, set SEARCH_URL in your .env file")
		os.Exit(0)
	}
	
	// Perform search and scrape
	fmt.Println("Searching and scraping profiles...")
	searchConfig := workflows.SearchConfig{
		SearchURL:    cfg.SearchURL,
		MaxProfiles:  20, // Scrape more than we'll process
		ScrollPauses: 5,
	}
	
	profiles, err := workflows.SearchAndScrape(page, searchConfig)
	if err != nil {
		log.Fatalf("Failed to search and scrape: %v", err)
	}
	
	fmt.Printf("✓ Found %d profiles\n", len(profiles))
	
	// Process profiles (send connection requests)
	connectionConfig := workflows.ConnectionConfig{
		Note: cfg.ConnectionNote,
	}
	
	actionsCount := 0
	for _, profile := range profiles {
		// Check safety limit
		if actionsCount >= cfg.MaxActions {
			fmt.Printf("Reached safety limit of %d actions. Stopping.\n", cfg.MaxActions)
			break
		}
		
		// Check if already processed
		processed, err := historyDB.IsProcessed(profile.URL)
		if err != nil {
			log.Printf("Error checking if profile processed: %v", err)
			continue
		}
		
		if processed {
			fmt.Printf("Skipping %s (already processed)\n", profile.URL)
			continue
		}
		
		// Send connection request
		fmt.Printf("Sending connection request to %s...\n", profile.URL)
		err = workflows.SendConnectionRequest(page, profile.URL, connectionConfig)
		if err != nil {
			log.Printf("Failed to send connection request to %s: %v\n", profile.URL, err)
			continue
		}
		
		// Mark as processed
		err = historyDB.MarkProcessed(profile.URL, "connection_request")
		if err != nil {
			log.Printf("Failed to mark profile as processed: %v", err)
		}
		
		actionsCount++
		fmt.Printf("✓ Connection request sent (%d/%d)\n", actionsCount, cfg.MaxActions)
		
		// Random delay between actions
		stealth.RandomSleep(5000, 10000)
	}
	
	fmt.Printf("\n✓ Completed %d actions\n", actionsCount)
}

