// Package regexgram provides a way to build a grammar using some of the
// operators common to regular expressions.
//
// Any operator can be applied to a single symbol or a group of symbols. A group
// can be formed by surrounding symbols or sub-groups in parenthesis.
//
// A symbol or group is marked as optional by following it with ?
//
// A symbol or group is marked as repeating using *
//
// A set of symbols or groups can be OR'd together with |
//
// The grammar also allows for full comments with //
package regexgram
