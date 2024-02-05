package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

type Post struct {
	Id      int      `json:"ID,string"`
	Content string   `json:"content"`
	Author  string   `json:"author"`
	Label   []string `json:"label,omitempty"`
}

type Post2 struct {
	Id      json.Number `json:"ID"`
	Content string      `json:"content"`
	Author  string      `json:"author"`
	Label   []string    `json:"label,omitempty"`
}

// 是否可以解析 驼峰到蛇形
func Test(t *testing.T) {
	data := []byte(`{"ID": "1", "content": "content", "author": "author", "label": ["label1", "label2"]}`)
	post := Post{}
	_ = json.Unmarshal(data, &post)
	fmt.Printf("%+v\n", post)

	raw, _ := json.Marshal(post)
	fmt.Printf("%s\n", raw)

	post2 := Post2{}
	_ = json.Unmarshal(data, &post2)
	fmt.Printf("%+v\n", post2)

	raw2, _ := json.Marshal(post2)
	fmt.Printf("%s\n", raw2)
}

type Person struct {
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age,omitempty"`
	Email string `json:"email,omitempty"`
	Post
}

func TestEmbed(t *testing.T) {
	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
		Post: Post{
			Content: "aaa",
		},
	}

	jsonData, err := json.Marshal(person)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(jsonData))
	// 输出：{"name":"Alice","age":30,"email":"alice@example.com"}

	person.Name = ""
	person.Age = 0
	jsonData, err = json.Marshal(person)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(jsonData))
	// 输出：{"age":0,"email":""}
}

type ReadRecordMessage struct {
	Date       string `json:"date,omitempty"`
	UID        uint64 `json:"uid,omitempty"`
	Project    string `json:"project,omitempty"`
	Event      string `json:"event,omitempty"`
	Extra      string `json:"extra,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type Extra struct {
	UID uint64 `json:"uid,omitempty"`
}

type ExtraStr struct {
	UID  string      `json:"uid,omitempty"`
	UID2 json.Number `json:"uid2,omitempty"`
}

/*
1.出现情况是：在使用interface{}进行unmarshal，超过9007199254740991的int会导致精度丢失。
2. 分析
最大的安全整数是52位尾数全为1且指数部分为最小 0x001F FFFF FFFF FFFF
float64可以存储的最大整数是52位尾数全位1且指数部分为最大 0x07FEF FFFF FFFF FFFF
复制
(0x001F FFFF FFFF FFFF)16 = (9007199254740991)10
(0x07EF FFFF FFFF FFFF)16 = (9218868437227405311)10
1.
2.
也就是理论上数值超过9007199254740991就可能会出现精度缺失。

10进制数值的有效数字是16位，一旦超过16位基本上缺失精度是没跑了，回过头看我处理的id是20位长度，所以必然出现精度缺失。
*/
func TestJsonUInt64(t *testing.T) {
	// 1.指定uid类型 不会有丢失问题
	rrm := &ReadRecordMessage{
		Date:       time.Now().Format("2006-01-02"),
		UID:        4817391248075129170,
		Project:    "demo",
		Event:      "read",
		Extra:      `{"uid":4817391248075129170}`,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	raw, _ := json.Marshal(rrm)
	fmt.Println(string(raw))

	rrm2 := &ReadRecordMessage{}
	json.Unmarshal(raw, rrm2)
	fmt.Println(rrm2)

	extra := &Extra{}
	json.Unmarshal([]byte(rrm2.Extra), extra)
	fmt.Println(extra)

	// 2.如果按照interface{}解析，会丢失精度
	// var extra2 interface{}
	var extra2 map[string]interface{}
	err := json.Unmarshal([]byte(rrm2.Extra), &extra2)
	if err != nil {
		fmt.Println("error:", err)
	}
	dealStr, err := json.Marshal(extra2)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("extra2:%v, dealStr:%s\n", extra2, string(dealStr)) // extra2:map[uid:4.817391248075129e+18], dealStr:{"uid":4817391248075129000}

	// 3.如果不得不使用 interface{}, 需要把字段改为string或 json.number,或使用专门的decoder
	var extra4 map[string]interface{}
	d := json.NewDecoder(strings.NewReader(rrm2.Extra))
	d.UseNumber()
	d.Decode(&extra4)
	dealStr, err = json.Marshal(extra4)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("extra4:%v, dealStr:%s", extra4, string(dealStr)) // extra2:map[uid:4.817391248075129e+18], dealStr:{"uid":4817391248075129000}

	// 4.使用json.number 其实和使用专门的decoder一个原理，decoder也是把flaot64转成了json.number
	// extra3 := &ExtraStr{
	// 	UID:  "4817391248075129170",
	// 	UID2: json.Number("4817391248075129170"),
	// }
}
