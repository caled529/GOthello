package main

import "fmt"

type Vec struct {
	x int
	y int
}

type Tile int

const (
	Empty Tile = iota
	Dark
	Light
	Legal
)

func main() {
	board := [8][8]Tile{}
	board[3][3] = Light
	board[3][4] = Dark
	board[4][3] = Dark
	board[4][4] = Light
	currentPlayer := Dark
	legalMoves := calcLegalMoves(&board, currentPlayer)
	for {
		propagateLegalMoves(&board, legalMoves)
		fmt.Printf("Current player: %s\n", playerString(currentPlayer))
		fmt.Print(boardToString(&board))
		placeTile(getMove(legalMoves), &board, currentPlayer)
		fmt.Println()
		clearLegalMoves(&board, legalMoves)
		currentPlayer = oppositePlayer(currentPlayer)
		legalMoves = calcLegalMoves(&board, currentPlayer)
		// If the current player has no legal moves, their turn is skipped
		if len(legalMoves) == 0 {
			currentPlayer = oppositePlayer(currentPlayer)
			legalMoves = calcLegalMoves(&board, currentPlayer)
			// If neither player has any legal moves, the game ends
			if len(legalMoves) == 0 {
				break
			}
		}
	}
	winner := calcWinner(&board)
	if winner == Empty {
		fmt.Println("Tie game!")
	} else {
		fmt.Printf("%s player wins!\n", playerString(winner))
	}

}

func oppositePlayer(currentPlayer Tile) Tile {
	if currentPlayer == Dark {
		return Light
	}
	return Dark
}

func playerString(player Tile) string {
	if player == Dark {
		return "Dark"
	}
	return "Light"
}

func propagateLegalMoves(board *[8][8]Tile, legalMoves []Vec) {
	for _, move := range legalMoves {
		board[move.x][move.y] = Legal
	}
}

func clearLegalMoves(board *[8][8]Tile, legalMoves []Vec) {
	for _, move := range legalMoves {
		// Check that Tile state is still legal to avoid clearing player move
		if board[move.x][move.y] == Legal {
			board[move.x][move.y] = Empty
		}
	}
}

// Need to document
func calcLegalMoves(board *[8][8]Tile, currentPlayer Tile) []Vec {
	var legalMoves []Vec
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			// Moves can only be made on Empty tiles
			if board[x][y] != Empty {
				continue
			}
			pos := Vec{x, y}
			for _, move := range circleTileCheck(pos, board, oppositePlayer(currentPlayer)) {
				dir := Vec{move.x - x, move.y - y}
				if lineTileCheck(pos, dir, board, currentPlayer) {
					legalMoves = append(legalMoves, Vec{x, y})
					break // Don't need to check legality multiple times
				}
			}
		}
	}
	return legalMoves
}

// circleTileCheck looks at the 8 tiles surrounding pos, then returns a slice
// of Vec containing coordinates of Tiles equal to target.
func circleTileCheck(pos Vec, board *[8][8]Tile, target Tile) []Vec {
	var matches []Vec
	for i := max(0, pos.x-1); i < min(len(board), pos.x+2); i++ {
		for j := max(0, pos.y-1); j < min(len(board[i]), pos.y+2); j++ {
			if board[i][j] == target {
				matches = append(matches, Vec{i, j})
			}
		}
	}
	return matches
}

// lineTileCheck steps from a position on the board (move) in increments of dir
// and returns true if a continous line of placed tiles can be made from pos to
// a Tile equal to target.
func lineTileCheck(pos Vec, dir Vec, board *[8][8]Tile, target Tile) bool {
	// One step made outside of loop to check if same color tiles are adjacent
	step := Vec{pos.x + dir.x, pos.y + dir.y}
	if board[step.x][step.y] == target {
		return false
	}
	for {
		step = Vec{step.x + dir.x, step.y + dir.y}
		// read: if out of bounds
		if step.x < 0 || step.x >= len(board) || step.y < 0 || step.y >= len(board[0]) {
			return false
		}
		if board[step.x][step.y] == Empty {
			return false
		}
		if board[step.x][step.y] == target {
			return true
		}
	}
}

// placeTile changes state of Tile of board at move to currentPlayer,
// and calls flipWalk to "flip" all opposite Tiles to the state of
// currentPlayer, in accordance with the game rules
// (https://wikipedia.org/wiki/Reversi#Rules).
func placeTile(move Vec, board *[8][8]Tile, currentPlayer Tile) {
	oppositePlayer := oppositePlayer(currentPlayer)
	board[move.x][move.y] = currentPlayer
	for _, opp := range circleTileCheck(move, board, oppositePlayer) {
		dir := Vec{opp.x - move.x, opp.y - move.y}
		if lineTileCheck(move, dir, board, currentPlayer) {
			flipWalk(move, dir, board, currentPlayer)
		}
	}
}

func flipWalk(move Vec, dir Vec, board *[8][8]Tile, currentPlayer Tile) {
	step := move
	for {
		step = Vec{step.x + dir.x, step.y + dir.y}
		if board[step.x][step.y] != currentPlayer {
			board[step.x][step.y] = currentPlayer
		} else {
			break
		}
	}
}

// getMove prints out a list of legal moves and allows the current player to
// select from one of them, returns the selected legal move.
//
// BUG: Scans until number is found in input, should scan once and clear input.
// BUG: Error message prints multiple times when non-integer inputs are given.
func getMove(legalMoves []Vec) Vec {
	fmt.Println("Your legal moves:")
	for i, move := range legalMoves {
		fmt.Printf("%d. %v\n", i+1, move)
	}
	var i int
	fmt.Print("Enter a number to select your move >> ")
	for {
		_, err := fmt.Scan(&i)
		if err == nil && i > 0 && i <= len(legalMoves) {
			break
		}
		fmt.Print("That is not a valid number, try again >> ")
	}
	return legalMoves[i-1]
}

// calcWinner counts the tiles of each colour on the board and determines who
// has the most pieces. Returns Empty if scores are the same, otherwise returns
// the color of the winner.
func calcWinner(board *[8][8]Tile) Tile {
	var darkScore, lightScore int
	for _, row := range board {
		for _, tile := range row {
			if tile == Dark {
				darkScore++
			}
			if tile == Light {
				lightScore++
			}
		}
	}
	if darkScore > lightScore {
		return Dark
	}
	if lightScore > darkScore {
		return Light
	}
	return Empty
}

func boardToString(board *[8][8]Tile) string {
	var boardString string = " 01234567\n"
	for j := 0; j < len(board[0]); j++ {
		boardString += fmt.Sprint(j)
		for i := 0; i < len(board); i++ {
			switch board[i][j] {
			case Empty:
				boardString += "┼"
			case Dark:
				boardString += "○"
			case Light:
				boardString += "●"
			case Legal:
				boardString += "?"
			}
		}
		boardString += "\n"
	}
	return boardString
}
