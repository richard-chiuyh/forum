package forumlogic

import (
	"context"

	"forum/dto"
	"forum/internal/logic"
	"forum/internal/svc"
	"forum/types/forum"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *forum.GetPostRequest) (*forum.GetPostResponse, error) {
	post, err := l.svcCtx.DBManager.PostDao.Get(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("GetPost Get err(%+v)", err)
		return &forum.GetPostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}
	return &forum.GetPostResponse{BaseResp: logic.SuccessResponse(), Post: dto.ToRpcPost(post, false, false)}, nil
}
