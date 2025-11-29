package game

import "sync"

type Attack struct {
	troopCount int32
	border     []int32
}

type AttackManager struct {
	// attack is map of playerId -> map of playerIds that they are currently attacking -> troopCount in attack
	attacks   map[string]map[string]*Attack
	attacksMu sync.Mutex

	board *Board
}

func NewAttackManager(board *Board) *AttackManager {
	return &AttackManager{
		board:   board,
		attacks: make(map[string]map[string]*Attack),
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
	border := am.board.FindBorder(tile.PlayerId(), playerId, tileId)
	if len(border) == 0 {
		return
	}

	am.attacksMu.Lock()
	defer am.attacksMu.Unlock()

	if attacks, exists := am.attacks[tile.PlayerId()]; exists {
		if attack, exists := attacks[playerId]; exists {
			if attack.troopCount > troopCount {
				attack.troopCount -= troopCount
			} else {
				delete(am.attacks[tile.PlayerId()], playerId)
				// in the case of 0 we don't want to add more troops
				if troopCount > attack.troopCount {
					am.startAttack(playerId, tile.PlayerId(), troopCount-attack.troopCount, border)
				}
			}
		}
	} else {
		am.startAttack(playerId, tile.PlayerId(), troopCount, border)
	}
}

func (am *AttackManager) startAttack(playerFrom, playerTo string, troopCount int32, border []int32) {
	if _, exists := am.attacks[playerFrom]; !exists {
		am.attacks[playerFrom] = make(map[string]*Attack)
	}
	if _, exists := am.attacks[playerFrom][playerTo]; !exists {
		am.attacks[playerFrom][playerTo] = &Attack{
			troopCount: 0,
			border:     border,
		}
	}
	am.attacks[playerFrom][playerTo].troopCount += troopCount
}
