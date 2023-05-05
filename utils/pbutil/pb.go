package pbutil

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromProtoTimestamp(t *timestamp.Timestamp) (tt time.Time) {
	return t.AsTime()
}

func ToProtoTimestamp(t time.Time) (tt *timestamp.Timestamp) {
	if t.IsZero() {
		return nil
	}

	return timestamppb.New(t)
}

func ToProtoString(str string) *wrappers.StringValue {
	return &wrappers.StringValue{Value: str}
}

func ToProtoUInt32(i uint32) *wrappers.UInt32Value {
	return &wrappers.UInt32Value{Value: i}
}

func ToProtoUInt64(i uint64) *wrappers.UInt64Value {
	return &wrappers.UInt64Value{Value: i}
}

func ToProtoInt64(i int64) *wrappers.Int64Value {
	return &wrappers.Int64Value{Value: i}
}

func ToProtoInt32(i int32) *wrappers.Int32Value {
	return &wrappers.Int32Value{Value: i}
}

func ToProtoBool(b bool) *wrappers.BoolValue {
	return &wrappers.BoolValue{Value: b}
}

func ToProtoBytes(bytes []byte) *wrappers.BytesValue {
	return &wrappers.BytesValue{Value: bytes}
}

func ToProtoFloat(f float32) *wrappers.FloatValue {
	return &wrappers.FloatValue{Value: f}
}
