package auth

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/subspace/linkedin-automation/internal/stealth"
)

const (
	linkedInLoginURL = "https://www.linkedin.com/login"
	maxLoginAttempts = 3
)

// LoginCredentials holds LinkedIn login information
type LoginCredentials struct {
	Email    string
	Password string
}

// Login attempts to log in to LinkedIn
// It first tries to load saved cookies, and if that fails, performs a fresh login
func Login(page *rod.Page, credentials LoginCredentials) error {
	// Try to load saved cookies first
	cookiesLoaded, err := LoadCookies(page)
	if err != nil {
		return fmt.Errorf("failed to load cookies: %w", err)
	}
	
	if cookiesLoaded {
		// Navigate to LinkedIn to verify session
		err = page.Navigate(linkedInLoginURL)
		if err != nil {
			return fmt.Errorf("failed to navigate: %w", err)
		}
		
		stealth.RandomSleep(2000, 3000)
		
		// Check if we're already logged in by looking for profile indicator
		currentURL := page.MustInfo().URL
		if currentURL != linkedInLoginURL && currentURL != linkedInLoginURL+"/" {
			// We're logged in!
			return nil
		}
		
		// Check for login form (if it exists, cookies expired)
		_, err = page.Timeout(3 * time.Second).Element("input[name='session_key']")
		if err == nil {
			// Login form found, cookies expired, proceed with fresh login
		} else {
			// No login form, assume we're logged in
			return nil
		}
	}
	
	// Perform fresh login
	return performLogin(page, credentials)
}

// performLogin performs the actual login process
func performLogin(page *rod.Page, credentials LoginCredentials) error {
	// Navigate to login page
	err := page.Navigate(linkedInLoginURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}
	
	stealth.RandomSleep(2000, 4000)
	
	// Wait for login form to be ready
	// NOTE: User may need to update these selectors if LinkedIn changes them
	emailSelector := "input[name='session_key']"
	passwordSelector := "input[name='session_password']"
	submitSelector := "button[type='submit']"
	
		// Wait for email field
	_, err = page.Timeout(10 * time.Second).Element(emailSelector)
	if err != nil {
		return fmt.Errorf("email field not found (selector may have changed): %w", err)
	}
	
	// Click email field and type
	err = stealth.HumanClick(page, emailSelector)
	if err != nil {
		return fmt.Errorf("failed to click email field: %w", err)
	}
	
	stealth.RandomSleep(200, 500)
	
	err = stealth.HumanType(page, emailSelector, credentials.Email)
	if err != nil {
		return fmt.Errorf("failed to type email: %w", err)
	}
	
	stealth.RandomSleep(500, 1000)
	
	// Click password field and type
	err = stealth.HumanClick(page, passwordSelector)
	if err != nil {
		return fmt.Errorf("failed to click password field: %w", err)
	}
	
	stealth.RandomSleep(200, 500)
	
	err = stealth.HumanType(page, passwordSelector, credentials.Password)
	if err != nil {
		return fmt.Errorf("failed to type password: %w", err)
	}
	
	stealth.RandomSleep(500, 1000)
	
	// Click submit button
	err = stealth.HumanClick(page, submitSelector)
	if err != nil {
		return fmt.Errorf("failed to click submit button: %w", err)
	}
	
	// Wait for navigation (either to feed or back to login if failed)
	stealth.RandomSleep(3000, 5000)
	
	// Check if login was successful
	currentURL := page.MustInfo().URL
	if currentURL == linkedInLoginURL || currentURL == linkedInLoginURL+"/" {
		// Still on login page, check for error message
		// NOTE: User may need to update error selector
		errorElement, err := page.Timeout(2 * time.Second).Element(".alert-error, .error-for-password, [role='alert']")
		if err == nil {
			errorText, _ := errorElement.Text()
			return fmt.Errorf("login failed: %s", errorText)
		}
		return fmt.Errorf("login failed: still on login page")
	}
	
	// Login successful, save cookies
	err = SaveCookies(page)
	if err != nil {
		// Log warning but don't fail - we're still logged in
		fmt.Printf("Warning: Failed to save cookies: %v\n", err)
	}
	
	return nil
}

