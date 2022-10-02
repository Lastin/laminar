module github.com/digtux/laminar

go 1.15

require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/aws/aws-sdk-go v1.31.11
	github.com/go-git/go-git/v5 v5.3.0
	github.com/gobwas/glob v0.2.3
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/labstack/echo/v4 v4.6.3
	github.com/labstack/gommon v0.3.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v0.9.3
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/buntdb v1.1.2
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v3 v3.0.1
)

// replace github.com/go-git/go-git/v5 => ../../zeripath/go-git
