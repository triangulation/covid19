# covid19

covid-19 extracts and transforms covid-19 time-series data, with the goal
of visualising it in an intuitive and performant way.

Data source: [github.com/datasets/covid-19/data/time-series-19-covid-combined.csv](https://github.com/datasets/covid-19/raw/master/data/time-series-19-covid-combined.csv)

Target: JSON or Proto with time series data of covid stats per country
optimised for size.

## Prerequisites

- [go](https://golang.org/doc/go1.14),
- [golangci-lint](https://github.com/golangci/golangci-lint/releases/tag/v1.23.6),
- GNU make

## Development

- Build with `make`
- View build options with `make help`

## Data visualisation notes

- simple time lines:
  - cumulative/new cases/deaths
  - rolling 3/5/7 day average
  - log or linear y-axes
  - absolute and per 1M/10K population number
  - choose country/ies (drop-down/search; keep pref in localStorage)
- bar charts, maybe stacked:
  - ranking of countries by monthly/weekly/rolling-month highest cases per capita
  - see https://twitter.com/RARohde/status/1289562024830173185?s=19
- geography:
  - selection on map
  - Austrialia: state, lga, suburb selection
- data:
  - compact JSON encoding: ~120K gzipped for all countries (too much)
  - try protos
  - try country pre-selection if needed
