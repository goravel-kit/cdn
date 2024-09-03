package cdn

import (
	"fmt"
	"time"

	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"
)

type BaiShan struct {
	Token string
}

type BaiShanRefreshResponse struct {
	Code uint `json:"code"`
	Data any  `json:"data"`
}

type BaiShanUsageResponse struct {
	Code int `json:"code"`
	Data map[string]struct {
		Domain string   `json:"domain"`
		Data   [][]uint `json:"data"`
	} `json:"data"`
}

// RefreshUrl 刷新URL
func (b *BaiShan) RefreshUrl(urls []string) error {
	client := req.C()
	client.SetTimeout(60 * time.Second)

	refreshURL := "https://cdn.api.baishan.com/v2/cache/refresh?token=" + b.Token
	data := map[string]any{
		"urls": urls,
		"type": "url",
	}

	var resp BaiShanRefreshResponse
	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).Post(refreshURL)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("刷新URL失败: %d", resp.Code)
	}

	return nil
}

// RefreshPath 刷新路径
func (b *BaiShan) RefreshPath(paths []string) error {
	client := req.C()
	client.SetTimeout(60 * time.Second)

	refreshURL := "https://cdn.api.baishan.com/v2/cache/refresh?token=" + b.Token
	data := map[string]any{
		"urls": paths,
		"type": "dir",
	}

	var resp BaiShanRefreshResponse
	_, err := client.R().SetBody(data).SetSuccessResult(&resp).SetErrorResult(&resp).Post(refreshURL)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("刷新路径失败: %d", resp.Code)
	}

	return nil
}

// GetUsage 获取使用量
func (b *BaiShan) GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error) {
	client := req.C()
	client.SetTimeout(60 * time.Second)

	var usage BaiShanUsageResponse
	_, err := client.R().SetQueryParams(map[string]string{
		"token":      b.Token,
		"domains":    domain,
		"start_time": startTime.ToDateString(),
		"end_time":   endTime.ToDateString(),
	}).SetSuccessResult(&usage).Get("https://cdn.api.baishan.com/v2/stat/request/eachDomain")
	if err != nil {
		return 0, err
	}

	if usage.Code != 0 {
		return 0, fmt.Errorf("获取用量失败: %d", usage.Code)
	}

	sum := uint(0)
	for _, domain := range usage.Data {
		for _, data := range domain.Data {
			sum += data[1]
		}
	}

	return sum, nil
}
