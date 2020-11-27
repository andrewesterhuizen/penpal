package lexer

const (
	TokenInstruction   = "INSTRUCTION"
	TokenLabel         = "LABEL"
	TokenDefine        = "DEFINE"
	TokenFileInclude   = "FILE_INCLUDE"
	TokenSystemInclude = "SYSTEM_INCLUDE"
)

type Token struct {
	Type       string
	Value      string
	LineNumber int
	FileName   string
	Args       []Arg
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

func newFileIncludeToken(value string) Token {
	return Token{
		Type:  TokenFileInclude,
		Value: value,
	}
}

func newSystemIncludeToken(value string) Token {
	return Token{
		Type:  TokenSystemInclude,
		Value: value,
	}
}
