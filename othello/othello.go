package othello

type Tile int

type Move struct {
	x int
	y int
}

const (
	Empty Tile = iota
	Dark
	Light
	Legal
)

func GetStartingBoard() *[8][8]Tile {
	board := [8][8]Tile{}
	board[3][3] = Light
	board[3][4] = Dark
	board[4][3] = Dark
	board[4][4] = Light

	return &board
}

func PropagateLegalMoves(board *[8][8]Tile, legalMoves []Move) {
	for _, move := range legalMoves {
		board[move.x][move.y] = Legal
	}
}

func ClearLegalMoves(board *[8][8]Tile, legalMoves []Move) {
	for _, move := range legalMoves {
		// Check that Tile state is still legal to avoid clearing player move
		if board[move.x][move.y] == Legal {
			board[move.x][move.y] = Empty
		}
	}
}

// Need to document
func CalcLegalMoves(board *[8][8]Tile, currentPlayer Tile) []Move {
	var legalMoves []Move
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			// Moves can only be made on Empty tiles
			if board[x][y] != Empty {
				continue
			}
			pos := Move{x, y}
			for _, move := range circleTileCheck(pos, board, OppositePlayer(currentPlayer)) {
				dir := [2]int{move.x - x, move.y - y}
				if lineTileCheck(pos, dir, board, currentPlayer) {
					legalMoves = append(legalMoves, Move{x, y})
					break // Don't need to check legality multiple times
				}
			}
		}
	}
	return legalMoves
}

// circleTileCheck looks at the 8 tiles surrounding pos, then returns a slice
// of Moves containing coordinates of Tiles equal to target.
func circleTileCheck(pos Move, board *[8][8]Tile, target Tile) []Move {
	var matches []Move
	for i := max(0, pos.x-1); i < min(len(board), pos.x+2); i++ {
		for j := max(0, pos.y-1); j < min(len(board[i]), pos.y+2); j++ {
			if board[i][j] == target {
				matches = append(matches, Move{i, j})
			}
		}
	}
	return matches
}

/*
lineTileCheck steps from a position on the board (a move) in increments of dir.

Returns true if a continous line of placed tiles can be made from pos to a Tile
equal to target, otherwise returns false.
*/
func lineTileCheck(pos Move, dir [2]int, board *[8][8]Tile, target Tile) bool {
	// One step made outside of loop to check if same color tiles are adjacent
	step := Move{pos.x + dir[0], pos.y + dir[1]}
	if board[step.x][step.y] == target {
		return false
	}
	for {
		step = Move{step.x + dir[0], step.y + dir[1]}
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

// PlaceTile changes state of Tile of board at move to currentPlayer,
// and calls flipWalk to "flip" all opposite Tiles to the state of
// currentPlayer, in accordance with the game rules
// (https://wikipedia.org/wiki/Reversi#Rules).
func PlaceTile(move Move, board *[8][8]Tile, currentPlayer Tile) {
	oppositePlayer := OppositePlayer(currentPlayer)
	board[move.x][move.y] = currentPlayer
	for _, opp := range circleTileCheck(move, board, oppositePlayer) {
		dir := [2]int{opp.x - move.x, opp.y - move.y}
		if lineTileCheck(move, dir, board, currentPlayer) {
			flipWalk(move, dir, board, currentPlayer)
		}
	}
}

func flipWalk(move Move, dir [2]int, board *[8][8]Tile, currentPlayer Tile) {
	step := move
	for {
		step = Move{step.x + dir[0], step.y + dir[1]}
		if board[step.x][step.y] != currentPlayer {
			board[step.x][step.y] = currentPlayer
		} else {
			break
		}
	}
}

func EndOfTurn(board *[8][8]Tile, legalMoves []Move, currentPlayer Tile) ([]Move, Tile, bool) {
	ClearLegalMoves(board, legalMoves)

	nextPlayer := OppositePlayer(currentPlayer)
	newLegalMoves := CalcLegalMoves(board, nextPlayer)
	// If the current player has no legal moves, their turn is skipped
	if len(newLegalMoves) == 0 {
		nextPlayer = OppositePlayer(currentPlayer)
		newLegalMoves = CalcLegalMoves(board, nextPlayer)
		// If neither player has any legal moves, the game ends
		if len(newLegalMoves) == 0 {
			return []Move{}, Empty, false
		}
	}
	return newLegalMoves, nextPlayer, true
}

// CalcWinner counts the tiles of each colour on the board and determines who
// has the most pieces. Returns Empty if scores are the same, otherwise returns
// the color of the winner.
func CalcWinner(board *[8][8]Tile) Tile {
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

func BoardToString(board *[8][8]Tile) string {
	var boardString string = "    0  1  2  3  4  5  6  7\n  ┌──┬──┬──┬──┬──┬──┬──┬──┐\n"
	for j := 0; j < len(board[0]); j++ {
		boardString += string(j+48) + " " // 48 = ASCII offset of 0
		for i := 0; i < len(board); i++ {
			boardString += "│"
			switch board[i][j] {
			case Empty:
				boardString += "  "
			case Dark:
				boardString += "⚫"
			case Light:
				boardString += "⚪"
			case Legal:
				boardString += "??"
			}
		}
		boardString += "│\n"
		if j < len(board[0])-1 {
			boardString += "  ├──┼──┼──┼──┼──┼──┼──┼──┤\n"
		}
	}
	return boardString + "  └──┴──┴──┴──┴──┴──┴──┴──┘\n"
}

func OppositePlayer(player Tile) Tile {
	if player == Dark {
		return Light
	}
	return Dark
}

func PlayerString(player Tile) string {
	if player == Dark {
		return "Dark"
	}
	return "Light"
}
