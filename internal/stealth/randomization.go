package stealth

import (
	"math/rand"
	"time"
)

// RandomSleep sleeps for a random duration between min and max milliseconds
// This mimics human behavior where actions are never perfectly timed
func RandomSleep(minMs, maxMs int) {
	if minMs < 0 {
		minMs = 0
	}
	if maxMs < minMs {
		maxMs = minMs
	}
	
	duration := time.Duration(minMs+rand.Intn(maxMs-minMs+1)) * time.Millisecond
	time.Sleep(duration)
}

// RandomSleepSeconds sleeps for a random duration between min and max seconds
func RandomSleepSeconds(minSec, maxSec float64) {
	if minSec < 0 {
		minSec = 0
	}
	if maxSec < minSec {
		maxSec = minSec
	}
	
	// Convert to milliseconds for more precision
	minMs := int(minSec * 1000)
	maxMs := int(maxSec * 1000)
	RandomSleep(minMs, maxMs)
}

// HumanTypingDelay returns a random delay between keystrokes (50-150ms)
// This simulates natural human typing speed
func HumanTypingDelay() time.Duration {
	return time.Duration(50+rand.Intn(100)) * time.Millisecond
}

// HumanScrollDelay returns a random delay for scroll actions (100-500ms)
func HumanScrollDelay() time.Duration {
	return time.Duration(100+rand.Intn(400)) * time.Millisecond
}

// HumanClickDelay returns a random delay before/after clicks (50-200ms)
func HumanClickDelay() time.Duration {
	return time.Duration(50+rand.Intn(150)) * time.Millisecond
}

// GetRandomUserAgent returns a realistic user agent string
func GetRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	}
	
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

