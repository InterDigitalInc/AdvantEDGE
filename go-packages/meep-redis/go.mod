module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/go-redis/redis v6.15.2+incompatible
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
