package main

import (
	"fmt"
	"github.com/cardrank/cardrank"
	"sort"
)

const (
	GameType    = cardrank.Badugi
	PlayerCount = 2
	DrawCount   = 1
)

func getChangeCount() []int {
	return []int{2, 1, 1}
}

type Board struct {
	Hand    []cardrank.Card
	Deck    []cardrank.Card
	Discard []cardrank.Card
}

type Boards []Board

func NewBoard(hand []cardrank.Card, deck []cardrank.Card, discard []cardrank.Card) Board {
	newDeck := removeCardsFromCards(deck, hand)

	return Board{hand, newDeck, discard}
}

func main() {
	deck := cardrank.NewDeck()

	boards := generateBoards(deck.All())

	resultBoards := make([]Board, 0)

	var foldCount = 0
	currentBoards := boards
	for i := 0; i < DrawCount; i++ {
		r, f := getDrawBoards(currentBoards, i+1)
		currentBoards = r
		foldCount += f
	}
	resultBoards = currentBoards

	hands := make([][]cardrank.Card, 0)
	for _, resultBoard := range resultBoards {
		hands = append(hands, resultBoard.Hand)
	}
	badugiHands := getBadugiHands(hands)
	badugiMap := make(map[cardrank.Rank]int, 0)
	for _, hand := range badugiHands {
		card := getHighCard(hand)
		badugiMap[card.Rank()]++
	}
	fmt.Printf("BadugiMap: %v\n", badugiMap)
	fmt.Printf("BadugiHands: %d\n", len(badugiHands))
	fmt.Printf("Boards: %d\n", len(resultBoards))
	fmt.Printf("FoldCount: %d\n", foldCount)
}

func getDrawBoards(boards []Board, count int) ([]Board, int) {
	resultBoards := make([]Board, 0)
	var foldCounts = 0
	for _, board := range boards {
		if count == 1 {
			fmt.Printf("Stating Hand: %s\n", board.Hand)
		} else {
			fmt.Printf("Current Hand: %s Draw count: %d\n", board.Hand, count)
		}
		discard := board.Discards()
		if len(discard) < 2 {
			drawBoards := board.Draw(discard)
			resultBoards = append(resultBoards, drawBoards...)
		} else {
			foldCounts++
		}
	}
	return resultBoards, foldCounts
}

func generateBoards(deck []cardrank.Card) []Board {
	var boards []Board
	//hand := []cardrank.Card{
	//	cardrank.New(cardrank.Ace, cardrank.Spade),
	//	cardrank.New(cardrank.Jack, cardrank.Spade),
	//	cardrank.New(cardrank.Queen, cardrank.Spade),
	//	cardrank.New(cardrank.King, cardrank.Spade),
	//}
	//
	//board := NewBoard(hand, deck, []cardrank.Card{})
	//boards = append(boards, board)

	cards := generateCombinations(deck, 4)

	for _, card := range cards {
		board := NewBoard(card, deck, []cardrank.Card{})
		boards = append(boards, board)
	}

	return boards
}

func getBadugiHands(hands [][]cardrank.Card) [][]cardrank.Card {
	var badugiHands [][]cardrank.Card
	for _, hand := range hands {
		run := &cardrank.Run{
			Pockets: [][]cardrank.Card{hand},
		}
		active := make(map[int]bool, 1)
		active[0] = true
		result := cardrank.NewResult(GameType, run, active, true)
		hi := result.Evals[0].Desc(false)
		if len(hi.Unused) == 0 {
			badugiHands = append(badugiHands, hand)
		}
	}
	return badugiHands
}

func removeCardsFromCards(cards []cardrank.Card, removeCards []cardrank.Card) []cardrank.Card {
	var newCards []cardrank.Card
	for _, card := range cards {
		if !containsCard(removeCards, card) {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

func containsCard(cards []cardrank.Card, card cardrank.Card) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}
	return false
}

func (b Board) Discards() []cardrank.Card {
	// スーツごとにカードを分類します。
	suitMap := make(map[cardrank.Suit][]cardrank.Card)
	for _, card := range b.Hand {
		suitMap[card.Suit()] = append(suitMap[card.Suit()], card)
	}

	// スーツごとにカードをランク順に並び替えます。
	for suit, cards := range suitMap {
		sort.Slice(cards, func(i, j int) bool {
			return convertRank(cards[i].Rank()) < convertRank(cards[j].Rank())
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
		if containsRank(b.Hand, minCard.Rank()) {
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
			if convertRank(minCard1.Rank()) != convertRank(minCard2.Rank()) {
				discardCards = append(discardCards, suitMap[suits[0]][1])
				discardCards = append(discardCards, suitMap[suits[1]][1])
			} else {
				if convertRank(suitMap[suits[0]][1].Rank()) < convertRank(suitMap[suits[1]][1].Rank()) {
					discardCards = append(discardCards, suitMap[suits[1]][1])
					discardCards = append(discardCards, suitMap[suits[0]][0])
				} else if convertRank(suitMap[suits[0]][1].Rank()) > convertRank(suitMap[suits[1]][1].Rank()) {
					discardCards = append(discardCards, suitMap[suits[0]][1])
					discardCards = append(discardCards, suitMap[suits[1]][0])
				} else {
					discardCards = append(discardCards, suitMap[suits[0]][1])
					discardCards = append(discardCards, suitMap[suits[1]][0])
				}
			}
		} else {
			minCard := suitMap[maxSuits][0]
			if containsRank(b.Hand, minCard.Rank()) {
				discardCards = append(discardCards, suitMap[maxSuits][1])
			} else {
				discardCards = append(discardCards, minCard)
				hand := removeCardsFromCards(b.Hand, discardCards)
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
		discardCards = removeCardsFromCards(b.Hand, keepCards)
	}

	return discardCards
}

func (b Board) Draw(discardCards []cardrank.Card) []Board {
	cards := generateCombinations(b.Deck, len(discardCards))

	var boards = make([]Board, 0)
	for _, card := range cards {
		board := NewBoard(b.Hand, b.Deck, b.Discard)
		board.Discard = append(board.Discard, discardCards...)
		board.Hand = removeCardsFromCards(board.Hand, discardCards)
		board.Hand = append(board.Hand, card...)
		board.Deck = removeCardsFromCards(board.Deck, card)

		boards = append(boards, board)
	}

	return boards
}

func containsSuit(suits []cardrank.Suit, suit cardrank.Suit) bool {
	for _, s := range suits {
		if s == suit {
			return true
		}
	}
	return false
}

func containsRank(cards []cardrank.Card, rank cardrank.Rank) bool {
	for _, card := range cards {
		if card.Rank() == rank {
			return true
		}
	}
	return false
}

func convertRank(rank cardrank.Rank) int {
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

func generateCombinations(deck []cardrank.Card, comboSize int) [][]cardrank.Card {
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

func getHighCard(hand []cardrank.Card) cardrank.Card {
	var highCard cardrank.Card
	for _, card := range hand {
		if convertRank(highCard.Rank()) < convertRank(card.Rank()) {
			highCard = card
		}
	}
	return highCard
}
