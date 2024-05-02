package game

import (
	"fmt"
	"github.com/cardrank/cardrank"
	"sort"
	"strings"

	"drawgame/util"
)

type DuceSeven struct {
}

func NewDuceSeven() DuceSeven {
	return DuceSeven{}
}

func (d DuceSeven) GetDiscard(hand []cardrank.Card, minRank cardrank.Rank) []cardrank.Card {
	// ランクごとにカードを分類する
	rankMap := make(map[cardrank.Rank][]cardrank.Card)
	for _, card := range hand {
		rankMap[card.Rank()] = append(rankMap[card.Rank()], card)
	}

	// スーツごとにカードを分類します。
	suitMap := make(map[cardrank.Suit][]cardrank.Card)
	for _, card := range hand {
		suitMap[card.Suit()] = append(suitMap[card.Suit()], card)
	}

	var discardCards = make([]cardrank.Card, 0)

	for _, cards := range rankMap {
		num := len(cards)
		if num >= 2 {
			for i := 0; i < num-1; i++ {
				discardCards = append(discardCards, cards[i])
			}
		}
	}

	currentHand := util.RemoveCardsFromCards(hand, discardCards)
	for _, card := range currentHand {
		if card.Rank() > minRank {
			discardCards = append(discardCards, card)
		}
	}

	return discardCards
}

func (d DuceSeven) ShowHands(hands [][]cardrank.Card) {
	lowballMap, otherMap := d.collectHands(hands)
	fmt.Printf("lowballMap: %v\n", lowballMap)
	for key, value := range lowballMap {
		fmt.Printf("%s:%d\n", key, value)
	}
	fmt.Printf("otherMap: %v\n", otherMap)
	for key, value := range otherMap {
		fmt.Printf("%s:%d\n", key, value)
	}
}

func (d DuceSeven) collectHands(hands [][]cardrank.Card) (map[string]int, map[string]int) {
	lowballHands, otherHands := d.getLowballHands(hands)
	fmt.Printf("lowballHands: %d\n", len(lowballHands))
	fmt.Printf("otherHands: %d\n", len(otherHands))
	lowballMap := make(map[string]int, 0)
	for _, hand := range lowballHands {
		card1, card2 := d.getHighCard(hand)
		key := fmt.Sprintf("%s-%s", card1.Rank(), card2.Rank())
		if card1.Rank() > cardrank.Nine {
			key = fmt.Sprintf("%s", card1.Rank())
		}
		lowballMap[key]++
	}
	otherMap := make(map[string]int, 0)

	return lowballMap, otherMap
}

func (d DuceSeven) getLowballHands(hands [][]cardrank.Card) ([][]cardrank.Card, [][]cardrank.Card) {
	var lowballHands [][]cardrank.Card
	var otherHands [][]cardrank.Card
	for _, hand := range hands {
		run := &cardrank.Run{
			Pockets: [][]cardrank.Card{hand},
		}
		active := make(map[int]bool, 1)
		active[0] = true
		result := cardrank.NewResult(cardrank.Lowball, run, active, true)
		hi := result.Evals[0].Desc(false)
		//fmt.Printf("hand: %s\n", hand)
		//fmt.Printf("hi: %s\n", hi)
		if strings.Contains(fmt.Sprintf("%v", hi), "low") {
			lowballHands = append(lowballHands, hi.Best)
		} else {
			otherHands = append(otherHands, hi.Best)
		}
	}
	return lowballHands, otherHands
}

func (d DuceSeven) getHighCard(hand []cardrank.Card) (cardrank.Card, cardrank.Card) {
	sort.Slice(hand, func(i, j int) bool {
		return hand[i].Rank() > hand[j].Rank()
	})

	return hand[0], hand[1]
}
