package lambdabase

type Config struct {
	LambdaServerPort int `env:"_lambda_server_port" required:"true"`
}
