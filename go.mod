module github.com/cloud-barista/cb-tumblebug

go 1.13

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/astaxie/beego v1.12.3
	github.com/bramvdbogaerde/go-scp v0.0.0-20200820121624-ded9ee94aef5
	github.com/cloud-barista/cb-log v0.3.0-espresso
	github.com/cloud-barista/cb-spider v0.3.0-espresso
	github.com/cloud-barista/cb-store v0.3.0-espresso
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/etcd-io/etcd v3.3.25+incompatible // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo/v4 v4.1.17
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.9.0
	github.com/shirou/gopsutil v3.20.11+incompatible
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
	github.com/tidwall/gjson v1.6.4
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
