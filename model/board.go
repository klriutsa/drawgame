package model

import (
	"drawgame/util"
	"github.com/cardrank/cardrank"
	"sort"
)

type Board struct {
	Hand    []cardrank.Card
	Deck    []cardrank.Card
	Discard []cardrank.Card
}

type Boards []Board

func NewBoard(hand []cardrank.Card, deck []cardrank.Card, discard []cardrank.Card) Board {
	newDeck := util.RemoveCardsFromCards(deck, hand)

	return Board{hand, newDeck, discard}
}

func NewBoards(deck []cardrank.Card) Boards {
	var boards Boards

	cards := util.GenerateCombinations(deck, 4)

	for _, card := range cards {
		board := NewBoard(card, deck, []cardrank.Card{})
		boards = append(boards, board)
	}

	return boards
}

func NewBoardsByHands(hands [][]cardrank.Card, deck []cardrank.Card) Boards {
	var boards Boards

	for _, hand := range hands {
		board := NewBoard(hand, deck, []cardrank.Card{})
		boards = append(boards, board)
	}

	return boards
}

func (b Board) Discards(minRank cardrank.Rank) []cardrank.Card {
	// スーツごとにカードを分類します。
	suitMap := make(map[cardrank.Suit][]cardrank.Card)
	for _, card := range b.Hand {
		suitMap[card.Suit()] = append(suitMap[card.Suit()], card)
	}

	// スーツごとにカードをランク順に並び替えます。
	for suit, cards := range suitMap {
		sort.Slice(cards, func(i, j int) bool {
			return util.ConvertRank(cards[i].Rank()) < util.ConvertRank(cards[j].Rank())
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
		hand := util.RemoveCardsFromCards(b.Hand, []cardrank.Card{minCard})
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
			if util.ConvertRank(minCard1.Rank()) != util.ConvertRank(minCard2.Rank()) {
				discardCards = append(discardCards, suitMap[suits[0]][1])
				discardCards = append(discardCards, suitMap[suits[1]][1])
			} else {
				if util.ConvertRank(suitMap[suits[0]][1].Rank()) < util.ConvertRank(suitMap[suits[1]][1].Rank()) {
					discardCards = append(discardCards, suitMap[suits[1]][1])
					discardCards = append(discardCards, suitMap[suits[0]][0])
				} else if util.ConvertRank(suitMap[suits[0]][1].Rank()) > util.ConvertRank(suitMap[suits[1]][1].Rank()) {
					discardCards = append(discardCards, suitMap[suits[0]][1])
					discardCards = append(discardCards, suitMap[suits[1]][0])
				} else {
					discardCards = append(discardCards, suitMap[suits[0]][1])
					discardCards = append(discardCards, suitMap[suits[1]][0])
				}
			}
		} else {
			minCard := suitMap[maxSuits][0]
			if util.ContainsRank(b.Hand, minCard.Rank()) {
				discardCards = append(discardCards, suitMap[maxSuits][1])
			} else {
				discardCards = append(discardCards, minCard)
				hand := util.RemoveCardsFromCards(b.Hand, discardCards)
				if hand[0].Rank() == hand[1].Rank() {
					discardCards = append(discardCards, hand[0])
				}
			}
		}
	} else if maxCount == 1 {
		var keepCards = make([]cardrank.Card, 0)
		var rankMap = make(map[cardrank.Rank][]cardrank.Card)
		for _, card := range b.Hand {
			rankMap[card.Rank()] = append(rankMap[card.Rank()], card)
		}
		for _, cards := range rankMap {
			if len(cards) > 1 {
				keepCards = append(keepCards, cards[0])
			} else {
				keepCards = append(keepCards, cards[0])
			}
		}
		discardCards = util.RemoveCardsFromCards(b.Hand, keepCards)
	}
	currentHand := util.RemoveCardsFromCards(b.Hand, discardCards)
	for _, card := range currentHand {
		if util.ConvertRank(card.Rank()) > util.ConvertRank(minRank) {
			discardCards = append(discardCards, card)
		}
	}
	return discardCards
}

func (b Board) Draw(discardCards []cardrank.Card) Boards {
	if len(discardCards) == 0 {
		return Boards{b}
	}
	b.Deck = util.RemoveCardsFromCards(b.Deck, discardCards)

	cards := util.GenerateCombinations(b.Deck, len(discardCards))

	var boards = make(Boards, 0)
	for _, card := range cards {
		board := NewBoard(b.Hand, b.Deck, b.Discard)
		board.Discard = append(board.Discard, discardCards...)
		board.Hand = util.RemoveCardsFromCards(board.Hand, discardCards)
		board.Hand = append(board.Hand, card...)
		board.Deck = util.RemoveCardsFromCards(board.Deck, card)

		boards = append(boards, board)
	}

	return boards
}
