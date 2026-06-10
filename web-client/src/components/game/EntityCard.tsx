import { cn } from "@/lib/utils"
import type { EntityMeta } from "@/lib/types"

interface EntityCardProps {
  entity: EntityMeta
  isSelected?: boolean
  isPlaced?: boolean
  isDestroyed?: boolean
  onClick?: () => void
}

export function EntityCard({
  entity,
  isSelected = false,
  isPlaced = false,
  isDestroyed = false,
  onClick,
}: EntityCardProps) {
  return (
    <button
      type="button"
      className={cn(
        "entity-card",
        isSelected && "entity-card-selected",
        isPlaced && "entity-card-placed",
        isDestroyed && "entity-card-destroyed",
        `entity-card-${entity.name.toLowerCase()}`
      )}
      style={{ "--entity-glow": entity.glowColor } as React.CSSProperties}
      onClick={onClick}
      disabled={isPlaced || isDestroyed}
    >
      <span className="entity-card-icon">{entity.icon}</span>
      <div className="entity-card-info">
        <span className="entity-card-name">{entity.name}</span>
        <span className="entity-card-desc">{entity.description}</span>
      </div>
      {isPlaced && <span className="entity-card-badge">Placed</span>}
      {isDestroyed && <span className="entity-card-badge entity-card-badge-dead">Lost</span>}
    </button>
  )
}
