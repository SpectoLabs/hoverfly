package faker

import (
	"fmt"
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// YouTube is a faker struct for YouTube
type YouTube struct {
	Faker *Faker
}

// GenerateVideoID returns a youtube video id
func (y YouTube) GenerateVideoID() (videoID string) {
	b := make([]byte, 11)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateFullURL returns a fake standard youtube video url
func (y YouTube) GenerateFullURL() string {
	return fmt.Sprintf("www.youtube.com/watch?v=%s", y.GenerateVideoID())
}

// GenerateShareURL returns a fake share youtube video url
func (y YouTube) GenerateShareURL() string {
	return fmt.Sprintf("youtu.be/%s", y.GenerateVideoID())
}

// GenerateEmbededURL returns a fake embedded youtube video url
func (y YouTube) GenerateEmbededURL() string {
	return fmt.Sprintf("www.youtube.com/embed/%s", y.GenerateVideoID())
}
