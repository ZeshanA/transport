module livedataloader

require (
	github.com/lib/pq v1.0.0
	github.com/tidwall/gjson v1.1.5
	github.com/tidwall/match v1.0.1 // indirect
	transport/lib v0.0.0
)

replace transport/lib => ../../lib
