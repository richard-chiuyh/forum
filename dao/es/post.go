package es

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"forum/dao/model"
	"forum/types/forum"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type postOpr struct {
	esClient  *elasticsearch.TypedClient
	tableName string
}

func newPostOpr(esClient *elasticsearch.TypedClient, env string) PostDao {
	return &postOpr{
		esClient:  esClient,
		tableName: fmt.Sprintf("%s_%s", model.TableNamePost, env),
	}
}

func initPostIndex(esClient *elasticsearch.TypedClient, env string) error {
	tableName := fmt.Sprintf("%s_%s", model.TableNamePost, env)
	exist, err := esClient.Indices.Exists(tableName).IsSuccess(context.Background())
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	res, err := esClient.Indices.Create(tableName).Mappings(model.Post{}.Mapping()).Do(context.Background())
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("es create index error: %v", res)
	}
	return nil
}

func (p *postOpr) Index(ctx context.Context, post *model.Post) (err error) {
	_, err = p.esClient.Index(p.tableName).Document(post).Id(strconv.FormatInt(post.ID, 10)).Do(ctx)
	return
}

func (p *postOpr) Delete(ctx context.Context, id int64) (err error) {
	_, err = p.esClient.Delete(p.tableName, strconv.FormatInt(id, 10)).Do(ctx)
	return
}

func (p *postOpr) Search(ctx context.Context, in *forum.ListPostRequest) ([]*model.Post, int64, error) {
	var filters, must []types.QueryVariant
	filters = append(filters, esdsl.NewTermQuery("staparent_idtus", esdsl.NewFieldValue().Int64(0)))
	if in.UserId != 0 {
		filters = append(filters, esdsl.NewTermQuery("user_id", esdsl.NewFieldValue().Int64(in.UserId)))
	}
	if in.PostFilter != nil {
		if in.PostFilter.Content != "" {
			must = append(must, esdsl.NewMatchQuery("content", in.PostFilter.Content))
		}
		if in.PostFilter.PostId != 0 {
			filters = append(filters, esdsl.NewTermQuery("id", esdsl.NewFieldValue().Int64(in.PostFilter.PostId)))
		}
		if len(in.PostFilter.Status) > 0 {
			vals := make(map[string]types.TermsQueryField, len(in.PostFilter.Status))
			for _, s := range in.PostFilter.Status {
				vals["status"] = esdsl.NewFieldValue().Int64(int64(s))
			}
			filters = append(filters, esdsl.NewTermsQuery().TermsQuery(vals))
		}
		if in.PostFilter.Account != "" {
			filters = append(filters, esdsl.NewWildcardQuery("account", "*"+in.PostFilter.Account+"*"))
		}

		if in.PostFilter.PublishTimeStart != 0 || in.PostFilter.PublishTimeEnd != 0 {
			rangeQuery := esdsl.NewLongNumberRangeQuery("publish_time")
			if in.PostFilter.PublishTimeStart != 0 {
				rangeQuery = rangeQuery.Gte(in.PostFilter.PublishTimeStart)
			}
			if in.PostFilter.PublishTimeEnd != 0 {
				rangeQuery = rangeQuery.Lte(in.PostFilter.PublishTimeEnd)
			}
			filters = append(filters, rangeQuery)
		}
	}
	req := p.esClient.Search().Index(p.tableName).Query(esdsl.NewBoolQuery().Filter(filters...).Must(must...))
	if in.Pagination != nil {
		req = req.From(int((in.Pagination.Page - 1) * in.Pagination.PageSize)).Size(int(in.Pagination.PageSize))
	}
	res, err := req.Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	posts := make([]*model.Post, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var post model.Post
		if err := json.Unmarshal(hit.Source_, &post); err != nil {
			return nil, 0, err
		}
		posts = append(posts, &post)
	}
	return posts, res.Hits.Total.Value, nil
}
