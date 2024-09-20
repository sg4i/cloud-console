module github.com/sg4i/cloud-console

go 1.22

require (
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.2.1
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1006
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts v1.0.1006
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
)

// 添加以下 replace 指令
replace github.com/sg4i/cloud-console => ./
