package main

import (
	"testing"
)

func TestJSONParse(t *testing.T) {
	ret := FromJsonString[Result[string]](`{"ok":true,"msg":"","data":"hello world"}`)
	if !ret.Ok {
		t.Errorf("FromJsonString error: %s", ret.Msg)
	}
	if ret.Msg != "" {
		t.Errorf("FromJsonString error: %s", ret.Msg)
	}
	if ret.Data != "hello world" {
		t.Errorf("FromJsonString error: %s", ret.Msg)
	}
}

type TestJSONData struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestJSONString(t *testing.T) {
	if ret := ToJsonString(OkResult[string]("x")); ret != "{\"ok\":true,\"msg\":\"\",\"data\":\"x\"}" {
		t.Errorf("ToJsonString error: %s", ret)
	}

	if ret := ToJsonString(TestJSONData{A: 1, B: "b"}); ret != "{\"a\":1,\"b\":\"b\"}" {
		t.Errorf("ToJsonString error: %s", ret)
	}

	if ret := ToJsonString(OkResult[TestJSONData](TestJSONData{A: 1, B: "b"})); ret != "{\"ok\":true,\"msg\":\"\",\"data\":{\"a\":1,\"b\":\"b\"}}" {
		t.Errorf("ToJsonString error: %s", ret)
	}
}
