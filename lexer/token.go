package lexer

const (
	TokenInstruction = "INSTRUCTION"
	TokenLabel       = "LABEL"
	TokenDefine      = "DEFINE"
)

type Token struct {
	Type  string
	Value string
	Args  []Arg
}

func newInstructionToken(value string, args []Arg) Token {
	return Token{
		Type:  TokenInstruction,
		Value: value,
		Args:  args,
	}
}

func newDefineToken(value string, args []Arg) Token {
	return Token{
		Type:  TokenDefine,
		Value: value,
		Args:  args,
	}
}

func newLabelToken(value string) Token {
	return Token{
		Type:  TokenLabel,
		Value: value,
	}
}
