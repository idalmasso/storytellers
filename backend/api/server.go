package api

import (
	"net/http"

	"github.com/idalmasso/storytellers/backend/api/router"
	"github.com/idalmasso/storytellers/backend/api/types"
)


type storyServer struct{
	router router.StoryServerRouter
}

func (s storyServer)Run(){
	
  http.ListenAndServe(":3333", s.router.Router)
}

func NewStoryServer(app types.StoryTellerAppInterface) storyServer{
	router:=router.CreateRouter(app)
	storyServer:=storyServer{router: router }
	
	return storyServer
}

