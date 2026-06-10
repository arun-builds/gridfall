import { useState, useCallback } from "react"
import { Button } from "@/components/ui/button"
import { GameBoard } from "./GameBoard"
import { EntityCard } from "./EntityCard"
import {
  ENTITIES,
  BOARD_SIZE,
  createEmptyBoard,
  type EntityId,
  type PlacementPayload,
} from "@/lib/types"

interface PlacementScreenProps {
  onPlaceEntities: (placements: PlacementPayload[]) => void
}

interface PlacedEntity {
  entity: EntityId
  x: number
  y: number
}

export function PlacementScreen({ onPlaceEntities }: PlacementScreenProps) {
  const [selectedEntity, setSelectedEntity] = useState<EntityId | null>(null)
  const [placements, setPlacements] = useState<PlacedEntity[]>([])
  const [submitted, setSubmitted] = useState(false)

  const placedEntityIds = new Set(placements.map((p) => p.entity))

  // Build the display board from placements
  const displayBoard = createEmptyBoard()
  for (const p of placements) {
    displayBoard[p.y][p.x] = p.entity
  }

  const handleCellClick = useCallback(
    (x: number, y: number) => {
      if (submitted) return
      if (x < 0 || x >= BOARD_SIZE || y < 0 || y >= BOARD_SIZE) return

      // If clicking a cell that already has a placed entity, remove it
      const existingIndex = placements.findIndex((p) => p.x === x && p.y === y)
      if (existingIndex >= 0) {
        setPlacements((prev) => prev.filter((_, i) => i !== existingIndex))
        return
      }

      if (!selectedEntity) return
      if (placedEntityIds.has(selectedEntity)) return

      setPlacements((prev) => [...prev, { entity: selectedEntity, x, y }])
      setSelectedEntity(null)
    },
    [selectedEntity, placements, placedEntityIds, submitted]
  )

  const handleEntitySelect = (entityId: EntityId) => {
    if (submitted) return
    if (placedEntityIds.has(entityId)) return
    setSelectedEntity((prev) => (prev === entityId ? null : entityId))
  }

  const handleSubmit = () => {
    if (placements.length !== ENTITIES.length) return
    setSubmitted(true)
    onPlaceEntities(
      placements.map((p) => ({
        entity: p.entity,
        x: p.x,
        y: p.y,
      }))
    )
  }

  const handleReset = () => {
    if (submitted) return
    setPlacements([])
    setSelectedEntity(null)
  }

  const allPlaced = placements.length === ENTITIES.length

  return (
    <div className="placement-screen">
      <div className="placement-header">
        <h1 className="placement-title">Deploy Your Forces</h1>
        <p className="placement-subtitle">
          {submitted
            ? "Waiting for opponent to deploy…"
            : selectedEntity
              ? `Click a cell to place your ${ENTITIES.find((e) => e.id === selectedEntity)?.name}`
              : "Select a unit, then click a cell to place it"}
        </p>
      </div>

      <div className="placement-body">
        {/* Board */}
        <div className="placement-board-area">
          <GameBoard
            board={displayBoard}
            showEntities={true}
            isInteractive={!submitted}
            onCellClick={handleCellClick}
            label="Your Grid"
          />
        </div>

        {/* Entity Sidebar */}
        <div className="placement-sidebar">
          <h3 className="placement-sidebar-title">Units</h3>
          <div className="placement-entity-list">
            {ENTITIES.map((entity) => (
              <EntityCard
                key={entity.id}
                entity={entity}
                isSelected={selectedEntity === entity.id}
                isPlaced={placedEntityIds.has(entity.id)}
                onClick={() => handleEntitySelect(entity.id)}
              />
            ))}
          </div>

          <div className="placement-actions">
            <Button
              className="placement-ready-btn"
              onClick={handleSubmit}
              disabled={!allPlaced || submitted}
            >
              {submitted ? (
                <>
                  <span className="placement-ready-spinner" />
                  Waiting…
                </>
              ) : (
                "Ready"
              )}
            </Button>
            <Button
              variant="outline"
              className="placement-reset-btn"
              onClick={handleReset}
              disabled={submitted || placements.length === 0}
            >
              Reset
            </Button>
          </div>
        </div>
      </div>

      {/* Submitted overlay spinner */}
      {submitted && (
        <div className="placement-waiting-overlay">
          <div className="placement-waiting-pulse" />
        </div>
      )}
    </div>
  )
}
