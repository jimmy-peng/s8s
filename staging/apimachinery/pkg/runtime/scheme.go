package runtime

type Scheme struct {
	schemeName string
}

func NewScheme() *Scheme {
	s := &Scheme{
		schemeName: "hello",
	}
	return s
}