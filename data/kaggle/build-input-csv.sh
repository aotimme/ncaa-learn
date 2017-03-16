#!/usr/bin/env bash

cat ./RegularSeasonDetailedResults.csv \
  | gocsv filter -c Season --regex "^2017$" \
  | gocsv select -c Wteam,Wscore,Lteam,Lscore \
  | gocsv join --left -c Wteam,Team_Id ./Teams.csv \
  | gocsv rename -c Team_Name --names "Team" \
  | gocsv select -c Wteam,Team_Id --exclude \
  | gocsv join --left -c Lteam,Team_Id ./Teams.csv \
  | gocsv rename -c Team_Name --names "Opponent" \
  | gocsv select -c Lteam,Team_Id --exclude \
  | gocsv rename -c Wscore,Lscore --names "Team Score,Opponent Score" \
  | gocsv select -c "Team,Team Score,Opponent,Opponent Score" > ./final-input.csv
