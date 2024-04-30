package model

import (
	"testing"

	"github.com/cardrank/cardrank"

	"github.com/google/go-cmp/cmp"
)

func TestBoard_Discards(t *testing.T) {
	type args struct {
		minRank cardrank.Rank
	}
	tests := []struct {
		name string
		args args
		hand []cardrank.Card
		want []cardrank.Card
	}{
		{
			name: "1枚交換",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Club),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
		},
		{
			name: "1枚交換その２",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Club),
				cardrank.New(cardrank.Four, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
		},
		{
			name: "2枚交換",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
		},
		{
			name: "3枚交換",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Two, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Two, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
		},
		{
			name: "2枚交換、9以上は交換する",
			args: args{
				minRank: cardrank.Eight,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Club),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Club),
			},
		},
		{
			name: "Aceが2枚あるとき、2tone",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
			},
		},
		{
			name: "2tone",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Spade),
				cardrank.New(cardrank.Three, cardrank.Diamond),
				cardrank.New(cardrank.Four, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Two, cardrank.Spade),
				cardrank.New(cardrank.Four, cardrank.Diamond),
			},
		},
		{
			name: "Aceが2枚、2tone",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Spade),
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Four, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Four, cardrank.Diamond),
			},
		},
		{
			name: "Aceが2枚、2が2枚あるとき、2tone",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Two, cardrank.Spade),
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Two, cardrank.Diamond),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Two, cardrank.Spade),
				cardrank.New(cardrank.Ace, cardrank.Diamond),
			},
		},
		{
			name: "すべて同じカード",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Ace, cardrank.Club),
				cardrank.New(cardrank.Ace, cardrank.Heart),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Ace, cardrank.Club),
				cardrank.New(cardrank.Ace, cardrank.Heart),
			},
		},
		{
			name: "3枚同じカード",
			args: args{
				minRank: cardrank.King,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Spade),
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Ace, cardrank.Club),
				cardrank.New(cardrank.Two, cardrank.Heart),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Ace, cardrank.Diamond),
				cardrank.New(cardrank.Ace, cardrank.Club),
			},
		},
		{
			name: "すべて同じカード、9以上は交換する",
			args: args{
				minRank: cardrank.Eight,
			},
			hand: []cardrank.Card{
				cardrank.New(cardrank.Nine, cardrank.Spade),
				cardrank.New(cardrank.Nine, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Club),
				cardrank.New(cardrank.Nine, cardrank.Heart),
			},
			want: []cardrank.Card{
				cardrank.New(cardrank.Nine, cardrank.Diamond),
				cardrank.New(cardrank.Nine, cardrank.Club),
				cardrank.New(cardrank.Nine, cardrank.Heart),
				cardrank.New(cardrank.Nine, cardrank.Spade),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := cardrank.NewDeck()
			board := NewBoard(tt.hand, deck.All(), []cardrank.Card{})
			got := board.Discards(tt.args.minRank)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("User value is mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
