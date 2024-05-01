package util

import (
	"github.com/cardrank/cardrank"
)

func RemoveCardsFromCards(cards []cardrank.Card, removeCards []cardrank.Card) []cardrank.Card {
	var newCards []cardrank.Card
	for _, card := range cards {
		if !ContainsCard(removeCards, card) {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

func ContainsCard(cards []cardrank.Card, card cardrank.Card) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}
	return false
}

func ContainsSuit(suits []cardrank.Suit, suit cardrank.Suit) bool {
	for _, s := range suits {
		if s == suit {
			return true
		}
	}
	return false
}

func ContainsRank(cards []cardrank.Card, rank cardrank.Rank) bool {
	for _, card := range cards {
		if card.Rank() == rank {
			return true
		}
	}
	return false
}

func GenerateCombinations(deck []cardrank.Card, comboSize int) [][]cardrank.Card {
	var result [][]cardrank.Card
	combos := make([]cardrank.Card, comboSize)
	var doGenerate func(int, int)
	doGenerate = func(start, index int) {
		if index == comboSize {
			// 組み合わせのコピーを作成し、結果に追加します。
			combo := make([]cardrank.Card, comboSize)
			copy(combo, combos)
			result = append(result, combo)
			return
		}
		for i := start; i <= len(deck)-(comboSize-index); i++ {
			combos[index] = deck[i]
			doGenerate(i+1, index+1)
		}
	}
	doGenerate(0, 0)
	return result
}
