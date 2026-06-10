import { BOARD_SIZE } from "@/lib/types"
import { GameCell } from "./GameCell"

interface GameBoardProps {
  board: number[][]
  showEntities?: boolean
  isInteractive?: boolean
  lastAttack?: { x: number; y: number; result: "hit" | "miss" } | null
  onCellClick?: (x: number, y: number) => void
  label?: string
}

const COL_LABELS = "ABCDEFGH".split("")

export function GameBoard({
  board,
  showEntities = false,
  isInteractive = false,
  lastAttack,
  onCellClick,
  label,
}: GameBoardProps) {
  return (
    <div className="game-board-wrapper">
      {label && <h2 className="game-board-label">{label}</h2>}

      <div className="game-board-container">
        {/* Column labels */}
        <div className="game-board-col-labels">
          <div className="game-board-corner" />
          {COL_LABELS.map((col) => (
            <div key={col} className="game-board-col-label">
              {col}
            </div>
          ))}
        </div>

        {/* Rows */}
        {Array.from({ length: BOARD_SIZE }, (_, y) => (
          <div key={y} className="game-board-row">
            <div className="game-board-row-label">{y + 1}</div>
            {Array.from({ length: BOARD_SIZE }, (_, x) => {
              const isLastAttack =
                lastAttack !== null &&
                lastAttack !== undefined &&
                lastAttack.x === x &&
                lastAttack.y === y

              return (
                <GameCell
                  key={`${x}-${y}`}
                  value={board[y]?.[x] ?? 0}
                  x={x}
                  y={y}
                  isInteractive={isInteractive}
                  isLastAttack={isLastAttack}
                  attackResult={isLastAttack ? lastAttack.result : undefined}
                  showEntities={showEntities}
                  onClick={onCellClick}
                />
              )
            })}
          </div>
        ))}
      </div>
    </div>
  )
}
