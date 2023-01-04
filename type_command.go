package telescope

import "github.com/sirupsen/logrus"

// Command @Bean
type Command struct {
	Command   string `json:"command,omitempty"`
	ExitCode  string `json:"exit_code,omitempty"`
	Arguments string `json:"arguments,omitempty"`
	Options   string `json:"options,omitempty"`
}

func (b Command) BindType() string {
	return "command"
}

func (b Command) Handler(entry *logrus.Entry) (*entries, []tag) {
	return &entries{}, nil
}
