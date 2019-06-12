module detector

require (
	github.com/VividCortex/ewma v1.1.1
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.2
	github.com/lib/pq v1.0.0
	github.com/pusher/pusher-http-go v4.0.0+incompatible
	github.com/stretchr/testify v1.3.0
	github.com/tidwall/gjson v1.2.1
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8 // indirect
	gopkg.in/guregu/null.v3 v3.4.0
	transport/lib v0.0.0
)

replace transport/lib => ../../lib
