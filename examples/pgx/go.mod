module github.com/networkteam/qrb/examples/pgx

require (
	github.com/jackc/pgx/v5 v5.6.0
	github.com/networkteam/qrb v0.0.0-00010101000000-000000000000
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.25.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/networkteam/qrb => ./../..

go 1.20
