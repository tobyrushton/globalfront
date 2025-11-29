package game

import (
	"math"
	"sync"

	v1 "github.com/tobyrushton/globalfront/pb/game/v1"
)

type Attack struct {
	troopCount int32
	border     []int32
}

type AttackManager struct {
	// attack is map of playerId -> map of playerIds that they are currently attacking -> troopCount in attack
	attacks   map[string]map[string]*Attack
	attacksMu sync.Mutex

	board   *Board
	players *map[string]*v1.Player

	// attacking exponential constant
	k float64
}

func NewAttackManager(board *Board, players *map[string]*v1.Player) *AttackManager {
	return &AttackManager{
		board:   board,
		attacks: make(map[string]map[string]*Attack),
		players: players,
		k:       1.5,
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

func (am *AttackManager) CalculateAttacks() {
	am.attacksMu.Lock()
	defer am.attacksMu.Unlock()

	for playerFrom, attacks := range am.attacks {
		for playerTo, attack := range attacks {
			defendingTroops := int32(0)
			// as wilderness won't have any troops
			if defender, exists := (*am.players)[playerTo]; exists {
				defendingTroops = defender.TroopCount
			}

			captured, remainingAttackers, _ := am.calculateAdvance(defendingTroops, attack.troopCount, int32(len(attack.border)))
			// update troop counts
			if remainingAttackers == 0 {
				delete(am.attacks[playerFrom], playerTo)
			} else {
				attack.troopCount = remainingAttackers
			}
			am.board.AdvancePlayer(attack.border, playerFrom, playerTo, captured)
		}
	}
}

func (am *AttackManager) calculateAdvance(defendingTroops, attackingTroops, borderLength int32) (int32, int32, int32) {
	a := math.Pow(float64(attackingTroops), am.k)
	d := math.Pow(float64(defendingTroops), am.k)
	ratio := a / (a + d)

	captured := int32(ratio * float64(borderLength))

	if captured > borderLength {
		captured = borderLength
	}

	frac := float64(captured) / float64(borderLength)

	remainingAttackers := int32(float64(attackingTroops) * (1 - frac))
	remainingDefenders := int32(float64(defendingTroops) * (1 - frac))

	return captured, remainingAttackers, remainingDefenders
}
