package lambdabase

type Config struct {
	LambdaServerPort int `env:"_LAMBDA_SERVER_PORT" required:"true"`
}
