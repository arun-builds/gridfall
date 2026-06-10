import { useReducer, useRef, useCallback, useEffect } from "react"
import { gameReducer } from "./gameReducer"
import { initialGameState } from "@/lib/types"
import type { PlacementPayload, ServerEvent } from "@/lib/types"

const WS_BASE = "ws://localhost:8080/ws"

export function useGameSocket() {
  const [state, dispatch] = useReducer(gameReducer, initialGameState)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const roomRef = useRef<string>("")

  const cleanup = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current)
      reconnectTimerRef.current = null
    }
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
  }, [])

  const handleMessage = useCallback(
    (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data) as ServerEvent

        switch (data.type) {
          case "game_state":
            dispatch({
              type: "SERVER_GAME_STATE",
              event: data as never,
            })
            break

          case "placement_success":
            dispatch({ type: "SERVER_PLACEMENT_SUCCESS" })
            break

          case "opponent_ready":
            dispatch({ type: "SERVER_OPPONENT_READY" })
            break

          case "battle_started":
            dispatch({
              type: "SERVER_BATTLE_STARTED",
              event: data as never,
            })
            // Fetch full board state now that battle has begun
            wsRef.current?.send(JSON.stringify({ type: "get_state" }))
            break

          case "attack_result":
            dispatch({
              type: "SERVER_ATTACK_RESULT",
              event: data as never,
            })
            break

          case "opponent_attacked":
            dispatch({
              type: "SERVER_OPPONENT_ATTACKED",
              event: data as never,
            })
            break

          case "game_over":
            dispatch({
              type: "SERVER_GAME_OVER",
              event: data as never,
            })
            break

          case "error":
            dispatch({
              type: "SERVER_ERROR",
              message: (data as { message: string }).message,
            })
            break

          default:
            console.warn("Unknown server event:", data.type)
        }
      } catch (err) {
        console.error("Failed to parse server message:", err)
      }
    },
    []
  )

  const connect = useCallback(
    (roomId: string) => {
      cleanup()
      roomRef.current = roomId
      dispatch({ type: "CONNECTING" })

      const ws = new WebSocket(`${WS_BASE}?room=${encodeURIComponent(roomId)}`)
      wsRef.current = ws

      ws.onopen = () => {
        // We don't know our ID until the server tells us via game_state,
        // but we mark as connected. The ID is derived from the first
        // game_state response's current_turn or attack_result context.
        dispatch({ type: "CONNECTED", myId: "" })

        // Request initial state
        ws.send(JSON.stringify({ type: "get_state" }))
      }

      ws.onmessage = handleMessage

      ws.onclose = () => {
        dispatch({ type: "DISCONNECTED" })
      }

      ws.onerror = (err) => {
        console.error("WebSocket error:", err)
      }
    },
    [cleanup, handleMessage]
  )

  const disconnect = useCallback(() => {
    cleanup()
    dispatch({ type: "RESET" })
  }, [cleanup])

  const sendPlacement = useCallback((placements: PlacementPayload[]) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: "place_entities",
          payload: { placements },
        })
      )
    }
  }, [])

  const sendAttack = useCallback((x: number, y: number) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: "attack",
          payload: { x, y },
        })
      )
    }
  }, [])

  const requestState = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: "get_state" }))
    }
  }, [])

  // Clear last attack animation after a delay
  useEffect(() => {
    if (state.lastAttack) {
      const timer = setTimeout(() => {
        dispatch({ type: "CLEAR_LAST_ATTACK" })
      }, 1200)
      return () => clearTimeout(timer)
    }
  }, [state.lastAttack])

  // Poll for state while waiting for opponent to join.
  // The server doesn't push a notification to player 1 when
  // player 2 joins, so we poll to detect the phase transition.
  useEffect(() => {
    if (
      state.connectionStatus === "connected" &&
      state.phase === "waiting" &&
      wsRef.current?.readyState === WebSocket.OPEN
    ) {
      const interval = setInterval(() => {
        wsRef.current?.send(JSON.stringify({ type: "get_state" }))
      }, 2000)
      return () => clearInterval(interval)
    }
  }, [state.connectionStatus, state.phase])

  // Cleanup on unmount
  useEffect(() => {
    return cleanup
  }, [cleanup])

  return {
    state,
    connect,
    disconnect,
    sendPlacement,
    sendAttack,
    requestState,
  }
}
