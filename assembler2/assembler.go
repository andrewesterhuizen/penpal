package assembler2

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

const EntryPointsTableSize = 12

type FileGetterFunc func(path string) (string, error)

func fileSystemFileGetterFunc(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

type Config struct {
	disableEntryPointsTable bool
	FileGetterFunc          FileGetterFunc
	SystemIncludes          map[string]string
	InteruptLabels          [3]string
}

type Assembler struct {
	config               Config
	getFile              FileGetterFunc
	systemIncludeSources map[string]string
	lexer                lexer_rewrite.Lexer
}

func New(config Config) Assembler {
	a := Assembler{config: config}

	if config.FileGetterFunc != nil {
		a.getFile = config.FileGetterFunc
	} else {
		a.getFile = fileSystemFileGetterFunc
	}

	if config.SystemIncludes != nil {
		a.systemIncludeSources = config.SystemIncludes
	} else {
		a.systemIncludeSources = make(map[string]string)
	}

	return a
}

func (a *Assembler) getIncludeTokens(filename string, tokens []lexer_rewrite.Token) ([]lexer_rewrite.Token, error) {
	out := []lexer_rewrite.Token{}

	for _, t := range tokens {
		switch t.Type {
		case lexer_rewrite.TokenTypeFileInclude:
			name := t.Value

			// get file
			f, err := a.getFile(name)
			if err != nil {
				return nil, err
			}

			// get tokens for file
			tokens, err := a.lexer.Run(name, f)
			if err != nil {
				return nil, err
			}

			// get tokens for files included in file
			includeTokens, err := a.getIncludeTokens(name, tokens)
			if err != nil {
				return nil, err
			}

			out = append(out, includeTokens...)
		case lexer_rewrite.TokenTypeSystemInclude:
			name := t.Value

			// get system include source
			source, exists := a.systemIncludeSources[name]
			if !exists {
				return nil, fmt.Errorf("no system include found for <%s>", name)
			}

			// get tokens for file
			tokens, err := a.lexer.Run(name, source)
			if err != nil {
				return nil, err
			}

			// get tokens for files included in file
			includeTokens, err := a.getIncludeTokens(name, tokens)
			if err != nil {
				return nil, err
			}

			out = append(out, includeTokens...)

		case lexer_rewrite.TokenTypeEndOfFile:
			// skip
		default:
			out = append(out, t)
		}
	}

	tokens = append(tokens, lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile})

	return out, nil
}

func (a *Assembler) getEntryPointTableTokens() ([]lexer_rewrite.Token, error) {
	buf := bytes.Buffer{}
	buf.WriteString("jump __start\n")

	for _, label := range a.config.InteruptLabels {
		if label != "" {
			buf.WriteString(fmt.Sprintf("jump %s\n", label))
		} else {
			buf.WriteString("db 0\n")
			buf.WriteString("db 0\n")
			buf.WriteString("db 0\n")
		}
	}

	return a.lexer.Run("", buf.String())
}

func (a *Assembler) GetProgram(filename string, source string) ([]uint8, error) {
	// get tokens for entry point table
	tokens := []lexer_rewrite.Token{}

	if !a.config.disableEntryPointsTable {
		entryPointTableTokens, err := a.getEntryPointTableTokens()
		if err != nil {
			return nil, err
		}

		entryPointTableTokens = entryPointTableTokens[:len(entryPointTableTokens)-1]

		tokens = append(tokens, entryPointTableTokens...)
	}

	// get tokens for entry point file
	entryPointTokens, err := a.lexer.Run(filename, source)
	if err != nil {
		return nil, err
	}

	// recursively gets tokens for each included file
	combinedTokens, err := a.getIncludeTokens(filename, entryPointTokens)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, combinedTokens...)

	p := NewParser()
	bin, err := p.Run(tokens)
	if err != nil {
		return nil, err
	}

	return bin, nil
}
