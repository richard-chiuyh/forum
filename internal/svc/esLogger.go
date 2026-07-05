package svc

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type EsLogger struct{}

func (l *EsLogger) LogRoundTrip(
	req *http.Request,
	res *http.Response,
	err error,
	start time.Time,
	dur time.Duration,
) error {
	ctx := req.Context()
	var reqBody string

	if req.Body != nil {
		buf, _ := io.ReadAll(req.Body)
		reqBody = string(buf)
		req.Body = io.NopCloser(bytes.NewBuffer(buf))
	}

	fields := []logx.LogField{
		logx.Field("method", req.Method),
		logx.Field("url", req.URL.String()),
		logx.Field("duration", dur.String()),
		logx.Field("body", reqBody),
	}

	if err != nil {
		logx.WithContext(ctx).WithFields(fields...).Errorf("ES_REQ_FAILED: %v", err)
		return nil
	}

	if res != nil {
		fields = append(fields, logx.Field("status", res.StatusCode))

		// 如果你还需要输出 Response Body（可选，生产环境建议关闭）
		// respBuf, _ := io.ReadAll(res.Body)
		// res.Body = io.NopCloser(bytes.NewBuffer(respBuf))
		// fields = append(fields, logx.Field("response_body", string(respBuf)))
	}

	logx.WithContext(ctx).WithFields(fields...).Info("ES_FULL_REQUEST")
	return nil
}

func (l *EsLogger) RequestBodyEnabled() bool {
	return true
}

func (l *EsLogger) ResponseBodyEnabled() bool {
	return true
}
