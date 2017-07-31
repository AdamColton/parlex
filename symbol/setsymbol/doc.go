// Package setsymbol provides a implementations of parlex.Symbol,
// parlex.Production and parlex.Productions. The underlying representation is
// as indexed members of a set. It was implemented to reduce the number of maps
// with strings as keys. Many of those maps can be replaced with slices.
package setsymbol
