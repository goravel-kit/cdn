package cdn

import (
	"fmt"

	"github.com/goravel/framework/support/carbon"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	"github.com/spf13/cast"
)

type HuaWei struct {
	AccessKey, SecretKey string // 密钥
}

// RefreshUrl 刷新URL
func (r *HuaWei) RefreshUrl(urls []string) error {
	auth, err := global.NewCredentialsBuilder().
		WithAk(r.AccessKey).
		WithSk(r.SecretKey).
		SafeBuild()
	if err != nil {
		return err
	}

	build, err := cdn.CdnClientBuilder().
		WithRegion(region.CN_NORTH_1).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return err
	}

	client := cdn.NewCdnClient(build)
	request := &model.CreateRefreshTasksRequest{}
	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().PREFIX
	modeRefreshTask := model.GetRefreshTaskRequestBodyModeEnum().ALL
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type: &typeRefreshTask,
		Mode: &modeRefreshTask,
		Urls: urls,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}

	response, err := client.CreateRefreshTasks(request)
	if err != nil {
		return err
	}

	if response.HttpStatusCode != 200 {
		return fmt.Errorf("刷新URL失败: %s", *response.RefreshTask)
	}

	return nil
}

// RefreshPath 刷新路径
func (r *HuaWei) RefreshPath(paths []string) error {
	auth, err := global.NewCredentialsBuilder().
		WithAk(r.AccessKey).
		WithSk(r.SecretKey).
		SafeBuild()
	if err != nil {
		return err
	}

	build, err := cdn.CdnClientBuilder().
		WithRegion(region.CN_NORTH_1).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return err
	}

	client := cdn.NewCdnClient(build)
	request := &model.CreateRefreshTasksRequest{}
	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().DIRECTORY
	modeRefreshTask := model.GetRefreshTaskRequestBodyModeEnum().ALL
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type: &typeRefreshTask,
		Mode: &modeRefreshTask,
		Urls: paths,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}

	response, err := client.CreateRefreshTasks(request)
	if err != nil {
		return err
	}

	if response.HttpStatusCode != 200 {
		return fmt.Errorf("刷新路径失败: %s", *response.RefreshTask)
	}

	return nil
}

// GetUsage 获取用量
func (r *HuaWei) GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error) {
	auth, err := global.NewCredentialsBuilder().
		WithAk(r.AccessKey).
		WithSk(r.SecretKey).
		SafeBuild()
	if err != nil {
		return 0, err
	}

	build, err := cdn.CdnClientBuilder().
		WithRegion(region.CN_NORTH_1).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return 0, err
	}

	client := cdn.NewCdnClient(build)
	request := &model.ShowDomainStatsRequest{}
	request.Action = "summary"
	request.StartTime = startTime.TimestampMilli()
	request.EndTime = endTime.TimestampMilli()
	request.DomainName = domain
	request.StatType = "req_num"
	response, err := client.ShowDomainStats(request)
	if err != nil {
		return 0, err
	}

	if response.HttpStatusCode != 200 {
		return 0, fmt.Errorf("获取用量失败: %v", response.Result)
	}

	if _, ok := response.Result["req_num"]; ok {
		return cast.ToUint(response.Result["req_num"]), nil
	}

	return 0, nil
}
