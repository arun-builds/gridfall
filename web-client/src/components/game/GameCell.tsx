import { cn } from "@/lib/utils"
import {
  EMPTY_CELL,
  MISS_CELL,
  UNKNOWN_CELL,
  isEntityCell,
  isDestroyedCell,
  getEntityMeta,
  getEntityMetaByAbsId,
} from "@/lib/types"

interface GameCellProps {
  value: number
  x: number
  y: number
  isInteractive?: boolean
  isLastAttack?: boolean
  attackResult?: "hit" | "miss"
  showEntities?: boolean
  onClick?: (x: number, y: number) => void
}

export function GameCell({
  value,
  x,
  y,
  isInteractive = false,
  isLastAttack = false,
  attackResult,
  showEntities = false,
  onClick,
}: GameCellProps) {
  const handleClick = () => {
    if (isInteractive && onClick) {
      onClick(x, y)
    }
  }

  const getCellContent = () => {
    if (value === EMPTY_CELL) return null
    if (value === UNKNOWN_CELL) return null

    if (value === MISS_CELL) {
      return <span className="cell-miss-marker">•</span>
    }

    if (isEntityCell(value) && showEntities) {
      const meta = getEntityMeta(value)
      if (meta) {
        return <span className="cell-entity-icon">{meta.icon}</span>
      }
    }

    if (isDestroyedCell(value)) {
      const meta = getEntityMetaByAbsId(value)
      if (meta) {
        return <span className="cell-destroyed-icon">{meta.icon}</span>
      }
      return <span className="cell-hit-marker">✕</span>
    }

    return null
  }

  const getEntityGlow = () => {
    if (isEntityCell(value) && showEntities) {
      const meta = getEntityMeta(value)
      return meta?.glowColor ?? undefined
    }
    return undefined
  }

  const entityMeta =
    isEntityCell(value) && showEntities
      ? getEntityMeta(value)
      : isDestroyedCell(value)
        ? getEntityMetaByAbsId(value)
        : undefined

  return (
    <button
      type="button"
      className={cn(
        "game-cell",
        value === EMPTY_CELL && "cell-empty",
        value === UNKNOWN_CELL && "cell-unknown",
        value === MISS_CELL && "cell-miss",
        isEntityCell(value) && showEntities && "cell-entity",
        isDestroyedCell(value) && "cell-destroyed",
        isInteractive && "cell-interactive",
        isLastAttack && attackResult === "hit" && "cell-anim-hit",
        isLastAttack && attackResult === "miss" && "cell-anim-miss",
        entityMeta && `cell-${entityMeta.name.toLowerCase()}`
      )}
      style={
        getEntityGlow()
          ? ({ "--entity-glow": getEntityGlow() } as React.CSSProperties)
          : undefined
      }
      onClick={handleClick}
      disabled={!isInteractive}
      aria-label={`Cell ${String.fromCharCode(65 + x)}${y + 1}`}
    >
      {getCellContent()}
    </button>
  )
}
