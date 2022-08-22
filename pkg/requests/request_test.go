package requests

import (
	"testing"
)

func Test_RegiterSub(t *testing.T) {
	x := Subscriber{
		CallBack: CallBack{
			HttpMethod:  "GET",
			CallbackUrl: "http://localhost:8083/pong",
		},
		Channel: "c4",
	}
	//y := "http://localhost:9090/register/subscriber"
	err := RegisterSub(x)
	if err != nil {
		t.Log(err.Error())
	}

}
func Test_RegiterPub(t *testing.T) {
	x := Publisher{Channel: "c1"}
	y := "http://localhost:9090/register/publisher"
	key, err := RegisterPub(y, x)
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(key)

}
func Test_UpdateSubs(t *testing.T) {
	x := Publisher{Channel: "c1"}
	y := "http://localhost:9090/register/publisher"
	z, err := RegisterPub(y, x)

	err = UpdateSubs("http://localhost:9090/publish", "HelloWorld", z, "c1")
	if err != nil {
		t.Log(err.Error())
	}

}
