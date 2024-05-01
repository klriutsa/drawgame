package model

import (
	"fmt"

	badugi "drawgame/game"
	"drawgame/util"

	"github.com/cardrank/cardrank"
)

type Board struct {
	Game    cardrank.Type
	Hand    []cardrank.Card
	Deck    []cardrank.Card
	Discard []cardrank.Card
}

type Boards []Board

func NewBoard(game cardrank.Type, hand []cardrank.Card, deck []cardrank.Card, discard []cardrank.Card) Board {
	newDeck := util.RemoveCardsFromCards(deck, hand)

	return Board{game, hand, newDeck, discard}
}

func NewBoards(game cardrank.Type, deck []cardrank.Card) Boards {
	var boards Boards

	var handSize = 0

	if game == cardrank.Badugi {
		handSize = 4
	}

	cards := util.GenerateCombinations(deck, handSize)

	for _, card := range cards {
		board := NewBoard(game, card, deck, []cardrank.Card{})
		boards = append(boards, board)
	}

	return boards
}

func NewBoardsByHands(game cardrank.Type, hands [][]cardrank.Card, deck []cardrank.Card) Boards {
	var boards Boards

	for _, hand := range hands {
		board := NewBoard(game, hand, deck, []cardrank.Card{})
		boards = append(boards, board)
	}

	return boards
}

func (b Board) Discards(minRank cardrank.Rank) []cardrank.Card {
	if b.Game == cardrank.Badugi {
		return badugi.GetDiscard(b.Hand, minRank)
	}
	return []cardrank.Card{}
}

func (b Board) Draw(discardCards []cardrank.Card) Boards {
	if len(discardCards) == 0 {
		return Boards{b}
	}
	b.Deck = util.RemoveCardsFromCards(b.Deck, discardCards)

	cards := util.GenerateCombinations(b.Deck, len(discardCards))

	var boards = make(Boards, 0)
	for _, card := range cards {
		board := NewBoard(b.Game, b.Hand, b.Deck, b.Discard)
		board.Discard = append(board.Discard, discardCards...)
		board.Hand = util.RemoveCardsFromCards(board.Hand, discardCards)
		board.Hand = append(board.Hand, card...)
		board.Deck = util.RemoveCardsFromCards(board.Deck, card)

		boards = append(boards, board)
	}

	return boards
}

func (b Boards) ExecDraw(game cardrank.Type, drawCount, changeCount int, minRank cardrank.Rank) {
	var foldCount = 0
	currentBoards := b

	for i := 0; i < drawCount; i++ {
		r, f := b.GetDrawBoards(drawCount, changeCount, minRank)
		currentBoards = r
		foldCount += f
	}

	if game == cardrank.Badugi {
		hands := currentBoards.GetHands()

		badugi.ShowHands(hands)
	}
}

func (b Boards) GetHands() [][]cardrank.Card {
	hands := make([][]cardrank.Card, 0)
	for _, board := range b {
		hands = append(hands, board.Hand)
	}
	return hands
}

func (b Boards) GetDrawBoards(drawCount, changeCount int, minRank cardrank.Rank) (Boards, int) {
	resultBoards := make(Boards, 0)
	var foldCounts = 0
	for _, board := range b {
		if drawCount == 1 {
			fmt.Printf("Stating Hand: %s\n", board.Hand)
		} else {
			fmt.Printf("Current Hand: %s Draw count: %d\n", board.Hand, drawCount)
		}
		discard := board.Discards(minRank)
		if len(discard) == changeCount {
			drawBoards := board.Draw(discard)
			resultBoards = append(resultBoards, drawBoards...)
		} else {
			foldCounts++
		}
	}
	return resultBoards, foldCounts
}
