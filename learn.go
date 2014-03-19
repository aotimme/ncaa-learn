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
  meanFor float64     // mean points for
  precFor float64     // precision of points for
  meanAgainst float64 // mean points against
  precAgainst float64 // precision of points against
  games []*Game
}

func (t *Team) PrintTeam() {
  fmt.Printf("%v [%v]: %v(%v) %v(%v)\n", t.name, len(t.games), t.meanFor, 1.0/t.precFor, t.meanAgainst, 1.0/t.precAgainst)
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


func RunMAP(teams []*Team) {
  iterations := 10
  mu := 50.0
  tau := 5.0
  alpha := 2.0
  beta := 5.0
  for _, team := range teams {
    team.meanFor = mu
    team.meanAgainst = mu
    team.precFor = alpha/beta
    team.precAgainst = alpha/beta
  }
  for i := 0; i < iterations; i++ {
    // update means
    for _, team := range teams {
      precForTilde := tau
      precAgainstTilde := tau
      meanForInner := mu * tau
      meanAgainstInner := mu * tau
      for _, game := range team.games {
        other := game.teamA
        scoreFor := float64(game.teamBScore)
        scoreAgainst := float64(game.teamAScore)
        if other == team {
          other = game.teamB
          tmp := scoreFor
          scoreFor = scoreAgainst
          scoreAgainst = tmp
        }
        denomFor := 1.0/team.precFor + 1.0/other.precAgainst
        denomAgainst := 1.0/team.precAgainst + 1.0/other.precFor
        precForTilde += 1.0/denomFor
        precAgainstTilde += 1.0/denomAgainst
        meanForInner += (2.0 * scoreFor - other.meanAgainst/2.0)/denomFor
        meanAgainstInner += (2.0 * scoreAgainst - other.meanFor/2.0)/denomAgainst
      }
      meanForTilde := meanForInner / precForTilde
      meanAgainstTilde := meanAgainstInner / precAgainstTilde
      // update
      team.meanFor = meanForTilde
      team.meanAgainst = meanAgainstTilde
    }
    // update precisions
    for _, team := range teams {
      alphaForTilde := alpha + float64(len(team.games))/2.0
      alphaAgainstTilde := alpha + float64(len(team.games))/2.0
      betaForTilde := beta
      betaAgainstTilde := beta
      for _, game := range team.games {
        other := game.teamA
        scoreFor := float64(game.teamBScore)
        scoreAgainst := float64(game.teamAScore)
        if other == team {
          other = game.teamB
          tmp := scoreFor
          scoreFor = scoreAgainst
          scoreAgainst = tmp
        }
        forMult := scoreFor - (team.meanFor + other.meanAgainst)/2.0
        againstMult := scoreAgainst - (team.meanAgainst + other.meanFor)/2.0
        betaForTilde += forMult * forMult
        betaAgainstTilde += againstMult * againstMult
      }
      team.precFor = alphaForTilde/betaForTilde
      team.precAgainst = alphaAgainstTilde/betaAgainstTilde
    }
  }
}

func main() {
  teams, games, err := ReadData()
  if err != nil {
    fmt.Println("error", err)
  }
  //for i := 0; i < 5; i++ {
  //  games[i].PrintGame()
  //}
  //for i := 0; i < 5; i++ {
  //  teams[i].PrintTeam()
  //}
  //maxGames := 0
  //for _, team := range teams {
  //  if len(team.games) > maxGames {
  //    maxGames = len(team.games)
  //  }
  //}
  //gamesHist := make([]int, maxGames)
  //for _, team := range teams {
  //  gamesHist[len(team.games)-1]++
  //}
  //for i, num := range gamesHist {
  //  fmt.Printf("%v:\t%v\n", i+1, num)
  //}
  fmt.Printf("Number of games: %v\n", len(games))
  fmt.Printf("Number of teams: %v\n", len(teams))
  RunMAP(teams)
  for i := 0; i < 5; i++ {
    teams[i].PrintTeam()
  }
  fmt.Println("Done!")
}