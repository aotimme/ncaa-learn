NCAA Learn
==========

Coded up a basic model that tries to determine every team's latent ability.
Where "latent ability" means mean number of points the team scores and
mean number of points the team allows.

### Usage

```bash
# Make sure you have `gocsv` installed
cd ./data/kaggle/
bash build-input-csv.sh

cd ../../compute
go run learn.go --input ../data/kaggle/final-input.csv --output ../data/kaggle/2017-rankings.csv

cd ../server
go run app.go --results ../data/kaggle/2017-rankings.csv

# Visit localhost:4000 in your browser.
```
