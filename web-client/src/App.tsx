import { useGameSocket } from "@/hooks/useGameSocket"
import { LobbyScreen } from "@/components/game/LobbyScreen"
import { PlacementScreen } from "@/components/game/PlacementScreen"
import { BattleScreen } from "@/components/game/BattleScreen"
import { GameOverScreen } from "@/components/game/GameOverScreen"

export function App() {
  const { state, connect, disconnect, sendPlacement, sendAttack } =
    useGameSocket()

  // Lobby: not yet connected or waiting for opponent
  if (
    state.connectionStatus === "idle" ||
    state.connectionStatus === "connecting" ||
    (state.connectionStatus === "connected" && state.phase === "waiting")
  ) {
    return (
      <LobbyScreen
        connectionStatus={state.connectionStatus}
        onJoinRoom={(roomId) => connect(roomId)}
      />
    )
  }

  // Disconnected mid-game
  if (state.connectionStatus === "disconnected") {
    return (
      <div className="disconnect-screen">
        <div className="disconnect-content">
          <div className="disconnect-icon">⚠</div>
          <h1 className="disconnect-title">Connection Lost</h1>
          <p className="disconnect-subtitle">
            You were disconnected from the server.
          </p>
          <button
            className="disconnect-btn"
            onClick={() => disconnect()}
          >
            Return to Lobby
          </button>
        </div>
      </div>
    )
  }

  // Placement phase
  if (state.phase === "placement") {
    return <PlacementScreen onPlaceEntities={sendPlacement} />
  }

  // Battle phase
  if (state.phase === "battle") {
    return <BattleScreen state={state} onAttack={sendAttack} />
  }

  // Game over
  if (state.phase === "game_over") {
    return (
      <GameOverScreen
        winner={state.winner}
        myId={state.myId}
        onPlayAgain={() => disconnect()}
      />
    )
  }

  return null
}

export default App
