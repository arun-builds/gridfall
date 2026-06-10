import { Button } from "@/components/ui/button"

interface GameOverScreenProps {
  winner: string
  myId: string
  onPlayAgain: () => void
}

export function GameOverScreen({ winner, myId, onPlayAgain }: GameOverScreenProps) {
  const isVictory = winner === myId

  return (
    <div className="gameover-screen">
      {/* Particle effects background */}
      <div className="gameover-particles" aria-hidden="true">
        {Array.from({ length: 30 }, (_, i) => (
          <div
            key={i}
            className={`gameover-particle ${isVictory ? "particle-victory" : "particle-defeat"}`}
            style={{
              left: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 2}s`,
              animationDuration: `${2 + Math.random() * 3}s`,
            }}
          />
        ))}
      </div>

      <div className="gameover-content">
        <div className={`gameover-icon ${isVictory ? "icon-victory" : "icon-defeat"}`}>
          {isVictory ? "👑" : "💀"}
        </div>

        <h1 className={`gameover-title ${isVictory ? "title-victory" : "title-defeat"}`}>
          {isVictory ? "VICTORY" : "DEFEAT"}
        </h1>

        <p className="gameover-subtitle">
          {isVictory
            ? "You destroyed all enemy forces!"
            : "Your forces have been eliminated."}
        </p>

        <Button className="gameover-btn" onClick={onPlayAgain}>
          Play Again
        </Button>
      </div>
    </div>
  )
}
