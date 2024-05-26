package server

import (
	"log"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/grpc_api/grpc_gateway/server/middleware"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/utils"
)

func (server *RestServer) SetupRouter(mux *http.ServeMux) {
	//EditProfile
	editProfile := middleware.ApplyAccessRoleMiddleware(server.token, server.log, utils.MERCHANT, http.HandlerFunc(server.EditProfile))
	mux.Handle("/api/merchant/profile", MethodHandler{Method: "PATCH", Handler: editProfile})

	//InsertProduct
	addProduct := middleware.ApplyAccessRoleMiddleware(server.token, server.log, utils.MERCHANT, http.HandlerFunc(server.InsertProduct))
	mux.Handle("/api/merchant/product", MethodHandler{Method: "POST", Handler: addProduct})

	//EditProduct
	editProduct := middleware.ApplyAccessRoleMiddleware(server.token, server.log, utils.MERCHANT, http.HandlerFunc(server.EditProduct))
	mux.Handle("/api/merchant/product/{id}", MethodHandler{Method: "PATCH", Handler: editProduct})
}

type MethodHandler struct {
	Method  string
	Handler http.Handler
}

func (mh MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " request method")
	log.Println(mh.Method, " my method")
	log.Println(r.URL)

	if r.Method != mh.Method {
		helpers.ErrorJson(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	mh.Handler.ServeHTTP(w, r)
}
