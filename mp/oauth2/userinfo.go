package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chanxuehong/wechat/oauth2"
)

const (
	Language_zh_CN = "zh_CN" // 简体中文
	Language_zh_TW = "zh_TW" // 繁体中文
	Language_en    = "en"    // 英文
)

const (
	SexUnknown = 0 // 未知
	SexMale    = 1 // 男性
	SexFemale  = 2 // 女性
)

type UserInfo struct {
	OpenId   string `json:"openid"`   // 用户的唯一标识
	Nickname string `json:"nickname"` // 用户昵称
	Sex      int    `json:"sex"`      // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	City     string `json:"city"`     // 普通用户个人资料填写的城市
	Province string `json:"province"` // 用户个人资料填写的省份
	Country  string `json:"country"`  // 国家, 如中国为CN

	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImageURL string `json:"headimgurl"`

	Privilege []string `json:"privilege"`         // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionId   string   `json:"unionid,omitempty"` // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
}

// 获取用户信息(需 scope 为 snsapi_userinfo).
//  lang 可能的取值是 zh_CN, zh_TW, en, 如果留空 "" 则默认为 zh_CN
//  httpClient 如果不指定则默认为 http.DefaultClient
func GetUserInfo(accessToken, openId, lang string, httpClient *http.Client) (info *UserInfo, err error) {
	switch lang {
	case "":
		lang = Language_zh_CN
	case Language_zh_CN, Language_zh_TW, Language_en:
	default:
		lang = Language_zh_CN
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	_url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) +
		"&openid=" + url.QueryEscape(openId) +
		"&lang=" + lang
	httpResp, err := httpClient.Get(_url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result struct {
		oauth2.Error
		UserInfo
	}
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	if result.ErrCode != oauth2.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.UserInfo
	return
}
