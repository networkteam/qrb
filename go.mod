module github.com/networkteam/qrb

go 1.20

// Note that this module is not using dependencies, besides tests and for the qrbpgx package.
require (
	// Used only for tests
	github.com/stretchr/testify v1.8.2
	// Used only for tests of qrbsql
	github.com/DATA-DOG/go-sqlmock v1.5.0
	// Used only in qrbpgx
	github.com/jackc/pgx/v5 v5.3.1
	// Used only for tests of qrbpx
	github.com/pashagolub/pgxmock/v2 v2.6.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
