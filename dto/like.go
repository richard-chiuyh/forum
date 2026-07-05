package dto

import (
	"forum/dao/model"
	"forum/types/forum"
)

func ToRpcLike(like *model.Like) *forum.Like {
	if like == nil {
		return nil
	}
	return &forum.Like{
		Id:         like.ID,
		TargetId:   like.TargetID,
		UserId:     like.UserID,
		Account:    like.Account,
		Status:     forum.Status(like.Status),
		IpAddress:  like.IPAddress,
		CreateTime: like.CreateTime,
	}
}

func ToRpcLikes(likes []*model.Like) []*forum.Like {
	res := make([]*forum.Like, 0, len(likes))
	for _, like := range likes {
		res = append(res, ToRpcLike(like))
	}
	return res
}
