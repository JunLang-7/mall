package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/service/do"
)

const (
	larkHost        = "https://open.feishu.cn"
	HttpAccept      = "application/json"
	HttpContentType = "application/json;charset=utf-8"
)

type ILark interface {
	GetLarkUserInfo(ctx context.Context, userAccessToken string) (*do.LarkUserInfo, error)
	GetLarkUserAccessToken(ctx context.Context, appCode int32, code, redirectUrl, scope string) (*do.LarkUserAccessToken, error)
	GetLarkTenantAccessToken(ctx context.Context, appCode int32) (*do.LarkUserAccessToken, error)
}

type GetTokenFunc func(ctx context.Context, force bool) (string, error)

type Lark struct {
	conf *config.Config
}

func NewLark(adaptor adaptor.IAdaptor) *Lark {
	return &Lark{conf: adaptor.GetConf()}
}

// GetLarkUserInfo 获取飞书用户信息
func (l *Lark) GetLarkUserInfo(ctx context.Context, userAccessToken string) (*do.LarkUserInfo, error) {
	url := fmt.Sprintf("%s/open-apis/authen/v1/user_info", larkHost)
	headers := map[string]string{
		"Content-Type":  HttpContentType,
		"Authorization": fmt.Sprintf("Bearer %s", userAccessToken),
		"Accept":        HttpAccept,
	}

	type Response struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data do.LarkUserInfo `json:"data"`
	}
	var resp Response
	// 发送 HTTP 请求并解析响应
	err := doRequest(ctx, http.MethodGet, url, headers, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark get user info failed: code=%d msg=%s", resp.Code, resp.Msg)
	}
	return &resp.Data, nil
}

// GetLarkUserAccessToken 获取飞书用户 access token
func (l *Lark) GetLarkUserAccessToken(ctx context.Context, appCode int32, code, redirectUrl, scope string) (*do.LarkUserAccessToken, error) {
	url := fmt.Sprintf("%s/open-apis/authen/v2/oauth/token", larkHost)
	appConf, ok := l.conf.AppConf[appCode]
	if !ok {
		return nil, fmt.Errorf("lark app config not found: appCode=%d", appCode)
	}

	body := map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     appConf.AppID,
		"client_secret": appConf.AppSecret,
		"code":          code,
		"redirect_uri":  redirectUrl,
		"scope":         scope,
	}
	headers := map[string]string{
		"Content-Type": HttpContentType,
		"Accept":       HttpAccept,
	}
	type Response struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data do.LarkUserAccessToken `json:"data"`
		do.LarkUserAccessToken
	}
	var resp Response
	// 发送 HTTP 请求并解析响应
	err := doRequest(ctx, http.MethodPost, url, headers, body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark get user access token failed: code=%d msg=%s", resp.Code, resp.Msg)
	}
	// 飞书接口在不同场景下返回的 access token 可能在不同字段，这里兼容两种情况
	if resp.Data.AccessToken != "" {
		return &resp.Data, nil
	}
	if resp.LarkUserAccessToken.AccessToken != "" {
		return &resp.LarkUserAccessToken, nil
	}
	return nil, errors.New("lark user access token is empty")
}

// GetLarkTenantAccessToken 获取飞书租户 access token
func (l *Lark) GetLarkTenantAccessToken(ctx context.Context, appCode int32) (*do.LarkUserAccessToken, error) {
	url := fmt.Sprintf("%s/open-apis/auth/v3/tenant_access_token/internal", larkHost)
	appConf, ok := l.conf.AppConf[appCode]
	if !ok {
		return nil, fmt.Errorf("lark app config not found: appCode=%d", appCode)
	}
	body := map[string]interface{}{
		"app_id":     appConf.AppID,
		"app_secret": appConf.AppSecret,
	}
	headers := map[string]string{
		"Content-Type": HttpContentType,
		"Accept":       HttpAccept,
	}
	type Response struct {
		Code              int    `json:"code"`
		Msg               string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
		Expire            int64  `json:"expire"`
	}
	var resp Response
	// 发送 HTTP 请求并解析响应
	err := doRequest(ctx, http.MethodPost, url, headers, body, &resp)
	if err != nil {
		return nil, err
	}
	// 飞书接口在不同场景下返回的 access token 可能在不同字段，这里兼容两种情况
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark get tenant access token failed: code=%d msg=%s", resp.Code, resp.Msg)
	}
	if resp.TenantAccessToken == "" {
		return nil, errors.New("lark tenant access token is empty")
	}
	return &do.LarkUserAccessToken{
		Code:              int64(resp.Code),
		AccessToken:       resp.TenantAccessToken,
		TenantAccessToken: resp.TenantAccessToken,
		ExpiresIn:         resp.Expire,
	}, nil
}

// doRequest 发送 HTTP 请求并解析响应
func doRequest(ctx context.Context, method, url string, headers map[string]string, body interface{}, resp interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(string(payload))
	}
	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}
	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// 发送 HTTP 请求
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	// 检查 HTTP 状态码和响应体
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return fmt.Errorf("lark request failed: status=%d body=%s", httpResp.StatusCode, string(respBody))
	}
	if len(respBody) == 0 {
		return errors.New("lark request body is empty")
	}
	return json.Unmarshal(respBody, resp)
}
