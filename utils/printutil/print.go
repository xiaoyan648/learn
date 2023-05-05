package printutil

import (
	"bytes"
	"encoding/json"
)

// TODO 上线前删除 直观的打印非数组的结构
func ToIndentJSON(obj interface{}) string {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, bs, "", "\t")
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
