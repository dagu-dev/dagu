package admin

// type route struct {
// 	method  string
// 	pattern string
// 	handler http.HandlerFunc
// }

// func defaultRoutes(cfg *config.Config) []*route {
// 	tc := &handlers.TemplateConfig{
// 		NavbarColor: cfg.NavbarColor,
// 		NavbarTitle: cfg.NavbarTitle,
// 	}
// 	return []*route{
// 		{http.MethodGet, `^/?$`, handlers.HandleGetList(
// 			&handlers.DAGListHandlerConfig{DAGsDir: cfg.DAGs},
// 			tc,
// 		)},
// 		{http.MethodPost, `^/?$`, handlers.HandlePostList(
// 			&handlers.DAGListHandlerConfig{DAGsDir: cfg.DAGs},
// 		)},
// 		{http.MethodGet, `^/dags/?$`, handlers.HandleGetList(
// 			&handlers.DAGListHandlerConfig{DAGsDir: cfg.DAGs},
// 			tc,
// 		)},
// 		{http.MethodPost, `^/dags/?$`, handlers.HandlePostList(
// 			&handlers.DAGListHandlerConfig{DAGsDir: cfg.DAGs},
// 		)},
// 		{http.MethodGet, `^/dags/([^/]+)/?.*`, handlers.HandleGetDAG(
// 			&handlers.DAGHandlerConfig{
// 				DAGsDir:            cfg.DAGs,
// 				LogEncodingCharset: cfg.LogEncodingCharset,
// 			}, tc,
// 		)},
// 		{http.MethodPost, `^/dags/([^/]+)$`, handlers.HandlePostDAG(
// 			&handlers.PostDAGHandlerConfig{
// 				DAGsDir: cfg.DAGs,
// 				Bin:     cfg.Command,
// 				WkDir:   cfg.WorkDir,
// 			},
// 		)},
// 		{http.MethodDelete, `^/dags/([^/]+)$`, handlers.HandleDeleteDAG(
// 			&handlers.DeleteDAGHandlerConfig{
// 				DAGsDir: cfg.DAGs,
// 			},
// 		)},
// 		{http.MethodGet, `^/search/?.*$`, handlers.HandleGetSearch(cfg.DAGs, tc)},
// 		{http.MethodGet, `^/assets/js/.*$`, handlers.HandleGetAssets("/web")},
// 		{http.MethodGet, `^/assets/css/.*$`, handlers.HandleGetAssets("/web")},
// 	}
// }
