package route

import (
	"cosplayrent/internal/delivery/http"
	"cosplayrent/internal/delivery/http/middleware"
	"github.com/julienschmidt/httprouter"
)

type RouteConfig struct {
	Router               *httprouter.Router
	UserController       *controller.UserController
	CostumeController    *controller.CostumeController
	OrderController      *controller.OrderController
	TopUpOrderController *controller.TopUpOrderController
	MidtransController   *controller.MidtransController
	AuthMiddleware       *middleware.AuthMiddleware
}

func (c *RouteConfig) SetupRoute() {
	c.Router.POST("/api/register", c.UserController.Register)
	c.Router.GET("/api/identitycard", c.AuthMiddleware.ServeHTTP(c.UserController.GetIdentityCard))
	c.Router.PUT("/api/identitycard", c.AuthMiddleware.ServeHTTP(c.UserController.AddOrUpdateIdentityCard))
	c.Router.GET("/api/emoney", c.AuthMiddleware.ServeHTTP(c.UserController.GetEMoneyAmount))
	c.Router.GET("/api/emoneyhistory", c.AuthMiddleware.ServeHTTP(c.UserController.GetEMoneyTransactionHistory))
	c.Router.POST("/api/login", c.UserController.Login)
	c.Router.GET("/api/userdetail", c.AuthMiddleware.ServeHTTP(c.UserController.FindByUUID))
	c.Router.GET("/api/user", c.AuthMiddleware.ServeHTTP(c.UserController.FindAll))
	c.Router.PATCH("/api/userdetail", c.AuthMiddleware.ServeHTTP(c.UserController.Update))
	//c.Router.DELETE("/api/useraccount", c.AuthMiddleware.ServeHTTP(c.UserController.Delete))
	c.Router.GET("/api/checkuserstatus/:costumeID", c.AuthMiddleware.ServeHTTP(c.UserController.CheckUserStatus))
	c.Router.GET("/api/selleraddress/checkout/:costumeID", c.AuthMiddleware.ServeHTTP(c.UserController.FindSellerAddressDetailByCostumeId))
	//c.Router.GET("/api/checkappversion", c.UserController.CheckAppVersion)

	c.Router.POST("/api/costume", c.AuthMiddleware.ServeHTTP(c.CostumeController.Create))
	c.Router.GET("/api/costume", c.CostumeController.FindAll)
	c.Router.GET("/api/seller", c.AuthMiddleware.ServeHTTP(c.CostumeController.FindSellerCostume))
	c.Router.GET("/api/costume/:costumeID", c.CostumeController.FindById)
	c.Router.GET("/api/syou are eller/:costumeID", c.AuthMiddleware.ServeHTTP(c.CostumeController.FindSellerCostumeByCostumeID)) // find by costume id
	c.Router.PATCH("/api/seller/:costumeID", c.AuthMiddleware.ServeHTTP(c.CostumeController.Update))
	c.Router.DELETE("/api/seller/:costumeID", c.AuthMiddleware.ServeHTTP(c.CostumeController.Delete))

	//c.Router.GET("/api/review", c.AuthMiddleware.ServeHTTP(c.ReviewController.FindUserReview))
	//c.Router.POST("/api/review", c.AuthMiddleware.ServeHTTP(c.ReviewController.Create))
	//c.Router.GET("/api/review/:reviewID", c.AuthMiddleware.ServeHTTP(c.ReviewController.FindUserReviewByReviewID))
	//c.Router.PUT("/api/review/:reviewID", c.AuthMiddleware.ServeHTTP(c.ReviewController.Update))
	//c.Router.DELETE("/api/review/:reviewID", c.AuthMiddleware.ServeHTTP(c.ReviewController.DeleteUserReviewByReviewID))
	//c.Router.GET("/api/costume/:costumeID/review", c.ReviewController.FindByCostumeId)

	c.Router.POST("/api/order/midtrans", c.AuthMiddleware.ServeHTTP(c.OrderController.Create))
	c.Router.GET("/api/checkorder/:orderID", c.OrderController.CheckStatusPayment)
	c.Router.GET("/api/order/seller", c.AuthMiddleware.ServeHTTP(c.OrderController.GetAllSellerOrder))
	//c.Router.PUT("/api/order/:orderID", c.AuthMiddleware.ServeHTTP(c.OrderController.UpdateSellerOrder))
	c.Router.GET("/api/orderdetail/:orderID", c.AuthMiddleware.ServeHTTP(c.OrderController.GetDetailOrderByOrderId))
	c.Router.GET("/api/userorder/:orderID", c.AuthMiddleware.ServeHTTP(c.OrderController.GetUserDetailOrder))
	c.Router.GET("/api/alluserorder", c.AuthMiddleware.ServeHTTP(c.OrderController.GetAllUserOrder))
	c.Router.POST("/api/checkbalancewithorderamount", c.AuthMiddleware.ServeHTTP(c.OrderController.CheckBalanceWithOrderAmount))

	c.Router.PUT("/api/topup", c.AuthMiddleware.ServeHTTP(c.TopUpOrderController.CreateTopUpOrder))
	//c.Router.GET("/api/checktopuporder/:orderID", c.TopUpOrderController.CheckTopUpOrderByOrderId)

	//c.Router.GET("/api/wishlist", c.AuthMiddleware.ServeHTTP(c.WishlistController.FindAllWishListByUserId))
	//c.Router.POST("/api/wishlist/:costumeID", c.AuthMiddleware.ServeHTTP(c.WishlistController.AddWishlist))
	//c.Router.DELETE("/api/wishlist/:costumeID", c.AuthMiddleware.ServeHTTP(c.WishlistController.DeleteWishlist))
	//
	//c.Router.GET("/api/provinces", c.AuthMiddleware.ServeHTTP(c.RajaOngkirController.FindProvince))
	//c.Router.GET("/api/city/:provinceID", c.AuthMiddleware.ServeHTTP(c.RajaOngkirController.FindCity))
	//c.Router.POST("/api/checkshippment", c.AuthMiddleware.ServeHTTP(c.RajaOngkirController.CheckShippment))

	c.Router.POST("/api/midtrans/callback", c.MidtransController.MidtransCallBack)
}
