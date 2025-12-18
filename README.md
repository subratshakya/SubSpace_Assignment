# LinkedIn Automation Proof-of-Concept

A sophisticated LinkedIn automation tool built in Go using the `go-rod` library. This project demonstrates advanced stealth techniques for browser automation, including human-like mouse movements, fingerprint masking, and session persistence.

**⚠️ Educational Purpose Only**: This tool is designed for educational purposes to demonstrate anti-detection techniques. Use responsibly and in accordance with LinkedIn's Terms of Service.

## Features

### 🎯 Stealth Engine (Critical Component)
- **Bézier Curve Mouse Movements**: Human-like mouse paths with overshoot and correction
- **Randomization**: All delays and actions are randomized to mimic human behavior
- **Browser Fingerprinting**: Masks `navigator.webdriver` and applies realistic browser properties

### 🔐 Authentication
- **Session Persistence**: Saves and loads cookies to avoid repeated logins
- **Human-like Typing**: Simulates natural typing speeds with random delays

### 📊 Workflows
- **Search & Scrape**: Navigates LinkedIn search results and extracts profile URLs
- **Connection Requests**: Sends connection requests with custom notes
- **Safety Limits**: Hardcoded limits to prevent account bans during testing

### 💾 Data Storage
- **SQLite Database**: Tracks processed profiles to prevent duplicates
- **Cookie Management**: Persistent session storage

## Project Structure

```
linkedin-automation/
├── cmd/
│   └── main.go                 # Main entry point
├── internal/
│   ├── stealth/                # Stealth engine
│   │   ├── bezier.go          # Bézier curve mouse movements
│   │   ├── randomization.go   # Random delays and user agents
│   │   └── fingerprint.go    # Browser fingerprint masking
│   ├── auth/                   # Authentication module
│   │   ├── login.go           # Login logic
│   │   └── session.go         # Cookie management
│   ├── workflows/              # Workflow modules
│   │   ├── search.go          # Search and scraping
│   │   └── connection.go      # Connection requests
│   └── storage/                # Data storage
│       └── history.go         # SQLite database operations
├── config/
│   └── config.go              # Configuration management
├── data/                       # Runtime data (created automatically)
│   ├── cookies.json           # Saved session cookies
│   └── history.db             # SQLite database
├── .env                        # Environment variables (create from .env.example)
├── .env.example               # Example environment file
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Installation

### Prerequisites
- Go 1.21 or higher
- Chrome/Chromium browser (for go-rod)

### Setup

1. **Clone or download this repository**

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Create `.env` file** from `.env.example`:
   ```bash
   cp .env.example .env
   ```

4. **Configure your `.env` file** (see Required User Inputs section below)

5. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

## Configuration

All configuration is done through environment variables in the `.env` file. See `.env.example` for the required variables.

## Usage

1. **Set up your `.env` file** with LinkedIn credentials and search URL
2. **Run the application**: `go run cmd/main.go`
3. The tool will:
   - Authenticate with LinkedIn (or use saved cookies)
   - Search for profiles based on your SEARCH_URL
   - Send connection requests (up to MAX_ACTIONS limit)
   - Track processed profiles to avoid duplicates

## Technical Details

### Stealth Techniques

1. **Bézier Curve Mouse Movements**: Instead of straight-line movements, the mouse follows cubic Bézier curves with control points that create natural overshoot and correction patterns.

2. **Randomization**: Every action includes randomized delays:
   - Typing: 50-150ms between keystrokes
   - Scrolling: 100-500ms pauses
   - Clicks: 50-200ms before/after delays

3. **Browser Fingerprinting**:
   - Masks `navigator.webdriver` property
   - Sets realistic User-Agent strings
   - Overrides automation-related properties

### Safety Features

- **Hardcoded Limits**: Maximum actions per run (default: 5)
- **Duplicate Prevention**: SQLite database tracks processed profiles
- **Session Persistence**: Reduces login frequency

## Troubleshooting

### Common Issues

1. **"Selector not found" errors**: LinkedIn frequently changes their HTML structure. You'll need to update CSS selectors in the code (see Required User Inputs section).

2. **Login failures**: 
   - Verify your credentials in `.env`
   - Clear `data/cookies.json` and try again
   - Check if LinkedIn requires 2FA (not currently supported)

3. **Browser not launching**:
   - Ensure Chrome/Chromium is installed
   - Check that the browser path is correct

## Development

### Adding New Workflows

1. Create a new function in `internal/workflows/`
2. Use the stealth functions from `internal/stealth/` for all interactions
3. Always include randomized delays
4. Update the main.go to call your new workflow

### Modifying Stealth Behavior

The stealth engine is modular:
- `bezier.go`: Adjust curve parameters for different mouse movement styles
- `randomization.go`: Modify delay ranges
- `fingerprint.go`: Add additional fingerprint masking

## License

This project is for educational purposes only. Use at your own risk and in accordance with LinkedIn's Terms of Service.

## Disclaimer

This tool is provided for educational purposes to demonstrate browser automation and anti-detection techniques. The authors are not responsible for any misuse of this software. Using automation tools may violate LinkedIn's Terms of Service and could result in account suspension or termination.

---

## **REQUIRED USER INPUTS**

### 1. Environment Variables (`.env` file)

You **MUST** create a `.env` file in the project root with the following variables:

```env
# Required: Your LinkedIn credentials
LINKEDIN_EMAIL=your-email@example.com
LINKEDIN_PASSWORD=your-password

# Optional: Browser settings (default: false)
HEADLESS=false

# Optional: Safety limit (default: 5)
MAX_ACTIONS=5

# Optional: Connection request note
CONNECTION_NOTE=Hi! I'd like to connect with you.

# Required: LinkedIn search URL
# To get this:
# 1. Go to LinkedIn and perform a search (e.g., search for people with keyword "software engineer")
# 2. Copy the full URL from your browser
# 3. Paste it here
SEARCH_URL=https://www.linkedin.com/search/results/people/?keywords=software%20engineer
```

**Important Notes:**
- The `.env` file is gitignored and will not be committed to version control
- Never share your `.env` file or commit it to a repository
- If LinkedIn requires 2FA, you may need to manually complete the login process

### 2. CSS Selectors (CRITICAL - May Need Updates)

LinkedIn frequently changes their HTML structure and CSS selectors. You **MUST** inspect LinkedIn's current HTML and update the following selectors in the code if they stop working:

#### Login Selectors (`internal/auth/login.go`)

**Current selectors (lines ~30-35):**
- Email field: `input[name='session_key']`
- Password field: `input[name='session_password']`
- Submit button: `button[type='submit']`
- Error message: `.alert-error, .error-for-password, [role='alert']`

**How to update:**
1. Open LinkedIn login page in your browser
2. Right-click on the email field → Inspect Element
3. Look for the `name` or `id` attribute
4. Update the selector in `internal/auth/login.go`

#### Search Result Selectors (`internal/workflows/search.go`)

**Current selectors (lines ~30-50):**
- Results container: `.search-results-container, .reusable-search__result-container, [data-test-id='search-result']`
- Profile links: `a[href*='/in/']:not([href*='#']):not([href*='?'])`
- Name element: `.entity-result__title-text a, .search-result__title a`
- Headline element: `.entity-result__primary-subtitle, .search-result__snippets`

**How to update:**
1. Perform a LinkedIn search
2. Inspect the search results HTML
3. Find the container and profile link elements
4. Update selectors in `internal/workflows/search.go`

#### Connection Request Selectors (`internal/workflows/connection.go`)

**Current selectors (lines ~28-42, 77-83, 110-116):**

**Connect Button** (try these in order):
- `button[aria-label*='Connect']`
- `button:has-text('Connect')`
- `.pvs-profile-actions button`
- `button[data-control-name='connect']`
- XPath: `//button[contains(text(), 'Connect')]`

**Note/Message Field** (for adding a note):
- `textarea[name='message']`
- `textarea[placeholder*='message']`
- `textarea[placeholder*='note']`
- `.send-invite__custom-message textarea`
- `#custom-message`

**Send Button**:
- `button[aria-label*='Send']`
- `button:has-text('Send')`
- `button[data-control-name='send_invite']`
- `.send-invite__actions button[type='submit']`
- XPath: `//button[contains(text(), 'Send')]`

**How to update:**
1. Navigate to a LinkedIn profile
2. Inspect the "Connect" button
3. Check if it has `aria-label`, `data-control-name`, or specific classes
4. Update the selectors array in `internal/workflows/connection.go` (lines 36-42)
5. If a modal appears after clicking Connect, inspect the note field and send button
6. Update selectors in `internal/workflows/connection.go` (lines 77-83 and 110-116)

#### Profile Page Selectors (`internal/workflows/connection.go`)

**Current selector (line ~28):**
- Profile container: `.pv-text-details__left-panel, .ph5, [data-test-id='profile-container']`

**How to update:**
1. Navigate to any LinkedIn profile
2. Inspect the main profile content area
3. Update the selector in `internal/workflows/connection.go`

### 3. Testing and Verification Steps

Before running the full automation:

1. **Test Login:**
   - Run the code and verify it can log in
   - Check if cookies are saved to `data/cookies.json`
   - Verify subsequent runs use saved cookies

2. **Test Search:**
   - Set a `SEARCH_URL` in `.env`
   - Run and verify it can find profile links
   - Check the console output for found profiles

3. **Test Connection Request (Manual):**
   - Comment out the connection request loop in `cmd/main.go`
   - Manually test clicking Connect on one profile
   - Verify the selectors work with current LinkedIn HTML

4. **Monitor for Errors:**
   - Watch for "selector not found" errors
   - These indicate LinkedIn has changed their HTML
   - Update the corresponding selectors immediately

### 4. Selector Inspection Tools

**Recommended Browser DevTools:**
- Chrome/Edge: F12 → Elements tab
- Firefox: F12 → Inspector tab
- Right-click element → "Inspect Element"

**Finding the Right Selector:**
1. Use unique attributes: `id`, `name`, `data-*` attributes
2. Prefer semantic selectors: `aria-label`, `role`
3. Avoid fragile selectors: class names that look auto-generated
4. Test selectors in browser console: `document.querySelector('your-selector')`

### 5. Common Selector Patterns

If LinkedIn uses React or similar frameworks, look for:
- `data-test-id` attributes
- `aria-label` attributes
- `data-control-name` attributes
- Stable class names (not ones with random hashes)

**Example of a good selector:**
```css
button[aria-label="Connect with John Doe"]
```

**Example of a fragile selector (avoid if possible):**
```css
.css-1a2b3c4d-button  /* Auto-generated class */
```

---

**⚠️ IMPORTANT REMINDER:** LinkedIn's HTML structure changes frequently. If the automation stops working, the first thing to check is whether the selectors need updating. Always test with `HEADLESS=false` first to see what's happening in the browser.

