package controllers

import (
	"backend/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"net/http"
	"strconv"
)

func GetPaymentController(e *echo.Group) {
	g := e.Group("/payment")
	g.POST("/:id", func(c echo.Context) error {
		id := c.Param("id")
		u64, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		n := uint(u64)
		res, err := GeneratePaymentLink(n)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to handle payment"})
		}
		return c.JSON(http.StatusOK, map[string]string{"clientSecret": res})
	})
}

func CalculateOrderCost(orderID uint) (int64, error) {
	order, err := GetOrderById(orderID)
	if err != nil {
		return -1, err
	}
	fmt.Print(order, "xd")

	amount := order.TotalPrice + 220.00

	return int64(amount), nil
}

func GeneratePaymentLink(orderID uint) (string, error) {
	amount, err := CalculateOrderCost(orderID)
	if err != nil {
		return "", err
	}

	stripe.Key = utils.GetValueFromEnv("STRIPE_SECRET", "secret")

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyPLN)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}
