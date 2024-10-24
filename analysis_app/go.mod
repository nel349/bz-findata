module github.com/nel349/coinbase-analysis

go 1.23

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	gopkg.in/go-jose/go-jose.v2 v2.6.3
)

require github.com/stretchr/testify v1.8.4 // indirect

replace github.com/nel349/bz-findata/pkg/exchange/coinbase v0.9.3 => /Users/norman/Development/go/bz-findata/pkg/exchange/coinbase

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
)
