package cdn

import (
	"errors"
	"fmt"
	"time"

	"github.com/goravel/framework/support/carbon"
	"github.com/imroc/req/v3"
)

type GoEdge struct {
	API, AccessKeyID, AccessKey string
}

type GoEdgeCommonResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type GoEdgeTokenResponse struct {
	Code int `json:"code"`
	Data struct {
		Token     string `json:"token"`
		ExpiresAt int    `json:"expiresAt"`
	} `json:"data"`
	Message string `json:"message"`
}

type GoEdgeServerResponse struct {
	Code int `json:"code"`
	Data struct {
		Servers []struct {
			Id              int    `json:"id"`
			IsOn            bool   `json:"isOn"`
			Name            string `json:"name"`
			FirstServerName string `json:"firstServerName"`
		} `json:"servers"`
	} `json:"data"`
	Message string `json:"message"`
}

type GoEdgeUsageResponse struct {
	Code int `json:"code"`
	Data struct {
		ServerDailyStat struct {
			ServerId      int `json:"serverId"`
			Bytes         int `json:"bytes"`
			CachedBytes   int `json:"cachedBytes"`
			CountRequests int `json:"countRequests"`
		} `json:"serverDailyStat"`
	} `json:"data"`
	Message string `json:"message"`
}

// RefreshUrl 刷新URL
func (r *GoEdge) RefreshUrl(urls []string) error {
	return r.RefreshPath(urls)
}

// RefreshPath 刷新路径
func (r *GoEdge) RefreshPath(paths []string) error {
	client, err := r.getClient()
	if err != nil {
		return err
	}

	var refreshResponse GoEdgeCommonResponse
	_, err = client.R().SetBodyJsonMarshal(map[string]any{
		"type":    "purge",
		"keyType": "prefix",
		"keys":    paths,
	}).
		SetSuccessResult(&refreshResponse).
		Post("/HTTPCacheTaskService/createHTTPCacheTask")
	if err != nil {
		return err
	}

	if refreshResponse.Code != 200 {
		return fmt.Errorf("刷新失败: %s", refreshResponse.Message)
	}

	return nil
}

// GetUsage 获取使用量
func (r *GoEdge) GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error) {
	id, err := r.getServerIDByDomain(domain)
	if err != nil {
		return 0, err
	}

	client, err := r.getClient()
	if err != nil {
		return 0, err
	}

	var request = map[string]any{
		"serverId": id,
		"dayFrom":  startTime.ToShortDateString(),
		"dayTo":    endTime.ToShortDateString(),
	}
	var usageResponse GoEdgeUsageResponse
	_, err = client.R().
		SetBodyJsonMarshal(request).
		SetSuccessResult(&usageResponse).
		Post("/ServerDailyStatService/sumServerDailyStats")
	if err != nil {
		return 0, err
	}

	if usageResponse.Code != 200 {
		return 0, fmt.Errorf("获取用量失败: %s", usageResponse.Message)
	}

	return uint(usageResponse.Data.ServerDailyStat.CountRequests), nil
}

func (r *GoEdge) getServerIDByDomain(domain string) (int, error) {
	client, err := r.getClient()
	if err != nil {
		return 0, err
	}

	var response GoEdgeServerResponse
	_, err = client.R().SetSuccessResult(&response).Post("/ServerService/findAllUserServers")
	if err != nil {
		return 0, err
	}

	if response.Code != 200 {
		return 0, errors.New(response.Message)
	}

	for _, server := range response.Data.Servers {
		if server.Name == domain {
			return server.Id, nil
		}
	}

	return 0, errors.New("未找到该域名对应的服务ID")
}

// getClient 获取客户端
func (r *GoEdge) getClient() (*req.Client, error) {
	client := req.C()
	client.SetTimeout(10 * time.Second)
	client.SetCommonRetryCount(2)
	client.ImpersonateSafari()
	client.SetBaseURL(r.API)

	// 换取Token
	data := map[string]string{
		"type":        "user",
		"accessKeyId": r.AccessKeyID,
		"accessKey":   r.AccessKey,
	}
	var tokenResponse GoEdgeTokenResponse
	_, err := client.R().
		SetBodyJsonMarshal(data).
		SetSuccessResult(&tokenResponse).
		Post("/APIAccessTokenService/getAPIAccessToken")
	if err != nil {
		return nil, err
	}
	if tokenResponse.Code != 200 {
		return nil, errors.New(tokenResponse.Message)
	}

	client.SetCommonHeader("X-Edge-Access-Token", tokenResponse.Data.Token)
	return client, nil
}
