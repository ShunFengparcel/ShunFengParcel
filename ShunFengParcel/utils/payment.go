package utils

import "fmt"

type Payment interface {
	Pay()
}

type AliPay struct {
}

func (AliPay) Pay(orderSn, price string) string {
	url := Alipayment(orderSn, price)
	fmt.Println("支付宝付款")
	return url

}

type WxPay struct {
}

func (WxPay) Pay() {
	fmt.Println("微信付款")
}
