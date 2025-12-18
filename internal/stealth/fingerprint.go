package stealth

import (
	"math/rand"
	"time"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// StealthBrowser creates a browser instance with anti-detection measures
// In internal/stealth/fingerprint.go, modify StealthBrowser function:

func StealthBrowser(headless bool) (*rod.Browser, error) {
	// Disable leakless by setting environment variable
	os.Setenv("ROD_LEAKLESS", "0")
	
	// Try to use existing Chrome installation first
	l := launcher.New().
		Headless(headless).
		// Disable automation flags
		Set("disable-blink-features", "AutomationControlled").
		Set("exclude-switches", "enable-automation").
		Set("useAutomationExtension", "false").
		// Additional stealth flags
		Set("disable-dev-shm-usage").
		Set("no-sandbox").
		Set("disable-setuid-sandbox").
		Set("disable-web-security").
		Set("disable-features", "IsolateOrigins,site-per-process")
	
	// Try to find Chrome in common locations
	chromePaths := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		os.Getenv("LOCALAPPDATA") + `\Google\Chrome\Application\chrome.exe`,
	}
	
	for _, path := range chromePaths {
		if _, err := os.Stat(path); err == nil {
			l = l.Bin(path)
			break
		}
	}
	
	browserURL, err := l.Launch()
	if err != nil {
		return nil, err
	}
	
	browser := rod.New().ControlURL(browserURL)
	if err := browser.Connect(); err != nil {
		return nil, err
	}
	
	return browser, nil
}

// ApplyStealthFingerprint applies anti-detection JavaScript to mask automation
func ApplyStealthFingerprint(page *rod.Page) error {
	// Mask webdriver property - wrap in arrow function for go-rod
	stealthScript := `() => {
		// Override navigator.webdriver
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});
		
		// Override chrome object
		window.chrome = {
			runtime: {}
		};
		
		// Override permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);
		
		// Override plugins
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});
		
		// Override languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});
		
		// Override platform
		Object.defineProperty(navigator, 'platform', {
			get: () => 'Win32'
		});
		
		// Remove automation indicators
		delete window.cdc_adoQpoasnfa76pfcZLmcfl_Array;
		delete window.cdc_adoQpoasnfa76pfcZLmcfl_Promise;
		delete window.cdc_adoQpoasnfa76pfcZLmcfl_Symbol;
	}`
	
	_, err := page.Eval(stealthScript)
	return err
}

// CreateStealthPage creates a new page with stealth features applied
func CreateStealthPage(browser *rod.Browser, userAgent string) (*rod.Page, error) {
	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}
	
	// Set user agent
	if userAgent == "" {
		userAgent = GetRandomUserAgent()
	}
	
	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      userAgent,
		AcceptLanguage: "en-US,en;q=0.9",
	})
	if err != nil {
		return nil, err
	}
	
	// Apply stealth fingerprinting
	err = ApplyStealthFingerprint(page)
	if err != nil {
		return nil, err
	}
	
	// Set viewport to common desktop size
	err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  1920,
		Height: 1080,
		DeviceScaleFactor: 1,
		Mobile: false,
	})
	
	return page, err
}

// HumanMouseMove moves the mouse along a Bézier curve path
func HumanMouseMove(page *rod.Page, start, end Point) error {
	path := GenerateOvershootPath(start, end)
	
	for _, point := range path {
		err := page.Mouse.MoveTo(proto.Point{
			X: point.X,
			Y: point.Y,
		})
		if err != nil {
			return err
		}
		// Small delay between movements
		RandomSleep(5, 15)
	}
	
	return nil
}

// HumanClick performs a human-like click with mouse movement and delays
func HumanClick(page *rod.Page, selector string) error {
	element, err := page.Element(selector)
	if err != nil {
		return err
	}
	
	// Get element position - Shape() returns quads, we need the first quad
	box, err := element.Shape()
	if err != nil {
		return err
	}
	
	// Extract coordinates from the quads (first quad, first point)
	if len(box.Quads) == 0 || len(box.Quads[0]) < 8 {
		return err
	}
	
	// Quads format: [x1, y1, x2, y2, x3, y3, x4, y4]
	// Use the first point (x1, y1) and calculate center from bounding box
	x1, y1 := box.Quads[0][0], box.Quads[0][1]
	x2, y2 := box.Quads[0][2], box.Quads[0][3]
	x3, y3 := box.Quads[0][4], box.Quads[0][5]
	x4, y4 := box.Quads[0][6], box.Quads[0][7]
	
	// Calculate bounding box
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
	
	// Calculate center
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	
	// Get current mouse position (approximate)
	currentPos := Point{X: 0, Y: 0} // Default to top-left
	
	// Move mouse to element using Bézier curve
	err = HumanMouseMove(page, currentPos, Point{X: centerX, Y: centerY})
	if err != nil {
		return err
	}
	
	// Small pause before click
	RandomSleep(50, 150)
	
	// Click with button and count
	err = element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Small pause after click
	RandomSleep(50, 150)
	
	return nil
}

// HumanScroll scrolls the page with human-like patterns
func HumanScroll(page *rod.Page, scrollAmount int) error {
	// Break scroll into smaller chunks with pauses
	chunkSize := 200 + rand.Intn(300)
	remaining := scrollAmount
	
	for remaining > 0 {
		currentScroll := chunkSize
		if currentScroll > remaining {
			currentScroll = remaining
		}
		
		// Scroll down using JavaScript (more reliable than Mouse API)
		_, err := page.Eval(`() => window.scrollBy(0, arguments[0])`, currentScroll)
		if err != nil {
			return err
		}
		
		// Random pause
		RandomSleep(100, 500)
		
		// Occasionally scroll up slightly (human behavior)
		if rand.Float64() < 0.2 {
			backScroll := 20 + rand.Intn(50)
			page.Eval(`() => window.scrollBy(0, arguments[0])`, -backScroll)
			RandomSleep(50, 150)
		}
		
		remaining -= currentScroll
	}
	
	return nil
}

// HumanType types text with human-like delays between keystrokes
func HumanType(page *rod.Page, selector, text string) error {
	element, err := page.Element(selector)
	if err != nil {
		return err
	}
	
	err = element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Clear existing text
	err = element.SelectAllText()
	if err != nil {
		// If SelectAllText fails, try manual clear
		element.Input("")
	}
	
	// Type character by character with delays
	for _, char := range text {
		err = element.Input(string(char))
		if err != nil {
			return err
		}
		time.Sleep(HumanTypingDelay())
	}
	
	return nil
}

