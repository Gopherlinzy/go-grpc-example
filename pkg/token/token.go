package tokenAuth

import (
	"go-grpc-example/proto/token"
	"golang.org/x/net/context"
	"strconv"
)

const IsTLS = true

// 定义一个认证的结构体，这里是因为我在porto写好了一个数据结构
// 也可以自定义认证字段
type TokenAuth struct {
	token.TokenValidateParam
}

func (x *TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	// 将 Credentials（认证凭证）存放在 metadata（元数据）中进行传递。
	return map[string]string{
		"uid":   strconv.FormatInt(int64(x.GetUid()), 10),
		"token": x.GetToken(),
	}, nil
}

func (x *TokenAuth) RequireTransportSecurity() bool {
	return IsTLS
}
