package grpc

// type Config struct {
// 	AuthGrpcClientAddress string `koanf:"auth_grpc_client_address"`
// }

type Config struct {
	Targets struct {
		Auth string `koanf:"auth"`
	} `koanf:"targets"`
}
