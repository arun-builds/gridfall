// ── Cell values (mirrors internal/ws/entities.go) ──────────────────────
export const EMPTY_CELL = 0
export const MISS_CELL = -99
export const UNKNOWN_CELL = -100

export const SCOUT = 1
export const BATTLESHIP = 2
export const MAGE = 3
export const ASSASSIN = 4

export const DESTROYED_SCOUT = -SCOUT
export const DESTROYED_BATTLESHIP = -BATTLESHIP
export const DESTROYED_MAGE = -MAGE
export const DESTROYED_ASSASSIN = -ASSASSIN

export const BOARD_SIZE = 8

// ── Entity metadata ────────────────────────────────────────────────────
export type EntityId = typeof SCOUT | typeof BATTLESHIP | typeof MAGE | typeof ASSASSIN

export interface EntityMeta {
  id: EntityId
  name: string
  description: string
  colorClass: string
  glowColor: string
  icon: string
}

export const ENTITIES: EntityMeta[] = [
  {
    id: SCOUT,
    name: "Scout",
    description: "Swift reconnaissance unit",
    colorClass: "text-cyan-400",
    glowColor: "rgba(34,211,238,0.6)",
    icon: "⚡",
  },
  {
    id: BATTLESHIP,
    name: "Battleship",
    description: "Heavy armored warship",
    colorClass: "text-amber-400",
    glowColor: "rgba(251,191,36,0.6)",
    icon: "⚓",
  },
  {
    id: MAGE,
    name: "Mage",
    description: "Arcane spellcaster",
    colorClass: "text-violet-400",
    glowColor: "rgba(167,139,250,0.6)",
    icon: "✦",
  },
  {
    id: ASSASSIN,
    name: "Assassin",
    description: "Deadly stealth operative",
    colorClass: "text-rose-400",
    glowColor: "rgba(251,113,133,0.6)",
    icon: "🗡",
  },
]

export function getEntityMeta(id: number): EntityMeta | undefined {
  return ENTITIES.find((e) => e.id === id)
}

export function getEntityMetaByAbsId(cellValue: number): EntityMeta | undefined {
  return ENTITIES.find((e) => e.id === Math.abs(cellValue))
}

export function isEntityCell(cell: number): boolean {
  return cell >= 1 && cell <= 4
}

export function isDestroyedCell(cell: number): boolean {
  return cell >= -4 && cell <= -1
}

// ── Match phase ────────────────────────────────────────────────────────
export type MatchPhase = "waiting" | "placement" | "battle" | "game_over"

// ── WebSocket event types ──────────────────────────────────────────────

// Events we send
export interface PlacementPayload {
  entity: number
  x: number
  y: number
}

export interface AttackPayload {
  x: number
  y: number
}

// Events we receive
export interface ServerEvent {
  type: string
  [key: string]: unknown
}

export interface GameStateEvent {
  type: "game_state"
  your_id: string
  phase: MatchPhase
  current_turn: string
  your_board: number[][]
  opponent_view: number[][]
}

export interface AttackResultEvent {
  type: "attack_result"
  result: "hit" | "miss"
  x: number
  y: number
  destroyed_entity?: string
}

export interface OpponentAttackedEvent {
  type: "opponent_attacked"
  result: "hit" | "miss"
  x: number
  y: number
}

export interface BattleStartedEvent {
  type: "battle_started"
  current_turn: string
}

export interface GameOverEvent {
  type: "game_over"
  winner: string
}

export interface ErrorEvent {
  type: "error"
  message: string
}

// ── Game state ─────────────────────────────────────────────────────────
export type ConnectionStatus = "idle" | "connecting" | "connected" | "disconnected"

export interface GameState {
  connectionStatus: ConnectionStatus
  phase: MatchPhase
  myId: string
  currentTurn: string
  yourBoard: number[][]
  opponentView: number[][]
  winner: string
  lastAttack: {
    board: "yours" | "opponent"
    x: number
    y: number
    result: "hit" | "miss"
  } | null
  error: string
}

export function createEmptyBoard(): number[][] {
  return Array.from({ length: BOARD_SIZE }, () => Array(BOARD_SIZE).fill(0))
}

export function createUnknownBoard(): number[][] {
  return Array.from({ length: BOARD_SIZE }, () => Array(BOARD_SIZE).fill(UNKNOWN_CELL))
}

export const initialGameState: GameState = {
  connectionStatus: "idle",
  phase: "waiting",
  myId: "",
  currentTurn: "",
  yourBoard: createEmptyBoard(),
  opponentView: createUnknownBoard(),
  winner: "",
  lastAttack: null,
  error: "",
}
