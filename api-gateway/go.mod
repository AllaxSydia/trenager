module api-gateway

go 1.25.0

require (
	github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	google.golang.org/grpc v1.81.1
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/AllaxSydia/trenager/proto => ../proto
