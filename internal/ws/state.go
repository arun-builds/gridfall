package ws

type MatchPhase string

const (
	PhaseWaiting   MatchPhase = "waiting"
	PhasePlacement MatchPhase = "placement"
	PhaseBattle    MatchPhase = "battle"
	PhaseGameOver  MatchPhase = "game_over"
)

type MatchState struct {
	Phase       MatchPhase
	CurrentTurn string

	Board1 [][]int
	Board2 [][]int

	Player1Ready bool
	Player2Ready bool
}
