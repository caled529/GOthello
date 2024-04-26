package cli

import (
	"fmt"
	"strconv"

	"github.com/caled529/gothello/othello"
)

type Move = othello.Move
type Tile = othello.Tile

func RunGame() {
	board := othello.GetStartingBoard()
	currentPlayer := othello.Dark
	legalMoves := othello.CalcLegalMoves(board, currentPlayer)

	for playing := true; playing; {
		othello.PropagateLegalMoves(board, legalMoves)

		PrintGame(board, currentPlayer)

		move, quit := GetMove(legalMoves)
		if quit {
			fmt.Println("Goodbye!")
			return
		}
		othello.PlaceTile(move, board, currentPlayer)

		// This might be excessive but I'm trying to remove as much Othello logic
		// as possible from this main program
		legalMoves, currentPlayer, playing = othello.EndOfTurn(board, legalMoves, currentPlayer)
	}

	ClearTerm()
	fmt.Print(othello.BoardToString(board))
	if winner := othello.CalcWinner(board); winner == othello.Empty {
		fmt.Println("Tie game!")
	} else {
		fmt.Printf("%s player wins!\n", othello.PlayerString(winner))
	}
}

func PrintGame(board *[8][8]Tile, currentPlayer Tile) {
	ClearTerm()
	fmt.Print(othello.BoardToString(board))
	fmt.Printf("    Current player: %s\n\n", othello.PlayerString(currentPlayer))
}

func ClearTerm() {
	fmt.Printf("\033[2J\033[H")
}

// GetMove prints out a list of legal moves and allows the current player to
// select from one of them, returns the selected legal move.
func GetMove(legalMoves []Move) (Move, bool) {
	fmt.Println("Your legal moves:")
	for i, move := range legalMoves {
		fmt.Printf("  %d. %v\n", i+1, move)
	}
	fmt.Printf("Enter a number to select your move (0 to quit) >> ")
	choice := GetUserInt(0, len(legalMoves))
	if choice == 0 {
		return Move{}, true
	}
	return legalMoves[choice-1], false
}

// GetUserInt reads input from a scanner and returns the result if it is an int
// in the range low <= x <= high. Prompts the user to attempt again if input
// cannot be intepreted as an int in given range.
func GetUserInt(low int, high int) int {
	var input string
	fmt.Scanln(&input)
	inputInt, err := strconv.Atoi(input)
	for err != nil || low > inputInt || inputInt > high {
		fmt.Printf("Error, please enter a number in the range %d to %d >> ", low, high)
		fmt.Scanln(&input)
		inputInt, err = strconv.Atoi(input)
	}
	return inputInt
}
