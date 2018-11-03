package serv

import (
	"github.com/gorilla/mux"
	"net/http"
	"web_service_GO/logger"
)

type Server interface {
	SetEndpoints() *http.ServeMux
	StartServe(int)
}

type DefaultServer struct {

}

func (ds * DefaultServer) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error.Println("Recovered", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (ds *DefaultServer) setRouters() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/submit", ds.handleSubmit).Methods("POST")
	r.HandleFunc("/check/{id}", ds.handleCheck).Methods("GET").GetError()
	r.HandleFunc("/", ds.handleRoot).Methods()
	return r
}

func (ds* DefaultServer) SetEndpoints() *http.ServeMux {
	r := ds.setRouters()
	flyMux := http.NewServeMux()
	flyMux.Handle("/", r)
	return flyMux
}



func (ds* DefaultServer) SetMiddlewares(mux *http.ServeMux) http.Handler {
	// Other middlewares may be easily added
	handler := ds.panicMiddleware(mux)
	return handler
}


func (ds* DefaultServer) StartServe(portNum int) {
	mux := ds.SetEndpoints()
	handler := ds.SetMiddlewares(mux)
	logger.Info.Println("Starting server at :", portNum)
	if err := http.ListenAndServe(string(portNum), handler); err != nil {
		logger.Error.Println("Can't start serving, check port num", err)
	}
}
