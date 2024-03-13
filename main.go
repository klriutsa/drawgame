package main

import (
	"drawgame/model"
	"drawgame/util"
	"fmt"
	"github.com/cardrank/cardrank"
)

const (
	GameType    = cardrank.Badugi
	PlayerCount = 2
	DrawCount   = 1
)

func getChangeCount() []int {
	return []int{2, 1, 1}
}

func main() {
	deck := cardrank.NewDeck()

	boards := model.NewBoards(deck.All())

	resultBoards := make([]model.Board, 0)

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
		card := util.GetHighCard(hand)
		badugiMap[card.Rank()]++
	}
	fmt.Printf("BadugiMap: %v\n", badugiMap)
	fmt.Printf("BadugiHands: %d\n", len(badugiHands))
	fmt.Printf("Boards: %d\n", len(resultBoards))
	fmt.Printf("FoldCount: %d\n", foldCount)
}

func getDrawBoards(boards []model.Board, count int) (model.Boards, int) {
	resultBoards := make(model.Boards, 0)
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
