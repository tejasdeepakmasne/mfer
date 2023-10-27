# mfer
interpreter (and maybe) compiler for the mayfair toylang

# Installation
## Dependencies
- golang

## process
1. Clone the repository and cd into the folder
2. run `go build`
3. copy the binary to any location in your `$PATH`
4. Run the binary

# Usage

The interpreter currently supports only lexing operations.
To run lexing operations the following commands can be used :

1. To run the lexer on a file run:
   `mfer lex <filepath>`
2. To start the REPL loop run :
   `mfer lex`
