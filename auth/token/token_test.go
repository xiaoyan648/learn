package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 新版本的token解析方法 对 老版本的token 的兼容性
// 测试结果：新增/删除字段都是兼容的，修改字段的类型不兼容
// 需要注意：添加字段的的场景，旧的token解析后新赠字段都为默认值，需要考虑到这种情况.
func TestUpdateTokenField(t *testing.T) {
	secretKey := "test-sk"

	// 老版本的token
	oldVersion := &StandardClaims{
		ID:            1,
		Issuer:        "test",
		IssuedAt:      1682661142,
		ExpiresAt:     -1,
		AccountID:     1,
		UID:           1,
		Project:       "test",
		DeviceID:      "test",
		DeviceChannel: "test",
		Role:          "test",
		VipExpireAt:   1682661142,
	}
	oldVersiontoken, err := generateToken(secretKey, oldVersion)
	assert.NoError(t, err)

	newVersion := &StandardClaimsV2{
		ID:        1,
		Issuer:    "test",
		IssuedAt:  1682661142,
		ExpiresAt: -1,
		AccountID: 1,
		UID:       1,
		// Project:       "test",
		DeviceID:      "test",
		DeviceChannel: "test",
		Role:          "test",
		VipExpireAt:   1682661142,
		// Ext:           "test",
		// VipType:       1,
	}

	got, err := ValidTokenV2(secretKey, oldVersiontoken)
	assert.NoError(t, err)
	assert.Equal(t, newVersion, got)
}

func TestMapClaimsToken(t *testing.T) {
	// jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"exp":  1682661142,
	// 	"id":   1,
	// 	"iat":  1682661142,
	// 	"iss":  "test",
	// 	"aid":  1,
	// })
}
