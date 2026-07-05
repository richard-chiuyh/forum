package forumlogic

import (
	"context"
	"errors"
	"time"

	"forum/dao/db"
	"forum/dao/model"
	"forum/dto"
	"forum/internal/logic"
	"forum/internal/svc"
	"forum/types/forum"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreatePost 创建内容：parent_id = 0 为帖子，parent_id > 0 为评论/回复
func (l *CreatePostLogic) CreatePost(in *forum.CreatePostRequest) (*forum.CreatePostResponse, error) {
	if in.Post == nil || in.Post.UserId == 0 {
		return &forum.CreatePostResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}
	if err := logic.ValidatePostContent(in.Post); err != nil {
		l.Logger.Errorf("CreatePost validatePostContent err(%+v)", err)
		return &forum.CreatePostResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}

	contentModel := dto.FromRpcPost(in.Post)
	if contentModel.ParentID == 0 {
		return l.createRootPost(contentModel)
	}
	return l.createComment(contentModel)
}

// createRootPost 创建帖子（直接发布，写 ES，更新用户发帖数）
func (l *CreatePostLogic) createRootPost(post *model.Post) (*forum.CreatePostResponse, error) {
	post.Status = int32(forum.Status_STATUS_PUBLISHED)
	publishTime := time.Now().UnixMilli()
	post.PublishTime = &publishTime

	if err := l.svcCtx.DBManager.PostDao.Create(l.ctx, post); err != nil {
		l.Logger.Errorf("CreatePost createRootPost err(%+v)", err)
		return &forum.CreatePostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}
	if err := l.svcCtx.ESManager.PostDao.Index(l.ctx, post); err != nil {
		l.Logger.Errorf("CreatePost Index err(%+v)", err)
		return &forum.CreatePostResponse{BaseResp: logic.ElasticsearchErrorResponse()}, err
	}
	if err := l.updateCommUserInfo(post); err != nil {
		l.Logger.Errorf("CreatePost updateCommUserInfo err(%+v)", err)
	}
	// logic.SendWsPostChangeMessage(l.ctx, l.svcCtx.Kafka, post)
	return &forum.CreatePostResponse{BaseResp: logic.SuccessResponse(), Post: dto.ToRpcPost(post, false, false)}, nil
}

// createComment 创建评论/回复：校验父节点存在且层级合法（树最多 2 层），父节点 comment_count + 1
func (l *CreatePostLogic) createComment(comment *model.Post) (*forum.CreatePostResponse, error) {
	parent, err := l.svcCtx.DBManager.PostDao.Get(l.ctx, comment.ParentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &forum.CreatePostResponse{BaseResp: logic.NoRecordErrorResponse()}, nil
		}
		l.Logger.Errorf("CreatePost get parent err(%+v)", err)
		return &forum.CreatePostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}
	if parent.ParentID != 0 {
		l.Logger.Errorf("CreatePost parent is a reply, not allowed! parent(%+v)", parent)
		return &forum.CreatePostResponse{BaseResp: logic.NoPermissionResponse()}, nil
	}

	comment.ID = 0 // 统一表使用 DB 自增主键
	comment.Status = int32(forum.Status_STATUS_PUBLISHED)
	if err := l.svcCtx.DBManager.PostDao.Create(l.ctx, comment); err != nil {
		l.Logger.Errorf("CreatePost createComment err(%+v)", err)
		return &forum.CreatePostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	l.svcCtx.DBManager.Executor.Add(db.CountTask{
		Table:  model.Post{}.TableName(),
		Id:     comment.ParentID,
		Column: db.COMMENT_COUNT,
		Delta:  1,
	})

	// logic.SendWsforumCommentChangeMessage(l.ctx, l.svcCtx.Kafka, comment)
	return &forum.CreatePostResponse{BaseResp: logic.SuccessResponse(), Post: dto.ToRpcPost(comment, false, false)}, nil
}

// updateCommUserInfo 更新社区用户信息（发帖数）
func (l *CreatePostLogic) updateCommUserInfo(post *model.Post) error {
	userInfoModel, err := l.svcCtx.DBManager.UserInfoDao.GetByUserID(l.ctx, post.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) || userInfoModel == nil {
		return l.svcCtx.DBManager.UserInfoDao.Create(l.ctx, &model.UserInfo{
			UserID:    post.UserID,
			Account:   post.Account,
			PostCount: 1,
		})
	}
	if err != nil {
		return err
	}
	userInfoModel.PostCount++
	return l.svcCtx.DBManager.UserInfoDao.Update(l.ctx, userInfoModel)
}
