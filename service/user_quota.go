package service

import (
	"context"
	"time"

	"github.com/sndies/chat_with_u/consts"
	"github.com/sndies/chat_with_u/db/dao"
	"github.com/sndies/chat_with_u/db/db_model"
	"github.com/sndies/chat_with_u/middleware/log"
)

func IsUserHasRunOutOfQuota(ctx context.Context, openId string) (bool, error) {
	// 起始时间今天0点,结束时间明天0点
	var (
		now      = time.Now()
		tomorrow = now.Add(24 * time.Hour)
	)
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	// 查询该用户今天产生的问答数量
	qnaList, err := dao.GetGptQNAByOpenId(ctx, openId, startTime, endTime)
	if err != nil {
		return false, err
	}
	// 统计一下有答案的个数
	var qnaListHavingAnswer []*db_model.GptQNA
	for _, qna := range qnaList {
		q := qna
		if q.Answer != "" {
			qnaListHavingAnswer = append(qnaListHavingAnswer, q)
		}
	}
	// 阈值目前配置在const.go里
	if len(qnaListHavingAnswer) >= consts.FreeQuotaPerDay {
		log.Infof(ctx, "[IsUserHasRunOutOfQuota] openId: %s has runOutOf quota", openId)
		return true, nil
	}
	return false, nil
}
