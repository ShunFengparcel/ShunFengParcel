package utils

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
)

func Alipayment(orderNo, price string) string {
	var privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC5+awnsflJyb+3KNUjgCT0rMZ7Q8vutr9A4Ex10qwv0hOaoiZP8PpAhe/mkf+SyonaymtMZiybiZ3oaXAZxvlIvNq6+O7uX/MLE+87/8UzI2aQqdblkZyym/ehbszPVMg/YjHfDvjq/RTyBYTrLva5gWTdlB7d6FPuBLLJP8ezEe+DSXzYbX0RJAiQHxs4Uf6u4cYQitv63zz1WaM0Jwhk8OMz+xf6+6Yrbc3bZ/y7MzcLyuTWL6t2oFxiQdhisFVv90jTfn2pwsLWCMvjEQrwoH+z1NH3bfdu7YO2PPSGTycuelgqpCIDDTMB8ykll0WHSDZVUunC501TetDpUP5PAgMBAAECggEBAKze52byqY4safggjZb1S+TgIZN9HrijyL3OgWRYU1QUS8LGjSRsFMMbfTdvaQkTGfd1hy26ICboUg8yy1A6w4gYfJ2mPbr5F5upiQzVoIH7myBHr4kvxF/GWPYHX3wuDAHXDhKRAVK1s92QwdA3lII1RAQv5k1R44pxdZtanQ/WuJHNTj4pktme+M7ULgNJmWXsflxNIjkc8tuzF9p6sgrkP6w8lSxR/UfJMHK1X1s2U7I5GgnPRGU7e765usPiA3CvtepDSDThQVhuRnbz3rUOxthRwos2SS6wr2pnqq1nj0K6G9aCYQyazoNjfPtnLpNWbI0pBpDm8rFJbOC05nECgYEA9lEz2GA2vPAFlQApt7UrjznfpfuIef5vQbUUrHaa7svUx0KhxXLkYSMmTF/DQFgVd/kXEsg3vAPdx/ZSsnV/HdL7kitixB1eAITM1xb1N0/UvbZ/NIgKxrXkXc/8JEkzDdR9rWa1WnYbA+wgp0hP8KR2C0KxBFwl+efbE01EFecCgYEAwUk5EmZa2BGuMYuGccXgO6lYHrnAA3vTRPPWkBahKffAuCYPvyySjvbhO8Abr5sDKKdQPvL9Ck++KLc7tCAps8ExoGJ0Tq3TPaXVE3CEpFPKefqCCKx7Mkgt5Y/S+BjksW59C4Q8z0dW95wek2ap7iZxxJa7yI6dtXevIPjxd1kCgYA/GbxTYQqEymRTsG//fO0EywmtRsvGnNS5m38JU3ULWbJPvZUdtPomnE+SXzHwyN/vFSPBDwOgKclmEYdL6me/Jy6FWpH5taBAN8UWEO6O1eelFhxuQ6+nCi/PjJmGXi4zC82KX3Z0Dy+KiLIwyIiaGDeZWONqP4UHCUuJHVEk3wKBgE/IldJOhbkiszCoUzqrXz/BSyqDqgrGFhMkQ7D+ZlAYgGiC7YUQNP7mUVqElekKp2ckiS8yxdh3yhqsZ+yWSiB04rM9cJz5i3Sq+yUnENlz7OQkz4AdEk1TFf7oO0FFpUDIRr12PFOjMvbKbqSRgBtZqyRmw+SpWdgKKzQFDkchAoGBALZYSAe/COTJ2d5WZW6+HGMTNkuzG7FDn0dkwfTcBcdaxfD/i2d3lRpe5S3HKUMTOTfX2bbaJTZYMG6Vx3YXWa92zhnnI0fXRIboYjye1bMQ3FW+AZWhAEFxYuJOC8G1IDcHSf0nu+1E1AXKvtLGkqhwWJGjM3n8hIww5HrOY5uD" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	client, err := alipay.New("2021000148652076", privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 加载支付宝公钥用于验证回调签名
	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuTeRNKVmgIl3iMzw3X7y76e5M0PWKT8LmfRXbOSC5/JMIV27MIy9Zns/ym98jW4dOI4W+0PO47k4JKdeZzzJBjOBVcXM1lru4m8aMclRqVGLhKWDVSn2CiL8NbOf2fDMflDYlW8eCrBBUTnoTGkhYijTVPL6av3GNCB5WhtGUjyRinMXBo43Xdujgz8SqU2V5Y+tRZoZeMu8Hz68OBCO+yLwT59JDnASsUaHZ6Axk71ZOkINruWMGoxMKL87/Q7+Zggh1tAcVm+UPnaIGf9DTkddBg45+mJL2KOC/J1b7FQqE8VfftR1Ybd88/fxLVF/1hxkvWAMyrpcSXkt/dP4qwIDAQAB"
	if err := client.LoadAliPayPublicKey(publicKey); err != nil {
		fmt.Println("加载支付宝公钥失败:", err)
		return ""
	}

	var p = alipay.TradeWapPay{}
	p.NotifyURL = "http://512140e1.r11.vip.cpolar.cn/payment/update"
	p.ReturnURL = "http://www.baidu.com"
	p.Subject = "顺丰速递"
	p.OutTradeNo = orderNo
	p.TotalAmount = price
	p.ProductCode = "QUICK_WAP_WAY"

	url, err := client.TradeWapPay(p)
	if err != nil {
		fmt.Println(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
