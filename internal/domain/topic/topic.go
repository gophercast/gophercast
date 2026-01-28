package topic

import (
	"errors"
	"regexp"
)

// Topic represents a named channel of communication in the pub/sub system.
// Topics are identified by strings and must follow naming rules.
type Topic struct {
	name string
}

// New creates a new Topic with the given name.
// The name must not be empty and must contain only letters, numbers, dots, and hyphens.
func New(name string) (Topic, error) {
	if name == "" {
		return Topic{}, errors.New("topic name cannot be empty")
	}

	if !isValidTopicName(name) {
		return Topic{}, errors.New("invalid topic name: must contain only letters, numbers, dots, and hyphens")
	}

	return Topic{name: name}, nil
}

// String returns the topic name.
func (t Topic) String() string {
	return t.name
}

// Equals returns true if two topics have the same name.
func (t Topic) Equals(other Topic) bool {
	return t.name == other.name
}

// isValidTopicName checks if the topic name follows the naming rules.
func isValidTopicName(name string) bool {
	// Only allow letters, numbers, dots, and hyphens
	pattern := `^[a-zA-Z0-9.-]+$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}
