package gfutils

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNum(t *testing.T) {

	price, err := decimal.NewFromString("136.02")
	if err != nil {
		panic(err)
	}

	quantity := decimal.NewFromInt(3)

	fee, _ := decimal.NewFromString(".035")

	taxRate, _ := decimal.NewFromString(".08875")

	// price * quantity
	subtotal := price.Mul(quantity)

	// subtotal * (fee + 1)
	preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1)))

	// preTax * (taxRate + 1)
	total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1)))

	fmt.Println("subtotal:", subtotal)                      // Subtotal: 408.06
	fmt.Println("Pre-tax:", preTax)                         // Pre-tax: 422.3421
	fmt.Println("Taxes:", total.Sub(preTax))                // Taxes: 37.482861375
	fmt.Println("Total:", total)                            // Total: 459.824961375
	fmt.Println("Tax rate:", total.Sub(preTax).Div(preTax)) // Tax rate: 0.08875

	// 设置精度
	decimal.DivisionPrecision = 3
	num, _ := decimal.NewFromString("1")
	// 验证除法
	fmt.Println(num.Div(decimal.NewFromFloatWithExponent(3, 0)).String())

}
