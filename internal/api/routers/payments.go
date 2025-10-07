package routers

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/service/payment"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/usecase"

	repositories "github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	"github.com/labstack/echo/v4"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
)

func PaymentsRouter(e *echo.Echo, dbConnection *dbadapter.DatabaseAdapter) {
	r := repositories.NewPaymentsRepository(dbConnection)
	or := repositories.NewOrderRepository(dbConnection)
	u := usecase.NewPaymentsUseCase(r, payment.NewMPService(), or)
	h := handlers.NewPaymentHandler(u)

	g := e.Group("/payments")

	g.GET("/:id", h.GetByID)
	g.POST("", h.CreatePayment)
	g.POST("/webhook", h.PaymentWebHook)
}
