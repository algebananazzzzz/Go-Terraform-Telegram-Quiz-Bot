package main

type Quiz struct {
	Id        int        `json:"id" yaml:"id"`
	Name      string     `json:"quiz_name" yaml:"quiz_name"`
	Length    int        `json:"quiz_length" yaml:"quiz_length"`
	Questions []Question `json:"questions" yaml:"questions"`
}

type Question struct {
	Id      int      `json:"id" yaml:"id"`
	Title   string   `json:"title" yaml:"title"`
	Options []string `json:"options" yaml:"options"`
	Answer  int8     `json:"answer" yaml:"answer"`
}

const (
	DefaultState = 0
	QuizState    = 1
)

const (
	unknownCommandMessage = "Unknown command. Available commands: /start to start bot, /exit to exit quiz."
)
