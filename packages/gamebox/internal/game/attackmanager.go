package game

import "sync"

type AttackManager struct {
	// attack is map of playerId -> map of playerIds that they are currently attacking -> troopCount in attack
	attacks   map[string]map[string]int32
	attacksMu sync.Mutex

	board *Board
}

func NewAttackManager(board *Board) *AttackManager {
	return &AttackManager{
		board:   board,
		attacks: make(map[string]map[string]int32),
	}
}

func (am *AttackManager) InitAttack(playerId string, tileId int32, troopCount int32) {
	tile := am.board.GetTile(tileId)
	if tile == nil {
		return
	}
	if tile.PlayerId() == playerId {
		return
	}

	am.attacksMu.Lock()
	defer am.attacksMu.Unlock()

	if attacks, exists := am.attacks[tile.PlayerId()]; exists {
		if attack, exists := attacks[playerId]; exists {
			if attack > troopCount {
				am.attacks[tile.PlayerId()][playerId] -= troopCount
			} else {
				delete(am.attacks[tile.PlayerId()], playerId)
				// in the case of 0 we don't want to add more troops
				if troopCount > attack {
					am.startAttack(playerId, tile.PlayerId(), troopCount-attack)
				}
			}
		}
	} else {
		am.startAttack(playerId, tile.PlayerId(), troopCount)
	}
}

func (am *AttackManager) startAttack(playerFrom, playerTo string, troopCount int32) {
	if _, exists := am.attacks[playerFrom]; !exists {
		am.attacks[playerFrom] = make(map[string]int32)
	}
	am.attacks[playerFrom][playerTo] += troopCount
}
