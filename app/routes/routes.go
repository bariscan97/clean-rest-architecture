package routes

import (
    "github.com/bariscan97/clean-rest-architecture/app/middleware"
    "github.com/bariscan97/clean-rest-architecture/internal/handler/post"
    "github.com/bariscan97/clean-rest-architecture/internal/handler/user"
    "github.com/go-chi/chi"
)

type Router struct {
    Mux         *chi.Mux
    userHandler user.Handler
    postHandler post.Handler
}

func NewRouter(uHandler user.Handler, pHandler post.Handler) *Router {
    return &Router{
        Mux:         chi.NewRouter(),
        userHandler: uHandler,
        postHandler: pHandler,
    }
}

func (r *Router) RegisterRoutes() {
    tokenMaker := r.userHandler.TokenMaker

    r.Mux.Route("/api/v1", func(api chi.Router) {

        api.Route("/posts", func(pr chi.Router) {
            pr.With(middleware.GetAuthMiddlewareFunc(tokenMaker)).Post("/", r.postHandler.CreatePost)
			pr.Get("/", r.postHandler.ListPosts)
			pr.Route("/{id}", func(idr chi.Router) {
                idr.Get("/comments", r.postHandler.GetCommentByPostID)

                idr.Group(func(gr chi.Router) {
                    gr.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))
                    gr.Patch("/", r.postHandler.UpdatePost)
                    gr.Delete("/", r.postHandler.DeletePostByID)
                })
            })
        })

        api.Route("/user", func(u chi.Router) {
            u.Patch("/", r.userHandler.UpdateUser)
            u.Delete("/", r.userHandler.DeleteUser)
            u.Get("/{id}", r.userHandler.GetUserByID)
        })

        api.Route("/auth", func(a chi.Router) {
            a.Post("/register", r.userHandler.CreateUser)
            a.Post("/login", r.userHandler.LoginUser)
        })
    })
}


