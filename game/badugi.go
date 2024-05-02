package game

import (
	"fmt"
	"sort"

	"drawgame/util"

	"github.com/cardrank/cardrank"
)

type Badugi struct{}

func NewBadugi() Badugi {
	return Badugi{}
}

func (b Badugi) GetDiscard(hand []cardrank.Card, minRank cardrank.Rank) []cardrank.Card {
	// スーツごとにカードを分類します。
	suitMap := make(map[cardrank.Suit][]cardrank.Card)
	for _, card := range hand {
		suitMap[card.Suit()] = append(suitMap[card.Suit()], card)
	}

	// スーツごとにカードをランク順に並び替えます。
	for suit, cards := range suitMap {
		sort.Slice(cards, func(i, j int) bool {
			return b.convertRank(cards[i].Rank()) < b.convertRank(cards[j].Rank())
		})
		suitMap[suit] = cards
	}

	var maxCount int = 0
	var maxSuits cardrank.Suit
	var discardCards = make([]cardrank.Card, 0)
	for _, cards := range suitMap {
		if len(cards) > maxCount {
			maxCount = len(cards)
			maxSuits = cards[0].Suit()
		}
	}

	if maxCount == 4 {
		discardCards = append(discardCards, suitMap[maxSuits][1])
		discardCards = append(discardCards, suitMap[maxSuits][2])
		discardCards = append(discardCards, suitMap[maxSuits][3])
	} else if maxCount == 3 {
		minCard := suitMap[maxSuits][0]
		hand := util.RemoveCardsFromCards(hand, []cardrank.Card{minCard})
		if util.ContainsRank(hand, minCard.Rank()) {
			discardCards = append(discardCards, suitMap[maxSuits][0])
			discardCards = append(discardCards, suitMap[maxSuits][2])
		} else {
			discardCards = append(discardCards, suitMap[maxSuits][1])
			discardCards = append(discardCards, suitMap[maxSuits][2])
		}
	} else if maxCount == 2 {
		var suits = make([]cardrank.Suit, 0)
		var isDoubleSuit = false
		var isDouble = false
		for _, cards := range suitMap {
			if len(cards) == 2 {
				if isDouble {
					isDoubleSuit = true
				}
				isDouble = true
			}
		}
		if isDoubleSuit {
			for suit, _ := range suitMap {
				suits = append(suits, suit)
			}
			minCard1 := suitMap[suits[0]][0]
			minCard2 := suitMap[suits[1]][0]
			if b.convertRank(minCard1.Rank()) != b.convertRank(minCard2.Rank()) {
				discardCards = append(discardCards, suitMap[suits[0]][1])
				discardCards = append(discardCards, suitMap[suits[1]][1])
			} else {
				if b.convertRank(suitMap[suits[0]][1].Rank()) < b.convertRank(suitMap[suits[1]][1].Rank()) {
					discardCards = append(discardCards, suitMap[suits[1]][1])
					discardCards = append(discardCards, suitMap[suits[0]][0])
				} else if b.convertRank(suitMap[suits[0]][1].Rank()) > b.convertRank(suitMap[suits[1]][1].Rank()) {
					discardCards = append(discardCards, suitMap[suits[0]][1])
					discardCards = append(discardCards, suitMap[suits[1]][0])
				} else {
					discardCards = append(discardCards, suitMap[suits[0]][1])
					discardCards = append(discardCards, suitMap[suits[1]][0])
				}
			}
		} else {
			minCard := suitMap[maxSuits][0]
			if util.ContainsRank(hand, minCard.Rank()) {
				discardCards = append(discardCards, suitMap[maxSuits][1])
			} else {
				discardCards = append(discardCards, minCard)
				hand := util.RemoveCardsFromCards(hand, discardCards)
				if hand[0].Rank() == hand[1].Rank() {
					discardCards = append(discardCards, hand[0])
				}
			}
		}
	} else if maxCount == 1 {
		var keepCards = make([]cardrank.Card, 0)
		var rankMap = make(map[cardrank.Rank][]cardrank.Card)
		for _, card := range hand {
			rankMap[card.Rank()] = append(rankMap[card.Rank()], card)
		}
		for _, cards := range rankMap {
			if len(cards) > 1 {
				keepCards = append(keepCards, cards[0])
			} else {
				keepCards = append(keepCards, cards[0])
			}
		}
		discardCards = util.RemoveCardsFromCards(hand, keepCards)
	}
	currentHand := util.RemoveCardsFromCards(hand, discardCards)
	for _, card := range currentHand {
		if b.convertRank(card.Rank()) > b.convertRank(minRank) {
			discardCards = append(discardCards, card)
		}
	}
	return discardCards
}

func (b Badugi) convertRank(rank cardrank.Rank) int {
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

func (b Badugi) getHighCard(hand []cardrank.Card) cardrank.Card {
	var highCard cardrank.Card
	for _, card := range hand {
		if b.convertRank(highCard.Rank()) < b.convertRank(card.Rank()) {
			highCard = card
		}
	}
	return highCard
}

func (b Badugi) ShowHands(hands [][]cardrank.Card) {
	badugiMap, daiMap := b.collectHands(hands)
	fmt.Printf("BadugiMap: %v\n", badugiMap)
	for key, value := range badugiMap {
		fmt.Printf("%s:%d\n", key, value)
	}
	fmt.Printf("DaiMap: %v\n", daiMap)
	for key, value := range daiMap {
		fmt.Printf("%s:%d\n", key, value)
	}
}

func (b Badugi) collectHands(hands [][]cardrank.Card) (map[cardrank.Rank]int, map[string]int) {
	badugiHands, daiHands := b.getBadugiHands(hands)
	fmt.Printf("BadugiHands: %d\n", len(badugiHands))
	fmt.Printf("DaiHands: %d\n", len(daiHands))
	badugiMap := make(map[cardrank.Rank]int, 0)
	for _, hand := range badugiHands {
		card := b.getHighCard(hand)
		badugiMap[card.Rank()]++
	}
	daiMap := make(map[string]int, 0)
	for _, hand := range daiHands {
		card := b.getHighCard(hand)
		key := fmt.Sprintf("%d-%s", len(hand), card.Rank())
		daiMap[key]++
	}

	return badugiMap, daiMap
}

func (b Badugi) getBadugiHands(hands [][]cardrank.Card) ([][]cardrank.Card, [][]cardrank.Card) {
	var badugiHands [][]cardrank.Card
	var daiHands [][]cardrank.Card
	for _, hand := range hands {
		run := &cardrank.Run{
			Pockets: [][]cardrank.Card{hand},
		}
		active := make(map[int]bool, 1)
		active[0] = true
		result := cardrank.NewResult(cardrank.Badugi, run, active, true)
		hi := result.Evals[0].Desc(false)
		if len(hi.Unused) == 0 {
			badugiHands = append(badugiHands, hand)
		} else {
			daiHands = append(daiHands, hi.Best)
		}
	}
	return badugiHands, daiHands
}
