package main

import (
  "net/http"
  "os"
  "io"
  "encoding/csv"
  "strconv"
  "encoding/json"
  "fmt"
  "math"
  "regexp"
)

type Team struct {
  Name string `json:"name"`
  MeanFor float64 `json:"mean_for"`         // mean points for
  SdFor float64 `json:"sd_for"`             // sd of points for
  MeanAgainst float64 `json:"mean_against"` // mean points against
  SdAgainst float64 `json:"sd_against"`     // sd of points against
  Wins int64 `json:"wins"`
  Losses int64 `json:"losses"`
}

var teams = make([]*Team, 0)
var teamsMap = make(map[string]*Team)

func ReadRankings(filename string) error {
  file, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer file.Close()
  reader := csv.NewReader(file)
  for {
    record, err := reader.Read()
    if err == io.EOF {
      break
    } else if err != nil {
      return err
    }
    name := record[0]
    wins, err := strconv.ParseInt(record[1], 0, 64)
    losses, err := strconv.ParseInt(record[2], 0, 64)
    meanFor, err := strconv.ParseFloat(record[3], 64)
    sdFor, err := strconv.ParseFloat(record[4], 64)
    meanAgainst, err := strconv.ParseFloat(record[5], 64)
    sdAgainst, err := strconv.ParseFloat(record[6], 64)
    if err != nil {
      return err
    }
    team := &Team{
      Name: name, Wins: wins, Losses: losses, MeanFor: meanFor,
      MeanAgainst: meanAgainst, SdFor: sdFor, SdAgainst: sdAgainst,
    }
    teams = append(teams, team)
    teamsMap[name] = team
  }
  return nil
}

func ShowRankings(w http.ResponseWriter, r *http.Request) {
  js, err := json.Marshal(teams)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func SearchTeams(w http.ResponseWriter, r *http.Request) {
  queryMap := r.URL.Query()
  names, _ := queryMap["name"]
  name := names[0]
  // TODO: check if name exists...
  results := make([]string, 0)
  numResults := 0
  for _, team := range teams {
    matched, err := regexp.MatchString("(?i).*" + name + ".*", team.Name)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    if matched {
    //if strings.Contains(team.Name, name) {
      numResults++
      results = append(results, team.Name)
      if numResults > 10 {
        break
      }
    }
  }
  js, err := json.Marshal(results)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func MatchupHandler(w http.ResponseWriter, r *http.Request) {
  queryMap := r.URL.Query()
  // TODO: check for existence...
  homes, _ := queryMap["home"]
  home := homes[0]
  aways, _ := queryMap["away"]
  away := aways[0]
  homeTeam := teamsMap[home]
  awayTeam := teamsMap[away]
  homeDiff := homeTeam.MeanFor - homeTeam.MeanAgainst
  awayDiff := awayTeam.MeanFor - awayTeam.MeanAgainst
  var winner, loser *Team
  if homeDiff > awayDiff {
    winner, loser = homeTeam, awayTeam
  } else {
    winner, loser = awayTeam, homeTeam
  }
  response := make(map[string]interface{})
  response["winner"] = winner
  response["loser"] = loser
  winnerScore := (winner.MeanFor + loser.MeanAgainst)/2.0
  loserScore := (loser.MeanFor + winner.MeanAgainst)/2.0
  mean := winnerScore - loserScore
  sd := math.Sqrt((winner.SdFor*winner.SdFor + winner.SdAgainst*winner.SdAgainst + loser.SdFor*loser.SdFor + loser.SdAgainst*loser.SdAgainst)/4.0)
  probability := (1.0 / 2.0) * (1 + math.Erf((mean)/(sd*math.Sqrt2)))

  score := make(map[string]float64)
  score["winner"] = winnerScore
  score["loser"] = loserScore
  score["probability"] = probability
  response["score"] = score
  js, err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
  return
}

func main() {
  ReadRankings("../data/2014_rankings.csv")
  fmt.Printf("Teams Loaded: %v\n", len(teams))
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.HandleFunc("/search", SearchTeams)
  http.HandleFunc("/matchup", MatchupHandler)
  http.HandleFunc("/", ServeIndex)
  http.ListenAndServe("localhost:4000", nil)
}
