package welcome

import (
	"testing"
	"time"
)

func TestWelcomeWait(t *testing.T) {
	now := time.Now()
	if got := welcomeWait(time.Time{}, now, welcomeInterval); got != 0 {
		t.Fatalf("first welcome should not wait, got %v", got)
	}
	if got := welcomeWait(now.Add(-time.Second), now, welcomeInterval); got < 2499*time.Millisecond || got > 2501*time.Millisecond {
		t.Fatalf("expected about 2.5s wait, got %v", got)
	}
	if got := welcomeWait(now.Add(-4*time.Second), now, welcomeInterval); got != 0 {
		t.Fatalf("welcome after the interval should not wait, got %v", got)
	}
}
