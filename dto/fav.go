package dto

import (
	"forum/dao/model"
	"forum/types/forum"
)

// ToRpcFav 收藏映射。target 为内容统一模型：parent_id==0 时是帖子，否则是评论。
// isTargetLiked 表示当前操作人是否点赞了该目标内容。
func ToRpcFav(fav *model.Fav, isTargetLiked bool) *forum.Fav {
	if fav == nil {
		return nil
	}
	rpc := &forum.Fav{
		Id:         fav.ID,
		TargetId:   fav.TargetID,
		UserId:     fav.UserID,
		Account:    fav.Account,
		Status:     forum.Status(fav.Status),
		CreateTime: fav.CreateTime,
	}
	if fav.Target != nil {
		// 目标内容统一用 Post 表示（帖子或评论）
		rpc.Post = ToRpcPost(fav.Target, isTargetLiked, true)
	}
	return rpc
}

func ToRpcFavs(favs []*model.Fav, likedMap map[int64]bool) []*forum.Fav {
	res := make([]*forum.Fav, 0, len(favs))
	for _, fav := range favs {
		res = append(res, ToRpcFav(fav, likedMap[fav.TargetID]))
	}
	return res
}
