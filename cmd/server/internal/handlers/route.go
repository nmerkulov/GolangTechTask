package handlers

import (
	"github.com/go-chi/chi"
)

func Routes(store Store) *chi.Mux {
	r := chi.NewRouter()

	bh := buffHandler{
		store: store,
	}
	sh := newStreamHandler(store)

	r.Route("/buff", func(r chi.Router) {
		r.Get("/{id}", bh.GetBuff)
		r.Post("/", bh.CreateBuff)
	})
	r.Route("/stream", func(r chi.Router) {
		r.Get("/{id}", sh.GetStream)
		r.Post("/", sh.CreateStream)
	})

	return r
}
