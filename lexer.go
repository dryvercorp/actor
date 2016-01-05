package actor

import (
	"bufio"
	"io"
	"strings"
)

type lineContent string
type lexerTree []*line

type line struct {
	line     int
	column   int
	content  lineContent
	children lexerTree
}

type lexer struct {
	reader io.Reader
	lines  lexerTree
}

func newLine(line_number, column int, content string) *line {
	return &line{
		line:     line_number,
		column:   column,
		content:  lineContent(content),
		children: make(lexerTree, 0),
	}
}

func (l *line) branch() *line {
	return newLine(l.line, l.content.indent(), strings.Trim(string(l.content), " \t"))
}

func newLexer(reader io.Reader) *lexer {

	lex := lexer{
		reader: reader,
		lines:  make(lexerTree, 0),
	}

	return &lex
}

func (s lineContent) indent() int {
	return len(s) - len(strings.TrimLeft(string(s), " \t"))
}

func (l *lexer) lex() (lines lexerTree, err error) {

	// Split to lines
	scanner := bufio.NewScanner(l.reader)
	scanner.Split(bufio.ScanLines)

	line_number := 0
	raw_lines := make(lexerTree, 0)

	for scanner.Scan() {
		line_number++

		text := strings.TrimRight(scanner.Text(), " \t")

		if text == "" {
			continue
		}

		raw_lines = append(raw_lines, newLine(line_number, 0, text))
	}

	// Make sure the first line has no indent
	if len(raw_lines) > 0 {
		index := 0
		l.indentLines(&index, raw_lines, &lines, raw_lines[0].content.indent())
	}

	return
}

func (l *lexer) indentLines(index *int, input lexerTree, output *lexerTree, indent int) {

	// Ends when there are no more lines
	if *index >= len(input) {
		return
	}

	var line_to_add *line

	for ; *index < len(input); *index++ {

		line_indent := input[*index].content.indent()

		if line_indent == indent {
			line_to_add = input[*index].branch()
			*output = append(*output, line_to_add)

		} else if line_indent > indent {
			l.indentLines(index, input, &line_to_add.children, line_indent)

		} else if line_indent < indent {
			*index--
			return
		}

	}

	return
}
