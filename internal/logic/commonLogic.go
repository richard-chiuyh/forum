package logic

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"forum/dao/db"
	"forum/dao/model"
	"forum/types/forum"
)

func PackageBaseResponse(code int32, msg string) *forum.BaseResponse {
	return &forum.BaseResponse{
		Code:    code,
		Message: msg,
	}
}

func SuccessResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.Success,
		// Message: svcErr.SuccessTxt,
	}
}

func NoPermissionResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorPermission,
		// Message: svcErr.ErrorPermissionTxt,
	}
}

func ParamErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorParam,
		// Message: svcErr.ErrorParamTxt,
	}
}

func MysqlErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorMysql,
		// Message: svcErr.ErrorMysqlTxt,
	}
}

func NoRecordErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorNoRecord,
		// Message: svcErr.ErrorNoRecordTxt,
	}
}

func RedisErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorRedis,
		// Message: svcErr.ErrorRedisTxt,
	}
}

func RpcErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorRpcService,
		// Message: svcErr.ErrorRpcServiceTxt,
	}
}

func KafkaErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorKafka,
		// Message: svcErr.ErrorKafkaTxt,
	}
}

func ElasticsearchErrorResponse() *forum.BaseResponse {
	return &forum.BaseResponse{
		// Code:    svcErr.SvcErrorElasticsearch,
		// Message: svcErr.ErrorElasticsearchTxt,
	}
}

func PaginationCheck(p *forum.Pagination) *forum.Pagination {
	if p == nil {
		return nil
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return p
}

func Offset(p *forum.Pagination) int {
	if p == nil {
		return -1
	}
	return int((p.Page - 1) * p.PageSize)
}

func Limit(p *forum.Pagination) int {
	if p == nil {
		return -1
	}
	return int(p.PageSize)
}

func PaginationResp(p *forum.Pagination, total int64) *forum.Pagination {
	if p == nil {
		return &forum.Pagination{
			Page:     1,
			PageSize: total,
			Total:    total,
		}
	}
	return &forum.Pagination{
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    total,
	}
}

// ValidatePostContent 校验单条内联内容（内容或图片二选一必填，内容≤300字，图片≤5张）
func ValidatePostContent(post *forum.Post) error {
	if strings.TrimSpace(post.Content) == "" && len(post.Images) == 0 {
		return errors.New("content or images is required")
	}
	if post.Content != "" {
		wordCount := utf8.RuneCountInString(post.Content)
		if wordCount > 300 {
			return errors.New("content exceeds maximum length of 300 characters")
		}
		post.WordCount = int32(wordCount)
	}
	if len(post.Images) > 5 {
		return errors.New("images exceeds maximum length of 5")
	}
	return nil
}

// func SendWsPostChangeMessage(ctx context.Context, kafka broker.Broker, post *model.Post) {
// 	data := map[string]interface{}{
// 		"id":          post.ID,
// 		"userId":      post.UserID,
// 		"account":     post.Account,
// 		"status":      post.Status,
// 		"publishTime": post.PublishTime,
// 	}
// 	var wsMsg = &mqMsg.WsMessage{
// 		Cmd:      app.WsforumPostChange,
// 		UserType: []int32{int32(users.UserType_UserTypeMember)},
// 		Data:     data,
// 	}
// 	bytes, _ := json.Marshal(wsMsg)
// 	go func() {
// 		if err := kafka.Publish(mqMsg.TopicWsSend, &broker.Message{
// 			Body: bytes,
// 		}); err != nil {
// 			logx.WithContext(ctx).Errorf("SendWsMessage send to kafka  error! in(%+v) err(%+v)", data, err)
// 		}
// 	}()
// }

// func SendWsforumCommentChangeMessage(ctx context.Context, kafka broker.Broker, comment *model.Post) {
// 	data := map[string]interface{}{
// 		"id":       strconv.FormatInt(comment.ID, 10),
// 		"userId":   comment.UserID,
// 		"account":  comment.Account,
// 		"parentId": strconv.FormatInt(comment.ParentID, 10),
// 		"status":   comment.Status,
// 	}
// 	var wsMsg = &mqMsg.WsMessage{
// 		Cmd:      app.WsforumCommentChange,
// 		UserType: []int32{int32(users.UserType_UserTypeMember)},
// 		Data:     data,
// 	}
// 	bytes, _ := json.Marshal(wsMsg)
// 	go func() {
// 		if err := kafka.Publish(mqMsg.TopicWsSend, &broker.Message{
// 			Body: bytes,
// 		}); err != nil {
// 			logx.WithContext(ctx).Errorf("SendWsMessage send to kafka  error! in(%+v) err(%+v)", data, err)
// 		}
// 	}()
// }

// GetLikedMap 返回 operator 已点赞(正常状态)的目标内容集合 targetId -> true
func GetLikedMap(ctx context.Context, likeDao db.LikeDao, targetIds []int64, operatorId int64) (map[int64]bool, error) {
	if len(targetIds) == 0 || operatorId == 0 {
		return map[int64]bool{}, nil
	}
	cond := &model.Like{
		UserID: operatorId,
		Status: int32(forum.Status_STATUS_PUBLISHED),
	}
	likes, _, err := likeDao.List(ctx, cond, db.TargetIdsInScope(targetIds))
	if err != nil {
		return nil, err
	}
	likedMap := make(map[int64]bool, len(likes))
	for _, like := range likes {
		likedMap[like.TargetID] = true
	}
	return likedMap, nil
}

// GetFavMap 返回 operator 已收藏(正常状态)的目标内容集合 targetId -> true
func GetFavMap(ctx context.Context, favDao db.FavDao, targetIds []int64, operatorId int64) (map[int64]bool, error) {
	if len(targetIds) == 0 || operatorId == 0 {
		return map[int64]bool{}, nil
	}
	cond := &model.Fav{
		UserID: operatorId,
		Status: int32(forum.Status_STATUS_PUBLISHED),
	}
	favs, _, err := favDao.List(ctx, cond, db.TargetIdsInScope(targetIds))
	if err != nil {
		return nil, err
	}
	favMap := make(map[int64]bool, len(favs))
	for _, fav := range favs {
		favMap[fav.TargetID] = true
	}
	return favMap, nil
}
