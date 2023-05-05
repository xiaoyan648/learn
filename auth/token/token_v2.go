package token

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type StandardClaimsV2 struct {
	ExpiresAt int64  `json:"expires_at"`
	ID        uint64 `json:"id"`
	IssuedAt  int64  `json:"issued_at"`
	Issuer    string `json:"issuer"`
	AccountID uint32 `json:"account_id"`
	UID       uint64 `json:"user_id"`

	// 删除字段
	// Project       string   `json:"project"`
	DeviceID      string `json:"device_id"`      // 用户当前的CurrentDeviceID
	DeviceChannel string `json:"device_channel"` // 渠道信息
	Role          string `json:"role"`           //
	VipExpireAt   int64  `json:"vip_expire_at"`

	// 新增字段
	Ext     string `json:"ext"`
	VipType int    `json:"vip_type"`
}

func (s *StandardClaimsV2) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func ValidTokenV2(secretKey string, accessToken string) (*StandardClaimsV2, error) {
	token, err := jwt.ParseWithClaims(accessToken, &StandardClaimsV2{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		e, ok := err.(*jwt.ValidationError)
		if !ok || e.Inner == nil {
			return nil, err
		}

		return nil, e.Inner
	}

	if claims, ok := token.Claims.(*StandardClaimsV2); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token invalid")
}

func (s *StandardClaimsV2) Valid() error {
	now := NowFunc()
	// 校验过期时间
	// 负数表示永不过期
	if s.ExpiresAt > 0 {
		expireAt := time.Unix(s.ExpiresAt, 0)
		if now.After(expireAt) || now.Equal(expireAt) {
			return ErrTokenExpired
		}
	}

	// 校验发布时间
	issueAt := time.Unix(s.IssuedAt, 0)
	if now.Before(issueAt) {
		return ErrTokenBeforeIssueAt
	}

	return nil
}
