package expr

import (
	"fmt"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/assert"
)

// 使用 expr 表达式
func TestExpr(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		expr string
		want interface{}
	}{
		{name: "运算", args: map[string]interface{}{"foo": 1, "bar": 2}, expr: "foo + bar", want: 3},
		{name: "判断1", args: map[string]interface{}{"foo": 1, "bar": 2}, expr: "foo == 1", want: true},
		{name: "判断2", args: map[string]interface{}{"foo": 1, "bar": 2}, expr: "foo==1 && bar==1", want: false},
		{name: "判断3", args: map[string]interface{}{"foo": 1, "bar": 2}, expr: "foo != 2 && bar==1", want: true},
		{name: "判断-字段不存在", args: map[string]interface{}{"foo": 1, "bar": 2}, expr: "foo != 2 && bar==1 && test=0", want: false},
		{name: "范围1", args: map[string]interface{}{"role": "user", "age": 18}, expr: `role in ["vistor", "admin"]`, want: false},
		{name: "范围2", args: map[string]interface{}{"role": "user", "age": 18}, expr: `role in ["user", "admin"]`, want: true},
		{name: "范围3", args: map[string]interface{}{"role": "user", "age": 18}, expr: `role in ["user", "admin"] && age > 18`, want: false},
		{name: "范围4", args: map[string]interface{}{"role": "user", "age": 18}, expr: `role in ["user", "admin"] && age in 12..18`, want: true},
		{name: "范围5", args: map[string]interface{}{"role": "user", "age": 18}, expr: `role in ["user", "admin"] && age in 12..17`, want: false},
		{name: "字符串匹配1", args: map[string]interface{}{"desc": "like a apple"}, expr: `desc contains "apple"`, want: true},
		{name: "字符串匹配2", args: map[string]interface{}{"desc": "panic: xxxxxx"}, expr: `desc startsWith "panic"`, want: true},
		{name: "正则表达式1", args: map[string]interface{}{"regex": "13978783452"}, expr: `regex matches "^[0-9]{11}$"`, want: true},
		{name: "正则表达式2", args: map[string]interface{}{"regex": "1397878345a"}, expr: `regex matches "^[0-9]{11}$"`, want: false},
		{name: "slice", args: map[string]interface{}{"slice": []int{1, 2, 3, 4}}, expr: `slice[1:3] == [2,3]`, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expr.Eval(tt.expr, tt.args)
			assert.NoErrorf(t, err, "expr.Eval error: %v", err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func ExampleExprFunctionStruct() {
	type Segment struct {
		Origin int
	}
	type Meta struct {
		Tags map[string]string
	}
	type Env struct {
		Meta
		Segments []*Segment
		Marker   string
	}
	code := `all(Segments, {.Origin > 1}) && len(Segments) > 0 && Tags["foo"] startsWith "bar"`
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := Env{
		Meta: Meta{
			Tags: map[string]string{
				"foo": "bar",
			},
		},
		Segments: []*Segment{
			{Origin: 10},
			{Origin: 11},
		},
		Marker: "test",
	}
	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleExprFunctionMap() {
	env := map[string]interface{}{
		"tags": []string{"foo", "bar"},
		"ids":  []int{1, 2, 3},
		"user": []struct{ Id int }{{1}, {2}},
	}

	// program, err := expr.Compile("all(user, {.Id >1 })", expr.Env(env)) // false
	// program, err := expr.Compile("filter(user, {.Id > 1})", expr.Env(env)) // [{2}]
	// program, err := expr.Compile("map(user, {.Id > 1})", expr.Env(env)) // [false true]
	// program, err := expr.Compile("one(user, {.Id > 1})", expr.Env(env)) // true
	program, err := expr.Compile("map(ids, {# > 1})", expr.Env(env)) // [false true true] # 代表当前元素
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}
