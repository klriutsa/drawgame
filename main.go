package main

import (
	"github.com/cardrank/cardrank"
	"math/rand"
	"time"

	"drawgame/model"
)

const (
	GameType    = cardrank.Badugi
	PlayerCount = 2
	DrawCount   = 1
	ChangeCount = 1
	MinRank     = cardrank.King
)

func getChangeCount() []int {
	return []int{2, 1, 1}
}

func main() {
	deck := cardrank.NewDeck()
	r := rand.New(rand.New(rand.NewSource(time.Now().UnixNano())))
	deck.Shuffle(r, 1)

	boards := model.NewBoards(GameType, deck.All())

	//boards := model.NewBoardsByHands(GameType, [][]cardrank.Card{
	//	{
	//		cardrank.New(cardrank.Ace, cardrank.Spade),
	//		cardrank.New(cardrank.Two, cardrank.Diamond),
	//		cardrank.New(cardrank.Eight, cardrank.Club),
	//		cardrank.New(cardrank.Four, cardrank.Diamond),
	//	},
	//}, deck.All())

	boards.ExecDraw(GameType, DrawCount, ChangeCount, MinRank)
}
