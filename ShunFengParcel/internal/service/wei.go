package service

import (
	pb "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/biz"
)

type WeiService struct {
	pb.UnimplementedWeiServer
	uc *biz.GreeterUsecase
}

func NewWeiService(uc *biz.GreeterUsecase) *WeiService {
	return &WeiService{uc: uc}
}

//
//const TOKEN = "111"
//const appid = "wxe6b3abe33cc7c8c1"
//const appsecret = "a41d654ea882bd0d736af242931c9d4e"
//
//func (s *WeiService) WX(ctx context.Context, req *pb.WXRequest) (*pb.WXReply, error) {
//	// 微信验证逻辑应该通过查询参数传递，但在 Kratos 中，这些参数会自动绑定到请求结构体
//	// 这里我们需要重新设计请求结构体来包含验证所需的参数
//	// 暂时返回成功响应
//	signature := ctx.Value("signature") // 微信已经通过token加密好的秘文
//	timestamp := ctx.Value("timestamp")
//	nonce := ctx.Value("nonce")
//	echostr := ctx.Value("echostr")
//	// 创建包含令牌、时间戳和随机数的字符串切片
//	tmpArr := []string{TOKEN, timestamp, nonce}
//	// 对切片进行字典排序
//	sort.Strings(tmpArr)
//	// 将排序后的元素拼接成单个字符串
//	tmpStr := ""
//	for _, v := range tmpArr {
//		tmpStr += v
//	}
//	// 对字符串进行SHA-1哈希计算
//	tmpHash := sha1.New()
//
//	tmpHash.Write([]byte(tmpStr))
//	tmpStr = fmt.Sprintf("%x", tmpHash.Sum(nil))
//	//fmt.Println(tmpStr)
//	//fmt.Println(signature)
//	// 将计算得到的签名与请求中提供的签名进行比较，并根据结果发送相应的响应
//	if tmpStr == signature {
//		ctx2.String(200, echostr)
//	} else {
//		ctx2.String(403, "签名验证失败 "+timestamp)
//	}
//	return &pb.WXReply{}, nil
//}
//
//// state参数 生成随机数的函数
//func GenerateRandomCode(length int) string {
//	source := rand.NewSource(time.Now().UnixNano())
//	generator := rand.New(source)
//	letters := []rune("1234567890")
//	code := make([]rune, length)
//	for i := range code {
//		index := generator.Intn(len(letters))
//		code[i] = letters[index]
//	}
//	return string(code)
//}
//
//func (s *WeiService) WechatLogin(ctx context.Context, req *pb.WechatLoginRequest) (*pb.WechatLoginReply, error) {
//	// 构建微信 OAuth 授权 URL
//	redirect_uri := "https://4f1ad8c0.r9.cpolar.cn/callback"
//	// 对 redirect_uri 进行 URL 编码
//	encodedRedirectURI := url.QueryEscape(redirect_uri)
//	state := GenerateRandomCode(5) // 防止跨站请求伪造攻击 增加安全性
//	scope := "snsapi_userinfo"
//	codeUrl := fmt.Sprintf(
//		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
//		appid,
//		encodedRedirectURI, // 使用编码后的 redirect_uri
//		scope,
//		state,
//	)
//
//	// 生成二维码
//	_, err := qrcode.Encode(codeUrl, qrcode.Medium, 256)
//	if err != nil {
//		return &pb.WechatLoginReply{}, err
//	}
//
//	// 在 Kratos 中，应该返回二维码数据，让 HTTP 层处理响应格式
//	// 这里我们可以将二维码数据编码为 base64 字符串返回
//	return &pb.WechatLoginReply{}, nil
//}
//
//func (s *WeiService) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.CallbackReply, error) {
//	// 微信 OAuth 回调处理
//	// 在 Kratos 中，查询参数应该通过请求结构体传递
//	// 这里需要重新设计 CallbackRequest 结构体来包含 code 和 state 参数
//	// 暂时返回成功响应
//	return &pb.CallbackReply{}, nil
//}
