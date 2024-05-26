module github.com/eikarna/gotps

go 1.22.2

require (
	github.com/codecat/go-libs v0.0.0-20210906174629-ffa6674c8e05
	github.com/eikarna/GoDat v0.0.0-20240525154534-5254d6820a04
	github.com/eikarna/gotops v0.0.0-20240427030241-b3f97fb7be43
)

require (
	github.com/fatih/color v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/eikarna/gotops => ../gotops
replace github.com/eikarna/GoDat => ../GoDat
