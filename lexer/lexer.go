package lexer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var sourceHasError bool = false

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type TokenType int

const (
	// single character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// ONE OR TWO CHARACTER TOKENS
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	//LITERALS
	IDENTIFIER
	STRING
	NUMBER

	//KEYWORDS
	VAR
	CONST
	FN
	RETURN
	IF
	ELSE
	FOR
	LOOP
	BREAK
	CONTINUE
	AND
	OR
	SWITCH
	CASE

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	line    int
}

func (t TokenType) EnumIndex() int {
	return int(t)
}

func (t TokenType) toString() string {
	return [...]string{"LEFT_PAREN",
		"RIGHT_PAREN",
		"LEFT_BRACE",
		"RIGHT_BRACE",
		"COMMA",
		"DOT",
		"MINUS",
		"PLUS",
		"SEMICOLON",
		"SLASH",
		"STAR",
		"BANG",
		"BANG_EQUAL",
		"EQUAL",
		"EQUAL_EQUAL",
		"GREATER",
		"GREATER_EQUAL",
		"LESS",
		"LESS_EQUAL",
		"IDENTIFIER",
		"STRING",
		"NUMBER",
		"VAR",
		"CONST",
		"FN",
		"RETURN",
		"IF",
		"ELSE",
		"FOR",
		"LOOP",
		"BREAK",
		"CONTINUE",
		"AND",
		"OR",
		"SWITCH",
		"CASE",
		"EOF",
	}[t]
}

var KeyWords = map[string]TokenType{
	"var":      VAR,
	"const":    CONST,
	"fn":       FN,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"loop":     LOOP,
	"break":    BREAK,
	"continue": CONTINUE,
	"and":      AND,
	"or":       OR,
	"switch":   SWITCH,
	"case":     CASE,
}

func report(line int, message string) {
	fmt.Printf("Error on line: %v : %v ", line, message)
}

func Error(line int, message string) {
	report(line, message)
	sourceHasError = true
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

// helper functions for scanner
func (scnr *Scanner) IsAtEnd() bool {
	return scnr.current >= len(scnr.source)
}

func (scnr *Scanner) IsADigit(char byte) bool {
	if char >= '0' && char <= '9' {
		return true
	} else {
		return false
	}
}

func (scnr *Scanner) IsAlpha(char byte) bool {
	if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '_' {
		return true
	} else {
		return false
	}
}

func (scnr *Scanner) IsAlphaNum(char byte) bool {
	return scnr.IsADigit(char) || scnr.IsAlpha(char)
}

func (scnr *Scanner) Advance() byte {
	current_byte := scnr.source[scnr.current]
	scnr.current = scnr.current + 1
	return current_byte
}

func (scnr *Scanner) AddToken(tokenType TokenType, lit string) {
	text := scnr.source[scnr.start:scnr.current]
	scnr.tokens = append(scnr.tokens, Token{tokenType, text, lit, scnr.line})
}

func (scnr *Scanner) AddTokenNoLiteral(tokenType TokenType) {
	scnr.AddToken(tokenType, "null")
}

func (scnr *Scanner) Match(expected byte) bool {
	if scnr.IsAtEnd() {
		return false
	}
	if scnr.source[scnr.current] != expected {
		return false
	}
	scnr.current++
	return true

}

func (scnr *Scanner) Peek() byte {
	if scnr.IsAtEnd() {
		return '\x00'
	}
	current_char := scnr.source[scnr.current]
	return current_char
}

func (scnr *Scanner) PeekNext() byte {
	if scnr.current+1 >= len(scnr.source) {
		return '\x00'
	}
	return scnr.source[scnr.current+1]
}

func (scnr *Scanner) AddString() {
	for scnr.Peek() != '"' && !scnr.IsAtEnd() {
		if scnr.Peek() == '\n' {
			scnr.line = scnr.line + 1
		}
		scnr.Advance()
	}

	if scnr.IsAtEnd() {
		Error(scnr.line, "Unterminated String")
		return
	}
	// closing "
	scnr.Advance()

	value := scnr.source[scnr.start+1 : scnr.current-1]
	scnr.AddToken(STRING, value)

}

func (scnr *Scanner) AddNumber(AlreadyFraction bool) {
	for scnr.IsADigit(scnr.Peek()) {
		scnr.Advance()
	}

	if !AlreadyFraction && scnr.Peek() == '.' && scnr.IsADigit(scnr.PeekNext()) {
		// consume the '.' if already not a fraction
		scnr.Advance()
		for scnr.IsADigit(scnr.Peek()) {
			scnr.Advance()
		}
	}

	scnr.AddToken(NUMBER, scnr.source[scnr.start:scnr.current])

}

func (scnr *Scanner) AddIdentifierOrKeyWord() {
	for scnr.IsAlphaNum(scnr.Peek()) {
		scnr.Advance()
	}

	text := scnr.source[scnr.start:scnr.current]
	typeOfText, isKeyWord := KeyWords[text]

	if isKeyWord {
		scnr.AddTokenNoLiteral(typeOfText)
	} else {
		scnr.AddToken(IDENTIFIER, text)
	}

}

// main scan function

func (scnr *Scanner) ScanToken() {
	character := scnr.Advance()

	switch character {
	case '(':
		scnr.AddTokenNoLiteral(LEFT_PAREN)
	case ')':
		scnr.AddTokenNoLiteral(RIGHT_PAREN)
	case '{':
		scnr.AddTokenNoLiteral(LEFT_BRACE)
	case '}':
		scnr.AddTokenNoLiteral(RIGHT_BRACE)
	case ',':
		scnr.AddTokenNoLiteral(COMMA)
	case '.':
		if scnr.IsADigit(scnr.PeekNext()) {
			scnr.AddNumber(true)
		}
		scnr.AddTokenNoLiteral(DOT)
	case '-':
		scnr.AddTokenNoLiteral(MINUS)
	case '+':
		scnr.AddTokenNoLiteral(PLUS)
	case ';':
		scnr.AddTokenNoLiteral(SEMICOLON)
	case '*':
		scnr.AddTokenNoLiteral(STAR)
	case '!':
		if scnr.Match('=') {
			scnr.AddTokenNoLiteral(BANG_EQUAL)
		} else {
			scnr.AddTokenNoLiteral(BANG)
		}
	case '=':
		if scnr.Match('=') {
			scnr.AddTokenNoLiteral(EQUAL_EQUAL)
		} else {
			scnr.AddTokenNoLiteral(EQUAL)
		}
	case '<':
		if scnr.Match('=') {
			scnr.AddTokenNoLiteral(LESS_EQUAL)
		} else {
			scnr.AddTokenNoLiteral(LESS)
		}
	case '>':
		if scnr.Match('=') {
			scnr.AddTokenNoLiteral(GREATER_EQUAL)
		} else {
			scnr.AddTokenNoLiteral(GREATER)
		}
	case '/':
		if scnr.Match('/') {
			// comment goes till the end of the line
			for scnr.Peek() != '\n' && !scnr.IsAtEnd() {
				scnr.Advance()
			}
		} else {
			scnr.AddTokenNoLiteral(SLASH)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		scnr.line++
	case '"':
		scnr.AddString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		scnr.AddNumber(false)
	default:
		if scnr.IsAlpha(character) {
			scnr.AddIdentifierOrKeyWord()
		} else {
			Error(scnr.line, "Unexpected Character")
		}

	}
}

func (scnr *Scanner) ScanTokensFromSource() []Token {

	for !scnr.IsAtEnd() {
		scnr.start = scnr.current
		scnr.ScanToken()
	}
	EOFToken := Token{EOF, "", "null", scnr.line}
	scnr.tokens = append(scnr.tokens, EOFToken)

	return scnr.tokens

}

func RunLexer(source string) []Token {
	scnr := Scanner{source, []Token{}, 0, 0, 1}
	outputTokens := scnr.ScanTokensFromSource()
	return outputTokens

}

func RunLexerFromFile(filePath string) []Token {
	inputFile, err := os.ReadFile(filePath)
	check(err)

	outputTokens := RunLexer(string(inputFile))
	return outputTokens
}

func RunLexerPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Mayfair lexer REPL")
	for {
		fmt.Print("\n-->> ")
		line, err := reader.ReadString('\n')
		check(err)
		if strings.TrimSpace(line) == "quit" {
			fmt.Println("Thanks for using Lex REPL")
			break
		}
		PrintTokens(RunLexer(line))
	}
}

func PrintTokens(tokens []Token) {
	fmt.Println("Line,TokenType,Literal,Lexeme")
	for i := range tokens {
		fmt.Printf("%v,%v,%v,%v\n", tokens[i].line, tokens[i].Type.toString(), tokens[i].Literal, tokens[i].Lexeme)
	}
}
