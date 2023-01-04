package telescope

import "github.com/sirupsen/logrus"

// Cache @Bean
type Cache struct {
	// hit missed set forget
	Type       string `json:"type"`
	Key        string `json:"key"`
	Value      string `json:"value,omitempty"`
	Expiration string `json:"expiration,omitempty"`
}

func (b Cache) BindType() string {
	return "cache"
}

func (b Cache) Handler(entry *logrus.Entry) (*entries, []tag) {
	return &entries{}, nil
}
