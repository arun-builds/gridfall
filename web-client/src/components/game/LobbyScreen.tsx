import { useState } from "react"
import { Button } from "@/components/ui/button"

interface LobbyScreenProps {
  connectionStatus: string
  onJoinRoom: (roomId: string) => void
}

function generateRoomId() {
  const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
  let result = ""
  for (let i = 0; i < 6; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

export function LobbyScreen({ connectionStatus, onJoinRoom }: LobbyScreenProps) {
  const [roomId, setRoomId] = useState(() => generateRoomId())
  const isConnecting = connectionStatus === "connecting"
  const isWaiting = connectionStatus === "connected"

  return (
    <div className="lobby-screen">
      {/* Animated grid background */}
      <div className="lobby-grid-bg" aria-hidden="true">
        {Array.from({ length: 64 }, (_, i) => (
          <div
            key={i}
            className="lobby-grid-cell"
            style={{ animationDelay: `${Math.random() * 3}s` }}
          />
        ))}
      </div>

      <div className="lobby-content">
        {/* Logo / Title */}
        <div className="lobby-hero">
          <h1 className="lobby-title">
            <span className="lobby-title-grid">GRID</span>
            <span className="lobby-title-fall">FALL</span>
          </h1>
          <p className="lobby-subtitle">1v1 Tactical Grid Combat</p>
        </div>

        {/* Room Join Panel */}
        <div className="lobby-panel">
          {isWaiting ? (
            <div className="lobby-waiting">
              <div className="lobby-waiting-spinner" />
              <p className="lobby-waiting-text">Waiting for opponent…</p>
              <p className="lobby-waiting-room">
                Room: <span className="lobby-room-code">{roomId}</span>
              </p>
              <p className="lobby-waiting-hint">Share this code with a friend</p>
            </div>
          ) : (
            <>
              <label className="lobby-label" htmlFor="room-input">
                Room Code
              </label>
              <div className="lobby-input-row">
                <input
                  id="room-input"
                  type="text"
                  className="lobby-input"
                  value={roomId}
                  onChange={(e) => setRoomId(e.target.value.toUpperCase())}
                  placeholder="Enter room code"
                  maxLength={12}
                  disabled={isConnecting}
                />
                <Button
                  className="lobby-btn"
                  onClick={() => onJoinRoom(roomId)}
                  disabled={!roomId.trim() || isConnecting}
                >
                  {isConnecting ? "Connecting…" : "Join Room"}
                </Button>
              </div>
              <button
                type="button"
                className="lobby-generate-btn"
                onClick={() => setRoomId(generateRoomId())}
                disabled={isConnecting}
              >
                Generate new code
              </button>
            </>
          )}
        </div>

        {/* Instructions */}
        <div className="lobby-instructions">
          <div className="lobby-instruction-item">
            <span className="lobby-instruction-num">1</span>
            <span>Create or join a room</span>
          </div>
          <div className="lobby-instruction-item">
            <span className="lobby-instruction-num">2</span>
            <span>Place your 4 units on the grid</span>
          </div>
          <div className="lobby-instruction-item">
            <span className="lobby-instruction-num">3</span>
            <span>Destroy all enemy units to win</span>
          </div>
        </div>
      </div>
    </div>
  )
}
