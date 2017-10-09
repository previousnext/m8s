package simpleyaml

import (
	"testing"
)

var data = []byte(`
name: smallfish
age: 99
float: 3.14159
bool: true
0: IntKey
emails:
   - xxx@xx.com
   - yyy@yy.com
bb:
    cc:
        dd:
            - 111
            - 222
            - 333
        ee: aaa
`)

func TestBool(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get("bool").Bool()
	if err != nil {
		t.Fatal("get yaml failed")
	}
	t.Log(v)
	if v != true {
		t.Fatal("match bool failed")
	}
}

func TestString(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get("name").String()
	if err != nil {
		t.Fatal("get yaml failed")
	}
	t.Log(v)
	if v != "smallfish" {
		t.Fatal("match name failed")
	}
}

func TestStringFromIntKey(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get(0).String()
	if err != nil {
		t.Fatal("get yaml failed")
	}

	t.Log(v)
	if v != "IntKey" {
		t.Fatal("match IntKey failed")
	}
}

func TestFloat(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}

	v, err := y.Get("float").Float()
	if err != nil {
		t.Fatal("get yaml failed", err)
	}
	t.Log(v)
	if v != 3.14159 {
		t.Fatal("match float failed")
	}
}

func TestInt(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get("age").Int()
	if err != nil {
		t.Fatal("get yaml failed")
	}
	t.Log(v)
	if v != 99 {
		t.Fatal("match age failed")
	}
}

func TestGetIndex(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get("bb").Get("cc").Get("dd").GetIndex(1).Int()
	t.Log(v)
	if err != nil {
		t.Fatal("match bb.cc.ee[1] failed")
	}
}

func TestString2(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get("bb").Get("cc").Get("ee").String()
	t.Log(v)
	if err != nil {
		t.Fatal("match bb.cc.ee failed")
	}
	if v != "aaa" {
		t.Fatal("bb.cc.ee not equal bbb")
	}
}

func TestGetPath(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.GetPath("bb", "cc", "ee").String()
	if err != nil {
		t.Fatal("get yaml failed")
	}
	t.Log(v)
	if v != "aaa" {
		t.Fatal("aa.bb.cc.ee not equal bbb")
	}
}

func TestArray(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	v, err := y.Get("emails").Array()
	if err != nil {
		t.Fatal("get yaml failed")
	}
	t.Log(v)
	if len(v) != 2 {
		t.Fatal("emails length not equal 2")
	}
}

func TestMap(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}
	if !y.IsMap() {
		t.Fatal("map check failed")
	}

	keys, err := y.GetMapKeys()
	if err != nil {
		t.Fatal("get keys from map is failed")
	}
	if len(keys) != 6 {
		t.Fatal("fail to check number of keys")
	}
}

func TestIsFound(t *testing.T) {
	y, err := NewYaml(data)
	if err != nil {
		t.Fatal("init yaml failed")
	}

	if !y.Get("name").IsFound() { // name is exists
		t.Fatal("found name failed")
	}

	if y.Get("xx").IsFound() { // xx is not exists
		t.Fatal("found xx failed")
	}

}
