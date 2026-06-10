import { useCallback } from "react"
import { GameBoard } from "./GameBoard"
import {
  ENTITIES,
  UNKNOWN_CELL,
  isDestroyedCell,
  type GameState,
} from "@/lib/types"

interface BattleScreenProps {
  state: GameState
  onAttack: (x: number, y: number) => void
}

export function BattleScreen({ state, onAttack }: BattleScreenProps) {
  const isMyTurn = state.currentTurn === state.myId

  const handleOpponentCellClick = useCallback(
    (x: number, y: number) => {
      if (!isMyTurn) return
      // Only allow clicking unknown cells
      const cellValue = state.opponentView[y]?.[x]
      if (cellValue !== UNKNOWN_CELL) return
      onAttack(x, y)
    },
    [isMyTurn, state.opponentView, onAttack]
  )

  // Determine which entities are destroyed on each side
  const yourDestroyedIds = new Set<number>()
  for (const row of state.yourBoard) {
    for (const cell of row) {
      if (isDestroyedCell(cell)) {
        yourDestroyedIds.add(Math.abs(cell))
      }
    }
  }

  const opponentDestroyedIds = new Set<number>()
  for (const row of state.opponentView) {
    for (const cell of row) {
      if (isDestroyedCell(cell)) {
        opponentDestroyedIds.add(Math.abs(cell))
      }
    }
  }

  const yourBoardAttack =
    state.lastAttack?.board === "yours" ? state.lastAttack : null
  const opponentBoardAttack =
    state.lastAttack?.board === "opponent" ? state.lastAttack : null

  return (
    <div className="battle-screen">
      {/* Turn indicator */}
      <div className={`battle-turn-bar ${isMyTurn ? "your-turn" : "enemy-turn"}`}>
        <div className="battle-turn-indicator" />
        <span className="battle-turn-text">
          {isMyTurn ? "Your Turn — Select a target" : "Opponent's Turn — Stand by…"}
        </span>
      </div>

      <div className="battle-body">
        {/* Your board */}
        <div className="battle-board-section">
          <GameBoard
            board={state.yourBoard}
            showEntities={true}
            isInteractive={false}
            lastAttack={yourBoardAttack}
            label="Your Grid"
          />

          {/* Your entity status */}
          <div className="battle-entity-status">
            {ENTITIES.map((entity) => (
              <div
                key={entity.id}
                className={`battle-entity-pip ${
                  yourDestroyedIds.has(entity.id) ? "pip-destroyed" : "pip-alive"
                } pip-${entity.name.toLowerCase()}`}
              >
                <span className="pip-icon">{entity.icon}</span>
                <span className="pip-name">{entity.name}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Divider */}
        <div className="battle-divider">
          <span className="battle-vs">VS</span>
        </div>

        {/* Opponent board */}
        <div className="battle-board-section">
          <GameBoard
            board={state.opponentView}
            showEntities={false}
            isInteractive={isMyTurn}
            lastAttack={opponentBoardAttack}
            onCellClick={handleOpponentCellClick}
            label="Enemy Grid"
          />

          {/* Opponent entity status */}
          <div className="battle-entity-status">
            {ENTITIES.map((entity) => (
              <div
                key={entity.id}
                className={`battle-entity-pip ${
                  opponentDestroyedIds.has(entity.id) ? "pip-destroyed" : "pip-alive"
                } pip-${entity.name.toLowerCase()}`}
              >
                <span className="pip-icon">{entity.icon}</span>
                <span className="pip-name">{entity.name}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
