package tokenAuth

import (
	"crypto/md5"
	"fmt"
	"go-grpc-example/proto/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

const IsTLS = true

type TokenAuth struct {
	token.TokenValidateParam
}

func (x *TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"uid":   strconv.FormatInt(int64(x.GetUid()), 10),
		"token": x.GetToken(),
	}, nil
}

func (x *TokenAuth) RequireTransportSecurity() bool {
	return IsTLS
}
func (x *TokenAuth) CheckToken(ctx context.Context) (*token.Response, error) {
	// 验证token
	md, b := metadata.FromIncomingContext(ctx)
	if !b {
		return nil, status.Error(codes.InvalidArgument, "token信息不存在")
	}

	var token, uid string
	// 取出token
	tokenInfo, ok := md["token"]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "token不存在")
	}

	token = tokenInfo[0]

	// 取出uid
	uidTmp, ok := md["uid"]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "uid不存在")
	}
	uid = uidTmp[0]

	//验证
	sum := md5.Sum([]byte(uid))
	md5Str := fmt.Sprintf("%x", sum)
	if md5Str != token {
		fmt.Println("md5Str:", md5Str)
		fmt.Println("uid:", uid)
		fmt.Println("token:", token)
		return nil, status.Error(codes.InvalidArgument, "token验证失败")
	}
	return nil, nil
}
