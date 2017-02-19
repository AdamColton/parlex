// Package parlex provides common types and interfaces for describing the
// commonents of parsing. These are grammars, lexers and parsers. This allows
// for easy experimentation and changing of components. For instance, a new
// parser could be written that reuses the provided grammar and lexer.
//
// There are subpackages that provide an instance of a Grammar, Lexer and
// Parser.
//
// This package is currently an experiment with parsing. It is far from
// optimized. Feedback is welcome and encouraged.
package parlex
