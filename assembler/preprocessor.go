package assembler

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type FileGetter func(path string) (string, error)

func FileSystemFileGetter(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

type Preprocessor struct {
	getFile FileGetter
}

func NewPreprocessor(fileGetter FileGetter) *Preprocessor {
	return &Preprocessor{getFile: fileGetter}
}

func splitAndRemoveWhiteSpace(line string) []string {
	split := strings.Split(line, " ")

	out := []string{}

	for _, s := range split {
		if strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}

	return out
}

func parseInclude(line string) (isFile bool, name string, err error) {
	parts := splitAndRemoveWhiteSpace(line)
	name = parts[1]

	switch name[0] {
	case '<':
		return false, name[1 : len(name)-1], nil
	case '"':
		return true, name[1 : len(name)-1], nil
	default:
		return false, "", fmt.Errorf("expected next char to be \" or < and got %s", string(name[0]))
	}
}

func (p *Preprocessor) getIncludes(source string) (map[string]string, error) {
	lines := strings.Split(source, "\n")

	includes := map[string]string{}

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#include") {
			isFile, name, err := parseInclude(line)
			if err != nil {
				return nil, fmt.Errorf("preprocessor failed: %w", err)
			}

			if isFile {
				file, err := p.getFile(name)
				if err != nil {
					return nil, err
				}

				fileSource, err := p.Process(file)
				if err != nil {
					return nil, err
				}

				includes[line] = fileSource

			} else {
				// handle "system" include
			}
		}
	}

	return includes, nil
}

func (p *Preprocessor) Process(source string) (string, error) {
	includes, err := p.getIncludes(source)
	if err != nil {
		return "", err
	}

	for line, file := range includes {
		source = strings.ReplaceAll(source, line, file)
	}

	return source, nil
}
