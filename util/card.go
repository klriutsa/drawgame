package util

import "github.com/cardrank/cardrank"

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

func ConvertRank(rank cardrank.Rank) int {
	switch rank {
	case cardrank.Ace:
		return 1
	case cardrank.Two:
		return 2
	case cardrank.Three:
		return 3
	case cardrank.Four:
		return 4
	case cardrank.Five:
		return 5
	case cardrank.Six:
		return 6
	case cardrank.Seven:
		return 7
	case cardrank.Eight:
		return 8
	case cardrank.Nine:
		return 9
	case cardrank.Ten:
		return 10
	case cardrank.Jack:
		return 11
	case cardrank.Queen:
		return 12
	case cardrank.King:
		return 13
	default:
		return 0
	}
}

func GetHighCard(hand []cardrank.Card) cardrank.Card {
	var highCard cardrank.Card
	for _, card := range hand {
		if ConvertRank(highCard.Rank()) < ConvertRank(card.Rank()) {
			highCard = card
		}
	}
	return highCard
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
