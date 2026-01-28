package topic_test

import (
	"testing"

	"github.com/gophercast/gophercast/internal/domain/topic"
)

func TestNewTopic(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    topic.Topic
		wantErr bool
	}{
		{
			name:    "valid simple name",
			input:   "users",
			want:    topic.Topic{},
			wantErr: false,
		},
		{
			name:    "valid dotted name",
			input:   "user.created",
			want:    topic.Topic{},
			wantErr: false,
		},
		{
			name:    "valid hyphenated",
			input:   "system-status",
			want:    topic.Topic{},
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			want:    topic.Topic{},
			wantErr: true,
		},
		{
			name:    "name with space",
			input:   "user created",
			want:    topic.Topic{},
			wantErr: true,
		},
		{
			name:    "name with slash",
			input:   "user/created",
			want:    topic.Topic{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topic.New(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.String() != tt.input {
				t.Errorf("New() = %v, want %v", got.String(), tt.input)
			}
		})
	}
}

func TestTopicEquals(t *testing.T) {
	tests := []struct {
		name     string
		topic1   string
		topic2   string
		expected bool
	}{
		{
			name:     "same name topics are equal",
			topic1:   "users",
			topic2:   "users",
			expected: true,
		},
		{
			name:     "different name topics are not equal",
			topic1:   "users",
			topic2:   "orders",
			expected: false,
		},
		{
			name:     "comparison is case-sensitive",
			topic1:   "Users",
			topic2:   "users",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1, _ := topic.New(tt.topic1)
			t2, _ := topic.New(tt.topic2)

			if t1.Equals(t2) != tt.expected {
				t.Errorf("Equals() = %v, want %v", t1.Equals(t2), tt.expected)
			}
		})
	}
}
