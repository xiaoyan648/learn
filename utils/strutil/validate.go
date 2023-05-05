package strutil

import "regexp"

// 校验手机号.
func ValidatePhoneNumber(phone string) bool {
	return regexp.MustCompile(`^[1]([3-9])[0-9]{9}$`).MatchString(phone)
}

// 校验邮箱(名称允许汉字、字母、数字，域名只允许英文域名).
func ValidateEmailNumber(email string) bool {
	return regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`).MatchString(email)
}

// 校验强密码(至少8-16个字符，至少1个大写字母，1个小写字母和1个数字，其他可以是任意字符).
func ValidateStrongPassword(s string) bool {
	return regexp.MustCompile(`.*[A-Z]+.*`).MatchString(s) && regexp.MustCompile(`.*[a-z]+.*`).MatchString(s) && regexp.MustCompile(`.*[0-9]+.*`).MatchString(s) && len(s) >= 8 && len(s) <= 16
}

func ValidateOnlyEnglish(s string) bool {
	return regexp.MustCompile("[A-Za-z_]+").MatchString(s)
}
