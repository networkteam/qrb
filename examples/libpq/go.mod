module github.com/networkteam/qrb/examples/libpq

require github.com/networkteam/qrb v0.0.0-00010101000000-000000000000

require (
	github.com/lib/pq v1.10.9
	github.com/urfave/cli/v2 v2.27.2
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
)

replace github.com/networkteam/qrb => ./../..

go 1.20
