package main

import (
	"fmt"
	"strconv"
)

type board [9]string
type Winner struct {
	player string
	title  string
}

func printBoard(values board) {
	fmt.Print("\033[H\033[2J") // clears the screen

	for i, value := range values {
		fmt.Printf(" %-1s ", value)

		switch {
		case i > 7:
			fmt.Print("\n")
		case (i+1)%3 == 0:
			fmt.Print("\n───┼───┼───\n")
		default:
			fmt.Print("│")
		}
	}
}

func findWinnerXY(values board, posCalc func(int, int) int) string {
	var player string
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			value := values[posCalc(i, j)]
			if player == "" {
				player = value
			} else if player != value {
				player = ""
				break
			}
		}

		if player != "" {
			break
		}
	}

	return player
}

func findWinnerDiagonal(values board, posCalc func(int) int) string {
	var player string
	for i := 0; i < 3; i++ {
		value := values[posCalc(i)]
		if player == "" {
			player = value
		} else if player != value {
			player = ""
			break
		}
	}

	return player
}

func findByRow(values board) Winner {
	player := findWinnerXY(values, func(i, j int) int {
		return (i * 3) + j
	})

	return Winner{player: player, title: "row"}
}

func findByColumn(values board) Winner {
	player := findWinnerXY(values, func(i, j int) int {
		return (j * 3) + i
	})

	return Winner{player: player, title: "column"}
}

func findByDiagonal1(values board) Winner {
	player := findWinnerDiagonal(values, func(i int) int {
		return (i * 3) + i
	})

	return Winner{player: player, title: "diagonal"}
}
func findByDiagonal2(values board) Winner {
	player := findWinnerDiagonal(values, func(i int) int {
		return (i * 3) + (2 - i)
	})
	return Winner{player: player, title: "diagonal"}
}

func findWinner(values board) Winner {
	finders := make(chan Winner)

	go func(values board) { finders <- findByRow(values) }(values)
	go func(values board) { finders <- findByColumn(values) }(values)
	go func(values board) { finders <- findByDiagonal1(values) }(values)
	go func(values board) { finders <- findByDiagonal2(values) }(values)

	winner := Winner{}
	for i := 0; i < 4; i++ {
		f := <-finders
		if f.player != "" {
			winner = f
		}
	}

	return winner
}

func validateInput(in int, values board) (bool, string) {
	switch {
	case in < 0 || in > 8:
		return false, "Invalid position! Choose between 1 and 9"
	case values[in] == "x" || values[in] == "o":
		return false, "Position already taken!"
	default:
		return true, ""
	}
}

func main() {
	var players = [2]string{"x", "o"}

	var values board
	for i := range values {
		values[i] = strconv.Itoa(i + 1)
	}

	printBoard(values)
	for i := 0; i < 9; i++ {
		player := players[i%2]
		fmt.Printf("Player %s chose a position: ", player)

		var position int
		fmt.Scanf("%d", &position)
		position-- // user inputs the position indexed by 1

		valid, errMsg := validateInput(position, values)
		if !valid {
			i--
			fmt.Println(errMsg, "Try again")
			continue
		}

		values[position] = player
		printBoard(values)

		// it needs at least five rounds to have a winner
		if i >= 4 {
			w := findWinner(values)
			if w.player != "" {
				fmt.Printf("%s wins the game by %s", w.player, w.title)
				return
			}
		}
	}

	fmt.Println(
		`No one win,
		no one lost,
		at the end
		we're all winners
		who have lost`)
}
