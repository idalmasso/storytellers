package api

import (
	"net/http"

	"github.com/idalmasso/storytellers/backend/api/router"
	"github.com/idalmasso/storytellers/backend/api/types"
)


type storyServer struct{
	router router.StoryServerRouter
}

func (s storyServer)Run(port string){
	
  http.ListenAndServe(":"+port, s.router.Router)
}

func NewStoryServer(app types.StoryTellerAppInterface) storyServer{
	router:=router.CreateRouter(app)
	storyServer:=storyServer{router: router }
	
	return storyServer
}

