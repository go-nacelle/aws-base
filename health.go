package lambdabase

type healthToken string

func (t healthToken) String() string {
	return "lambda-init"
}
