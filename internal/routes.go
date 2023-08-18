package internal

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"queue-manager/internal/controllers"
	"queue-manager/internal/providers"
	"queue-manager/internal/structures"
)

func dashboardAuthenticationMiddleware(conf *structures.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !conf.DashboardAuthentication.Required {
			next.ServeHTTP(w, r)
			return
		}

		session, err := r.Cookie("session")
		if err == nil {

			cookieValue := session.Value

			token, err := jwt.Parse(cookieValue, func(token *jwt.Token) (interface{}, error) {
				// Make sure the signing method is correct
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(conf.DashboardAuthentication.Secret), nil
			})

			if err != nil {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			// Check if the token is valid
			if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				next.ServeHTTP(w, r)
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitRoutes(apiController *controllers.ApiController, adminController *controllers.AdminController, conf *structures.Config) providers.RouterProviderInterface {
	routers := providers.NewRouterProvider()

	// swagger routes
	routers.Get("/schema", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, conf.Template.TplDir+"swaggerui/index.html")
	}))

	// http client API
	// swagger:route GET /api/v1/queue/list api listQueues
	// List queues
	// responses:
	//   200
	routers.Get("/api/v1/queue/list", http.HandlerFunc(apiController.List))

	// swagger:route POST /api/v1/queue/create api createQueue
	// Create a queue
	// Responses:
	//   201
	// Parameters:
	// + name: queueName
	//   in: body
	//   description: Queue name
	//   required: true
	//   schema:
	//     type: CreateQueuePayload
	routers.Post("/api/v1/queue/create", http.HandlerFunc(apiController.CreateQueue))
	// swagger:route POST /api/v1/message/create api createMessage
	// Create a message
	// parameters:
	//   + name: message
	//     in: body
	//     description: Message object
	//     required: true
	//     schema:
	//       type: CreateMessagePayload
	// responses:
	//   201
	routers.Post("/api/v1/message/create", http.HandlerFunc(apiController.CreateMessage))
	// swagger:route GET /api/v1/message/list api listMessagesInQueue
	// List messages in a queue
	// parameters:
	//   + name: queue_name
	//     in: query
	//     description: Queue name
	//     required: true
	// responses:
	//   200
	routers.Get("/api/v1/message/list", http.HandlerFunc(apiController.ListMessagesInQueue))
	// swagger:route GET /api/v1/message/get api getMessage
	// Get a message from a queue
	// parameters:
	//   + name: queue_name
	//     in: query
	//     description: Queue name
	//     required: true
	// responses:
	//   200
	routers.Get("/api/v1/message/get", http.HandlerFunc(apiController.GetMessage))
	// static for dashboard
	routers.Get("/login", http.HandlerFunc(adminController.Login))
	routers.Post("/login/do_login", http.HandlerFunc(adminController.DoLogin))
	routers.Get("/logout", http.HandlerFunc(adminController.Logout))
	routers.Get("/", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(adminController.Index)))
	routers.Get("/static/", http.HandlerFunc(adminController.Static))
	routers.Get("/queues", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(adminController.Queues)))
	routers.Get("/queue", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(adminController.Queue)))
	routers.Get("/queues/create", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(adminController.CreateQueue)))
	routers.Get("/queues/message", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(adminController.SendMessage)))
	//api for dashboard
	routers.Get("/api/dashboard/queues", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(apiController.List)))
	routers.Get("/api/dashboard/clients", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(apiController.TCPClients)))
	routers.Get("/api/dashboard/messages", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(apiController.MessagesCount)))
	routers.Get("/api/dashboard/queues/detailed", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(apiController.DetailedQueuesList)))
	routers.Get("/api/dashboard/queue/delete", dashboardAuthenticationMiddleware(conf, http.HandlerFunc(apiController.DeleteQueue)))

	return routers
}
