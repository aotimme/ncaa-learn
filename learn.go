package main

import (
  "fmt"
  "encoding/csv"
  "os"
  "io"
  "strconv"
)

type Game struct {
  teamA *Team
  teamB *Team
  teamAScore int64
  teamBScore int64
}

func (g *Game) PrintGame() {
  fmt.Printf("%v (%v) - %v (%v)\n", g.teamA.name, g.teamAScore, g.teamB.name, g.teamBScore)
}

type Team struct {
  name string
  forMean float64     // mean points for
  forVar float64      // variance of points for
  againstMean float64 // mean points against
  againstVar float64  // variance of points against
  games []*Game
}

func (t *Team) PrintTeam() {
  fmt.Printf("%v [%v]: %v(%v) %v(%v)\n", t.name, len(t.games), t.forMean, t.forVar, t.againstMean, t.againstVar)
}

func ReadData() (teams []*Team, games []*Game, err error) {
  //games = []*Game{}
  //teams = []*Team{}
  teamsMap := make(map[string]*Team)
  file, err := os.Open("2014_game_results.csv")
  if err != nil {
    return nil, nil, err
  }
  defer file.Close()
  reader := csv.NewReader(file)
  hasSeenHeader := false
  first := true
  for {
    record, err := reader.Read()
    if err == io.EOF {
      break
    } else if err != nil {
      return nil, nil, err
    }
    if hasSeenHeader {
      if first {
        teamAName := record[1]
        teamA, ok := teamsMap[teamAName]
        if !ok {
          teamA = &Team{name: teamAName}
          teamsMap[teamAName] = teamA
          teams = append(teams, teamA)
        }
        teamBName := record[4]
        teamB, ok := teamsMap[teamBName]
        if !ok {
          teamB = &Team{name: teamBName}
          teamsMap[teamBName] = teamB
          teams = append(teams, teamB)
        }
        teamAScore, _ := strconv.ParseInt(record[3], 0, 64)
        teamBScore, _ := strconv.ParseInt(record[5], 0, 64)
        game := &Game{teamA: teamA, teamB: teamB, teamAScore: teamAScore, teamBScore: teamBScore}
        games = append(games, game)
        teamA.games = append(teamA.games, game)
        teamB.games = append(teamB.games, game)
        first = false
      } else {
        first = true
      }
    } else {
      hasSeenHeader = true
    }
  }
  return
}

func main() {
  teams, games, err := ReadData()
  if err != nil {
    fmt.Println("error", err)
  }
  for i := 0; i < 5; i++ {
    games[i].PrintGame()
  }
  for i := 0; i < 5; i++ {
    teams[i].PrintTeam()
  }
  fmt.Printf("Number of games: %v\n", len(games))
  fmt.Printf("Number of teams: %v\n", len(teams))
  fmt.Println("Done!")
}
