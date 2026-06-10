import type {
  GameState,
  AttackResultEvent,
  OpponentAttackedEvent,
  BattleStartedEvent,
  GameOverEvent,
  GameStateEvent,
} from "@/lib/types"
import { createEmptyBoard, createUnknownBoard } from "@/lib/types"

// ── Actions ────────────────────────────────────────────────────────────
export type GameAction =
  | { type: "CONNECTING" }
  | { type: "CONNECTED"; myId: string }
  | { type: "DISCONNECTED" }
  | { type: "SERVER_GAME_STATE"; event: GameStateEvent }
  | { type: "SERVER_PLACEMENT_SUCCESS" }
  | { type: "SERVER_OPPONENT_READY" }
  | { type: "SERVER_BATTLE_STARTED"; event: BattleStartedEvent }
  | { type: "SERVER_ATTACK_RESULT"; event: AttackResultEvent }
  | { type: "SERVER_OPPONENT_ATTACKED"; event: OpponentAttackedEvent }
  | { type: "SERVER_GAME_OVER"; event: GameOverEvent }
  | { type: "SERVER_ERROR"; message: string }
  | { type: "CLEAR_LAST_ATTACK" }
  | { type: "RESET" }

// ── Reducer ────────────────────────────────────────────────────────────
export function gameReducer(state: GameState, action: GameAction): GameState {
  switch (action.type) {
    case "CONNECTING":
      return {
        ...state,
        connectionStatus: "connecting",
        error: "",
      }

    case "CONNECTED":
      return {
        ...state,
        connectionStatus: "connected",
        myId: action.myId,
        error: "",
      }

    case "DISCONNECTED":
      return {
        ...state,
        connectionStatus: "disconnected",
      }

    case "SERVER_GAME_STATE":
      return {
        ...state,
        myId: action.event.your_id || state.myId,
        phase: action.event.phase,
        currentTurn: action.event.current_turn,
        yourBoard: action.event.your_board,
        opponentView: action.event.opponent_view,
      }

    case "SERVER_PLACEMENT_SUCCESS":
      return {
        ...state,
      }

    case "SERVER_OPPONENT_READY":
      return {
        ...state,
      }

    case "SERVER_BATTLE_STARTED":
      return {
        ...state,
        phase: "battle",
        currentTurn: action.event.current_turn,
      }

    case "SERVER_ATTACK_RESULT": {
      const { x, y, result } = action.event
      const newOpponentView = state.opponentView.map((row) => [...row])

      if (result === "miss") {
        newOpponentView[y][x] = -99 // MISS_CELL
      } else if (result === "hit") {
        // Map destroyed entity name to negative value
        const entityMap: Record<string, number> = {
          scout: -1,
          battleship: -2,
          mage: -3,
          assassin: -4,
        }
        const entity = action.event.destroyed_entity ?? ""
        newOpponentView[y][x] = entityMap[entity] ?? -1
      }

      return {
        ...state,
        opponentView: newOpponentView,
        currentTurn: "", // Turn switches, will be set by next state update
        lastAttack: { board: "opponent", x, y, result },
      }
    }

    case "SERVER_OPPONENT_ATTACKED": {
      const { x, y, result } = action.event
      const newYourBoard = state.yourBoard.map((row) => [...row])

      if (result === "miss") {
        newYourBoard[y][x] = -99 // MISS_CELL
      } else if (result === "hit") {
        // Negate the entity value to mark as destroyed
        const currentValue = newYourBoard[y][x]
        if (currentValue > 0) {
          newYourBoard[y][x] = -currentValue
        }
      }

      return {
        ...state,
        yourBoard: newYourBoard,
        currentTurn: state.myId, // It's now our turn
        lastAttack: { board: "yours", x, y, result },
      }
    }

    case "SERVER_GAME_OVER":
      return {
        ...state,
        phase: "game_over",
        winner: action.event.winner,
      }

    case "SERVER_ERROR":
      return {
        ...state,
        error: action.message,
      }

    case "CLEAR_LAST_ATTACK":
      return {
        ...state,
        lastAttack: null,
      }

    case "RESET":
      return {
        ...state,
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

    default:
      return state
  }
}
