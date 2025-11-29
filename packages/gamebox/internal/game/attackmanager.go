package game

type AttackManager struct {
	// attack is map of playerId -> map of playerIds that they are currently attacking -> troopCount in attack
	attacks map[string]map[string]int32

	board *Board
}

func NewAttackManager(board *Board) *AttackManager {
	return &AttackManager{
		board:   board,
		attacks: make(map[string]map[string]int32),
	}
}

func (am *AttackManager) InitAttack(playerId string, tileId int32, troopCount int32) {
}
