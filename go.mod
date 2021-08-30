module github.com/jrcasso/tugboat

replace github.com/jrcasso/tugboat => /workspaces/tugboat/

go 1.17

require (
	github.com/google/go-github v17.0.0+incompatible
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
	google.golang.org/appengine v1.1.0 // indirect
)
