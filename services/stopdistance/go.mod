module stopdistance

require (
	github.com/stretchr/testify v1.3.0
	github.com/tidwall/gjson v1.2.1
	googlemaps.github.io/maps v0.0.0-20190311183511-743053230cec
	transport/lib v0.0.0
	transport/services/labeller v0.0.0
)

replace transport/lib => ../../lib

replace transport/services/labeller => ../../services/labeller
