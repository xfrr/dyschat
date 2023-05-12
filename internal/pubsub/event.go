package pubsub

import "regexp"

type Subject string

func (s Subject) Match(other Subject) bool {
	regex := regexp.MustCompile(string(s))
	return regex.MatchString(string(other))
}

type Event interface {
	Subject() Subject
	SubjectRegex() Subject
	UnmarshalJSON([]byte) error
}
