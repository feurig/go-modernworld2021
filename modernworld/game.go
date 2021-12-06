package modernworld

import (
	"time"

	tl "github.com/JoelOtter/termloop"
)

type modernworld struct {
	*tl.Entity
	Game               *tl.Game
	Level              *tl.BaseLevel
	Arena              *Arena
	GameOverZone       *GameOverZone
	Hud                *Hud
	Hero               *Hero
	AlienCluster       *AlienCluster
	AlienLaserVelocity float64
	TimeDelta          float64
	RefreshSpeed       time.Duration
	Score              int
	Started            bool
	ScreenSizeNotOK    bool
}

func NewGame() *modernworld {
	modernworld := modernworld{
		Entity:             tl.NewEntity(0, 0, 1, 1),
		Game:               tl.NewGame(),
		Level:              tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack, Fg: tl.ColorWhite}),
		AlienLaserVelocity: 0.04,
		RefreshSpeed:       20,
		Score:              0,
	}

	modernworld.Game.Screen().SetFps(60)
	modernworld.Game.SetEndKey(tl.KeyBackspace)
	modernworld.Game.Screen().SetLevel(modernworld.Level)
	modernworld.Level.AddEntity(&modernworld)

	return &modernworld
}

func (modernworld *modernworld) Start() {
	go ShowTitleScreen(modernworld)
	modernworld.Game.Start()
}

func (modernworld *modernworld) Tick(event tl.Event) {
	if modernworld.Started == false && modernworld.ScreenSizeNotOK == false && event.Type == tl.EventKey && event.Key == tl.KeyEnter {
		go modernworld.initializeGame()
	}
}

func (modernworld *modernworld) initializeGame() {
	prepareScreen(modernworld)

	modernworld.Started = true

	modernworld.initHero()
	modernworld.initAliens()
	modernworld.initGameOverZone()
	modernworld.gameLoop()
}

func (modernworld *modernworld) initArena() {
	screenWidth, screenHeight := modernworld.getScreenSize()
	modernworld.Arena = newArena(screenWidth, screenHeight)
	modernworld.Level.AddEntity(modernworld.Arena)
}

func (modernworld *modernworld) initHud() {
	modernworld.Hud = NewHud(modernworld.Arena, modernworld.Level)
}

func (modernworld *modernworld) getScreenSize() (int, int) {
	screenWidth, screenHeight := modernworld.Game.Screen().Size()

	for screenWidth == 0 && screenHeight == 0 {
		time.Sleep(100 * time.Millisecond)
		screenWidth, screenHeight = modernworld.Game.Screen().Size()
	}

	return screenWidth, screenHeight
}

func (modernworld *modernworld) initHero() {
	modernworld.Hero = NewHero(modernworld.Arena)
	modernworld.Level.AddEntity(modernworld.Hero)
}

func (modernworld *modernworld) initAliens() {
	modernworld.AlienCluster = NewAlienCluster()
	SetPositionAndRenderAliens(modernworld.AlienCluster.Aliens, modernworld.Level, modernworld.Arena)
}

func (modernworld *modernworld) initGameOverZone() {
	modernworld.GameOverZone = CreateGameOverZone(modernworld.Arena, modernworld.Hero)
	modernworld.Level.AddEntity(modernworld.GameOverZone)
}

func (modernworld *modernworld) gameLoop() {
	for {
		if modernworld.Hero.IsDead() || modernworld.AlienCluster.IsAllAliensDead() {
			modernworld.Hero.animateHeroEndGame(modernworld.Level)
			modernworld.Started = false
			break
		}

		modernworld.updateLaserPositions()
		modernworld.RemoveDeadAliensAndIncrementScore()
		modernworld.updateAlienClusterPosition()
		modernworld.updateScore()
		modernworld.verifyGameOverZone()

		time.Sleep(modernworld.RefreshSpeed * time.Millisecond)
	}

	if modernworld.Hero.IsDead() {
		ShowGameOverScreen(modernworld)
	}

	if modernworld.AlienCluster.IsAllAliensDead() {
		modernworld.initializeGame()
	}
}

func (modernworld *modernworld) updateScore() {
	modernworld.Hud.UpdateScore(modernworld.Score)
}

func (modernworld *modernworld) updateAlienClusterPosition() {
	modernworld.AlienCluster.UpdateAliensPositions(modernworld.Game.Screen().TimeDelta(), modernworld.Arena)
	modernworld.AlienCluster.Shoot()
}

func (modernworld *modernworld) RemoveDeadAliensAndIncrementScore() {
	points := modernworld.AlienCluster.RemoveDeadAliensAndGetPoints(modernworld.Level)
	modernworld.addScore(points)
}

func (modernworld *modernworld) updateLaserPositions() {
	modernworld.updateHeroLasers()
	modernworld.updateAlienLasers()
	modernworld.removeLasers()
}

func (modernworld *modernworld) updateHeroLasers() {
	modernworld.updateLasers(modernworld.Hero.Lasers)
}

func (modernworld *modernworld) updateAlienLasers() {
	modernworld.TimeDelta += modernworld.Game.Screen().TimeDelta()

	if modernworld.TimeDelta >= modernworld.AlienLaserVelocity {
		modernworld.TimeDelta = 0
		modernworld.updateLasers(modernworld.AlienCluster.Lasers)
	}
}

func (modernworld *modernworld) updateLasers(lasers []*Laser) {
	for _, laser := range lasers {
		if laser.IsNew {
			modernworld.renderNewLaser(laser)
			continue
		}

		x, y := laser.Position()
		laser.SetPosition(x, y-laser.Direction)
	}
}

func (modernworld *modernworld) renderNewLaser(laser *Laser) {
	laser.IsNew = false
	modernworld.Level.AddEntity(laser)
}

func (modernworld *modernworld) removeLasers() {
	_, arenaY := modernworld.Arena.Position()
	_, arenaH := modernworld.Arena.Size()

	upperLimit := arenaY
	bottomLimit := arenaY + arenaH - 1

	modernworld.Hero.Lasers = modernworld.removeLaserOf(modernworld.Hero.Lasers, upperLimit)
	modernworld.AlienCluster.Lasers = modernworld.removeLaserOf(modernworld.AlienCluster.Lasers, bottomLimit)
}

func (modernworld *modernworld) removeLaserOf(lasers []*Laser, arenaLimit int) []*Laser {
	for index, laser := range lasers {
		_, y := laser.Position()
		isEndOfArena := y == arenaLimit

		if isEndOfArena || laser.HasHit {
			modernworld.Level.RemoveEntity(laser)

			if laser.HitAlienLaser {
				modernworld.addScore(laser.Points)
			}

			if index < len(lasers) {
				lasers = append(lasers[:index], lasers[index+1:]...)
			}
		}
	}

	return lasers
}

func (modernworld *modernworld) addScore(points int) {
	modernworld.Score += points
}

func (modernworld *modernworld) verifyGameOverZone() {
	if modernworld.GameOverZone.EnteredZone {
		modernworld.Hero.IsAlive = false
	}
}
