package workflows

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/subspace/linkedin-automation/internal/stealth"
)

// ConnectionConfig holds configuration for connection requests
type ConnectionConfig struct {
	Note string
}

// SendConnectionRequest sends a connection request to a LinkedIn profile
func SendConnectionRequest(page *rod.Page, profileURL string, config ConnectionConfig) error {
	// Navigate to profile
	err := page.Navigate(profileURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}
	
	stealth.RandomSleep(2000, 4000)
	
	// Wait for page to load
	// NOTE: User may need to update this selector
	profileSelector := ".pv-text-details__left-panel, .ph5, [data-test-id='profile-container']"
	_, err = page.Timeout(10 * time.Second).Element(profileSelector)
	if err != nil {
		return fmt.Errorf("profile page not loaded (selector may have changed): %w", err)
	}
	
	// Look for Connect button
	// NOTE: User MUST update these selectors - LinkedIn changes them frequently
	connectSelectors := []string{
		"button[aria-label*='Connect']",
		"button:has-text('Connect')",
		".pvs-profile-actions button",
		"button[data-control-name='connect']",
		"//button[contains(text(), 'Connect')]",
	}
	
	var connectButton *rod.Element
	for _, selector := range connectSelectors {
		var err error
		if selector[0:2] == "//" {
			// XPath selector
			connectButton, err = page.Timeout(3 * time.Second).ElementX(selector)
		} else {
			connectButton, err = page.Timeout(3 * time.Second).Element(selector)
		}
		
		if err == nil {
			break
		}
	}
	
	if connectButton == nil {
		return fmt.Errorf("Connect button not found - profile may already be connected or selectors need updating")
	}
	
		// Click Connect button using human-like movement
		box, err := connectButton.Shape()
		if err == nil && len(box.Quads) > 0 && len(box.Quads[0]) >= 8 {
			// Extract coordinates from quads
			x1, y1 := box.Quads[0][0], box.Quads[0][1]
			x2, y2 := box.Quads[0][2], box.Quads[0][3]
			x3, y3 := box.Quads[0][4], box.Quads[0][5]
			x4, y4 := box.Quads[0][6], box.Quads[0][7]
			
			// Calculate bounding box center
			minX := x1
			if x2 < minX { minX = x2 }
			if x3 < minX { minX = x3 }
			if x4 < minX { minX = x4 }
			
			maxX := x1
			if x2 > maxX { maxX = x2 }
			if x3 > maxX { maxX = x3 }
			if x4 > maxX { maxX = x4 }
			
			minY := y1
			if y2 < minY { minY = y2 }
			if y3 < minY { minY = y3 }
			if y4 < minY { minY = y4 }
			
			maxY := y1
			if y2 > maxY { maxY = y2 }
			if y3 > maxY { maxY = y3 }
			if y4 > maxY { maxY = y4 }
			
			centerX := (minX + maxX) / 2
			centerY := (minY + maxY) / 2
			
			err = stealth.HumanMouseMove(page, stealth.Point{X: 0, Y: 0}, stealth.Point{X: centerX, Y: centerY})
			if err == nil {
				stealth.RandomSleep(50, 150)
			}
		}
	
	// Perform the click
		// Perform the click
	err = connectButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click Connect button: %w", err)
	}
	
	stealth.RandomSleep(1000, 2000)
	
	// Check if a modal/dialog appeared for adding a note
	// NOTE: User may need to update these selectors
	noteSelectors := []string{
		"textarea[name='message']",
		"textarea[placeholder*='message']",
		"textarea[placeholder*='note']",
		".send-invite__custom-message textarea",
		"#custom-message",
	}
	
	var noteField *rod.Element
	for _, selector := range noteSelectors {
		var err error
		noteField, err = page.Timeout(3 * time.Second).Element(selector)
		if err == nil {
			break
		}
	}
	
	if noteField != nil && config.Note != "" {
		// Click the note field first
		err = noteField.Click(proto.InputMouseButtonLeft, 1)
		if err != nil {
			return fmt.Errorf("failed to click note field: %w", err)
		}
		
		stealth.RandomSleep(200, 500)
		
		// Clear existing text
		noteField.Input("")
		
		// Type the note character by character with human-like delays
		for _, char := range config.Note {
			err = noteField.Input(string(char))
			if err != nil {
				return fmt.Errorf("failed to type note: %w", err)
			}
			time.Sleep(stealth.HumanTypingDelay())
		}
		
		stealth.RandomSleep(500, 1000)
	}
	
	// Find and click Send button
	// NOTE: User may need to update these selectors
	sendSelectors := []string{
		"button[aria-label*='Send']",
		"button:has-text('Send')",
		"button[data-control-name='send_invite']",
		".send-invite__actions button[type='submit']",
		"//button[contains(text(), 'Send')]",
	}
	
	var sendButton *rod.Element
	for _, selector := range sendSelectors {
		var err error
		if selector[0:2] == "//" {
			sendButton, err = page.Timeout(3 * time.Second).ElementX(selector)
		} else {
			sendButton, err = page.Timeout(3 * time.Second).Element(selector)
		}
		
		if err == nil {
			break
		}
	}
	
	if sendButton != nil {
		err = sendButton.Click(proto.InputMouseButtonLeft, 1)
		if err != nil {
			return fmt.Errorf("failed to click Send button: %w", err)
		}
		
		stealth.RandomSleep(1000, 2000)
	}
	
	return nil
}

