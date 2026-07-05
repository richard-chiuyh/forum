package dto

import (
	"forum/dao/model"
	"forum/types/forum"
)

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func int64Val(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func int32Val(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func FromRpcPost(post *forum.Post) *model.Post {
	return &model.Post{
		ID:          post.Id,
		ParentID:    post.ParentId,
		UserID:      post.UserId,
		Account:     post.Account,
		Status:      int32(post.Status),
		Remark:      &post.Remark,
		Content:     post.Content,
		Images:      post.Images,
		WordCount:   post.WordCount,
		ViewCount:   post.ViewCount,
		IPAddress:   post.IpAddress,
		PublishTime: &post.PublishTime,
	}
}

func ToRpcPost(post *model.Post, isLiked bool, isFav bool) *forum.Post {
	if post == nil {
		return nil
	}
	rpc := &forum.Post{
		Id:           post.ID,
		UserId:       post.UserID,
		Account:      post.Account,
		Status:       forum.Status(post.Status),
		Remark:       strVal(post.Remark),
		Content:      post.Content,
		Images:       post.Images,
		WordCount:    post.WordCount,
		ViewCount:    post.ViewCount,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		FavCount:     post.FavCount,
		PublishTime:  int64Val(post.PublishTime),
		IsLiked:      isLiked,
		IsFav:        isFav,
		ParentId:     post.ParentID,
		IpAddress:    post.IPAddress,
	}
	if post.Parent != nil {
		rpc.Parent = ToRpcPost(post.Parent, false, false)
	}
	return rpc
}

func ToRpcPosts(posts []*model.Post, likedMap map[int64]bool, favMap map[int64]bool) []*forum.Post {
	res := make([]*forum.Post, 0, len(posts))
	for _, post := range posts {
		res = append(res, ToRpcPost(post, likedMap[post.ID], favMap[post.ID]))
	}
	return res
}
