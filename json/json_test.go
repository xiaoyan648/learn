package json

import (
	"encoding/json"
	"fmt"
	"testing"
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
