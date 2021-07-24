package token

type Token interface {
	GetValue() string
	GetOwner() string
}

type token struct {
	Value string
	Owner string
}

func NewToken(value, owner string) Token {
	return &token{Value: value, Owner: owner}
}

func (t *token) GetValue() string {
	return t.Value
}

func (t *token) GetOwner() string {
	return t.Owner
}
