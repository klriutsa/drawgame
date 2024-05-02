package game

import "github.com/cardrank/cardrank"

type Game interface {
	GetDiscard(hand []cardrank.Card, minRank cardrank.Rank) []cardrank.Card
	ShowHands(hands [][]cardrank.Card)
}
