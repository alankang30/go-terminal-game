package main

import (
  "fmt"
  "log"
  "math/rand"
  "time"
  "github.com/gdamore/tcell/v2"
)

func drawString(screen tcell.Screen, x, y int, msg string) {
  for index, char := range msg {
    screen.SetContent(x+index, y, char, nil, tcell.StyleDefault)
  }
}

func setupCoins(level int, xmax, ymax int) []*Sprite {
  coins := make([]*Sprite, level+2)
  for index := range level+2 {
    invalid := true
    x, y := -1, -1
    for invalid {
      x = rand.Intn(xmax)
      y = rand.Intn(ymax)
      // ensure they aren't spawned inside of the text area
      if (x > 16 || y > 4) {
        break
      }
    }


    coins[index] = NewSprite(
      '©',
      x,
      y,
      tcell.StyleDefault.Foreground(tcell.ColorYellow),
    )
  }
  return coins
}

func setupBugs(level int, xmax, ymax, xspeed, yspeed int, ) []*Projectile {
  count := level + 1
  bugs := make([]*Projectile, count)
  for index := range count {
    sx := xspeed
    sy := yspeed
    randomDirection := rand.Intn(4)
    if randomDirection > 0 {
      sx *= -1
    }
    if randomDirection > 1 {
      sy *= -1
    }
    if randomDirection > 2 {
      sx *= -1
    }

    bugs[index] = NewProjectile('𓆣', rand.Intn(xmax), rand.Intn(ymax), 
    tcell.StyleDefault.Foreground(tcell.ColorLimeGreen), sx, sy)

  }
  return bugs
}


func game(screen tcell.Screen, xmax, ymax int) (int, int) {
  // screen, err := tcell.NewScreen()
  // if err != nil {
  //   log.Fatal(err)
  // }



  // err = screen.Init();
  // if err != nil {
  //   log.Fatal(err)
  // }

  // game init section
  player := NewSprite('𖨆', 10, 10, tcell.StyleDefault.Foreground(tcell.ColorBlue))
  //xmax, ymax := screen.Size()
  coins := setupCoins(1, xmax, ymax)

  //bug := NewSprite('𓆣', 15, 20, tcell.StyleDefault.Foreground(tcell.ColorRed))
  bugs := setupBugs(1, xmax, ymax, 1, 1)

  score := 0
  level := 1

	// fps counter initialization
	fps := 0
	frameCount := 0
	lastFPSUpdate := time.Now()
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

  //game loop

  playerAlive:= true
  for playerAlive {

    screen.Clear()

    //draw

    player.Draw(screen)

    for _, bug := range bugs {
      bug.Sprite.Draw(screen);
    }

    for _, coin := range coins {
      coin.Draw(screen)
    }

    //ui

    drawString( screen, 1, 1, fmt.Sprintf("Press Q to quit."),)
    drawString( screen, 1, 2, fmt.Sprintf("Score: %d", score),)
    drawString( screen, 1, 3, fmt.Sprintf("Level: %d", level),) 
    drawString(screen, 1, 4, fmt.Sprintf("FPS: %d", fps))

    screen.Show()

			// fps counter logic
			frameCount++
			if time.Since(lastFPSUpdate) >= time.Second {
				fps = frameCount
				frameCount = 0
				lastFPSUpdate = time.Now()
			}

			<-ticker.C

    //update logic

    playerMoved := false

      //player movement
      playerMoved = false
      if screen.HasPendingEvent() {
        ev := screen.PollEvent()
        switch ev := ev.(type) {
        case *tcell.EventKey:
          switch ev.Key() {
          case tcell.KeyEscape:
            playerAlive = false
            level = -1;
          }
          switch ev.Rune() {
          case 'q':
            playerAlive = false
            level = -1;
          case 'k', 'w':
            if player.Y > 0 {
              player.Y--
              playerMoved = true
            }
          case 'j', 's':
            if player.Y < ymax - 1 {
              player.Y++
              playerMoved = true
            }
          case 'h', 'a':
            if player.X > 0{
              player.X--
              playerMoved = true
            }
          case 'l', 'd':
            if player.X < xmax - 1 {
              player.X++
              playerMoved = true
            }
          }
        }
      } // player movement

      // coin collision check
      if playerMoved {
        coinCollectedIndex := -1
        for index, coin := range coins {
          if coin.X == player.X && coin.Y == player.Y {
            //collect the coin
            coinCollectedIndex = index
            score++
          }
        }
        // handle coin collision
        if coinCollectedIndex > -1 {
          // swap target with last 
          coins[coinCollectedIndex] = coins[len(coins)-1]

          //trim off last item
          coins = coins[0:len(coins)-1]

          if len(coins) == 0 {
            level++
            coins = setupCoins(level, xmax, ymax)
            bugs = setupBugs(level, xmax, ymax, 1, 1)
          }
        } 
      }// coin collision


      //bug collision check
      bug_idx := 0
      for _, bug := range bugs {
        bug.Update()

        if (bug.Sprite.X == player.X && bug.Sprite.Y == player.Y) {
          //draw an explosion
          boom := NewSprite('X', bug.Sprite.X, bug.Sprite.Y, tcell.StyleDefault.Foreground(tcell.ColorIndianRed));
          if (bug.Sprite.X < xmax - 2 ) {

            boom = NewSprite('X', bug.Sprite.X + 1, bug.Sprite.Y, tcell.StyleDefault.Foreground(tcell.ColorIndianRed));
            boom.Draw(screen)
          }
          if (bug.Sprite.X > 0) {

            boom = NewSprite('X', bug.Sprite.X - 1, bug.Sprite.Y, tcell.StyleDefault.Foreground(tcell.ColorIndianRed));
            boom.Draw(screen)
          }
          if (bug.Sprite.Y > 0) {

            boom = NewSprite('X', bug.Sprite.X, bug.Sprite.Y - 1, tcell.StyleDefault.Foreground(tcell.ColorIndianRed));
            boom.Draw(screen)
          }
          if (bug.Sprite.Y < ymax - 2 ) {

            boom = NewSprite('X', bug.Sprite.X, bug.Sprite.Y + 1, tcell.StyleDefault.Foreground(tcell.ColorIndianRed));
            boom.Draw(screen)
          }

          
          playerAlive = false;
          playerMoved = false;
          break;
        }

        //update bug 
        if bug.Sprite.X >= xmax || bug.Sprite.X <= 0{
          bug.SpeedX *= -1;
        }
        if bug.Sprite.Y >= ymax || bug.Sprite.Y <= 0{
          bug.SpeedY *= -1;
        }
        bugs[bug_idx] = bug
        bug_idx++
      } // bug update and collision

  } // game loop

  return score, level;
}

func main() {

  menu_screen, err := tcell.NewScreen();
  if err != nil {
    log.Fatal(err);
  }

  err = menu_screen.Init();
  if err != nil {
    log.Fatal(err)
  }

  xmax, ymax := menu_screen.Size()

  /* clean up screen */
  defer menu_screen.Fini()
  // defer calls at the end


  playing := true;
  for playing {

    //run the game 
    score, level := game(menu_screen, xmax, ymax );

    // quit during game
    if level == -1 {
      break;
    }

    //menu_screen.Clear();

    xmid := xmax / 2;
    ymid := ymax / 2;

    drawString( menu_screen,
      xmid - 17,
      ymid,
      fmt.Sprintf("You died with score %d on level %d!\n", score, level),
    )
    drawString( menu_screen,
      xmid - 8,
      ymid + 1,
      fmt.Sprintf("Press Q to quit."),
    )
    drawString( menu_screen,
      xmid - 9,
      ymid + 2,
      fmt.Sprintf("Press R to restart."),
    )
    menu_screen.Show();

    pending := true;
    for pending {
      ev := menu_screen.PollEvent()
      switch ev := ev.(type) {
      case *tcell.EventKey:
        switch ev.Rune() {
        case 'q', 'Q':
          pending = false;
          playing = false;
          break;
        case 'r', 'R':
          pending = false;
        default:
          continue;
        }
      }
    }

  }
}









