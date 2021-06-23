package router

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idalmasso/storytellers/backend/api/types"
	"github.com/idalmasso/storytellers/backend/common"
)
type mapKeyString string
const(
	STORY_ID mapKeyString ="story"
)
type StoryServerRouter struct{
	Router chi.Router
	app types.StoryTellerAppInterface
}

func CreateRouter(app types.StoryTellerAppInterface) StoryServerRouter{
	router:=StoryServerRouter{app: app}

	router.init()
	return router
}

func (r *StoryServerRouter) init(){
	r.Router=chi.NewRouter()
	r.Router.Use(middleware.RequestID)
  r.Router.Use(middleware.RealIP)
  r.Router.Use(middleware.Logger)
  r.Router.Use(middleware.Recoverer)
	r.Router.Use(middleware.Timeout(60 * time.Second))

  r.Router.Route("/api/stories", func(router chi.Router) {
    router.With(paginate).Get("/", r.listStories)                           // GET /api/stories
    //r.With(paginate).Get("/{month}-{day}-{year}", listArticlesByDate) // GET /api/stories/01-16-2017

    router.Post("/", r.createStory)                                        // POST /api/articles
    //r.Get("/search", searchArticles)                                  // GET /api/articles/search

    // Regexp url parameters:
    //r.Get("/{articleSlug:[a-z-]+}", getArticleBySlug)                // GET /api/articles/home-is-toronto

    // Subrouters:
    router.Route("/{storyId}", func(router chi.Router) {
      router.Use(r.StoryCtx)
      router.Get("/", r.getStory)                                          // GET /api/articles/123
      router.Put("/", r.updateStory)                                       // PUT /api/articles/123
      router.Delete("/", r.deleteStory)                                    // DELETE /api/articles/123
    })
  })

  // Mount the admin sub-router
  //r.Mount("/admin", adminRouter())

}


func (sRouter StoryServerRouter)StoryCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    storyId := chi.URLParam(r, "storyId")
    story, err := sRouter.app.FindStory(r.Context(),storyId)
    if err != nil {
      http.Error(w, http.StatusText(404), 404)
      return
    }
    ctx := context.WithValue(r.Context(), STORY_ID, story)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

//paginate handles the pagination of requests
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func (sRouter StoryServerRouter) getStory(w http.ResponseWriter, r *http.Request){
	story := r.Context().Value(STORY_ID).(common.Story)
	err:=json.NewEncoder(w).Encode(story)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func  (sRouter StoryServerRouter) createStory(w http.ResponseWriter, r *http.Request){
	var story common.Story
	err := json.NewDecoder(r.Body).Decode(&story)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s, err:=sRouter.app.CreateStory(r.Context(), story.User, story.Title)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err =json.NewEncoder(w).Encode(s)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func  (sRouter StoryServerRouter) listStories(w http.ResponseWriter, r *http.Request){
	s, err:=sRouter.app.FindAllStories(r.Context())
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
	}
	err =json.NewEncoder(w).Encode(s)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func  (sRouter StoryServerRouter) updateStory(w http.ResponseWriter, r *http.Request){
	var story common.Story
	sOld := r.Context().Value(STORY_ID).(common.Story)
	err := json.NewDecoder(r.Body).Decode(&story)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	story, err=sRouter.app.UpdateStory(r.Context(),sOld.ID, story.User, story.Title, story.Text )
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err =json.NewEncoder(w).Encode(story)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func  (sRouter StoryServerRouter) deleteStory(w http.ResponseWriter, r *http.Request){
	story := r.Context().Value(STORY_ID).(common.Story)
	err:=sRouter.app.DeleteStory(r.Context(), story.ID)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
