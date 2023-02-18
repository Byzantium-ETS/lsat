package proxy

import (
	"fmt"
	"net/http"
)

type Server struct {
	database map[string]int
}

func HttpServer() {
	fmt.Println("Server launched!")
	http.HandleFunc("/", HandlePaymentRequest) // Retourné dans le cas où la platforme reçoit un Token invalide
	http.ListenAndServe(":8080", nil)          // Ca devrait etre lié à la platforme?
}

func HandlePaymentRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusPaymentRequired)
	w.Header().Set("WWW-Authenticate", "application/json") // A revoir le second argument

	// 1. Ici on devrait générer un invoice à l'aide de notre node.
	// 2. Si le client à son macaroon, on devrait aller le chercher dans la request, sinon on devrait en générer un.
	// m := PaymentRequest{"", ""}
	// b, _ := json.Marshal(m)

	// w.Write(b) // Envoyer la requete JSON? Je suis habitué à devoir flush.
}
