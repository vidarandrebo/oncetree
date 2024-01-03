package greeter

type Greeting struct {
	Text   string
	Number int
}

func NewGreeting() *Greeting {
	return &Greeting{Text: "Hello there"}
}
