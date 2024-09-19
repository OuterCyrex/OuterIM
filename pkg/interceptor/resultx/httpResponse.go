package resultx

import (
	"OuterIM/pkg/xerr"
	"context"
	"errors"
	pkgErr "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, data any) any {
	return Success(data)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errCode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errCode)

		causeErr := pkgErr.Cause(err)
		var e *zrpcErr.CodeMsg
		if errors.As(causeErr, &e) {
			errCode = e.Code
			errmsg = e.Msg
		} else {
			if gStatus, ok := status.FromError(err); ok {
				errCode = int(gStatus.Code())
				errmsg = gStatus.Message()
			}
		}

		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)

		return http.StatusBadRequest, Fail(errCode, errmsg)
	}
}
