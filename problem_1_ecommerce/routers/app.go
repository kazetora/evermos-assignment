package routers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/controllers"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/database"
)

func apiRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	orderController := controllers.NewOrderController(db)

	r.Route("/order", func(r chi.Router) {
		r.Post("/addToCart", orderController.AddToCart)
	})

	r.Route("/transaction", func(r chi.Router) {
		r.Get("/{id}", controllers.GetTransactionStatus)
	})

	return r
}

// InitRouter function to init API router
func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	db := database.GetDatabase()

	r.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	r.Mount("/api/v1", apiRouter(db))

	return r
}
