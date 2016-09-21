package vault

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	httptransport "github.com/go-kit/kit/transport/http"
// 	"golang.org/x/net/context"
// )

// func main() {
// 	ctx := context.Background()
// 	srv := hasherService{}
// 	http.Handle("/json/hash", httptransport.NewServer(
// 		ctx,
// 		makeHashEndpoint(srv),
// 		decodeHashRequest,
// 		encodeResponse,
// 	))
// 	http.Handle("/json/validate", httptransport.NewServer(
// 		ctx,
// 		makeValidateEndpoint(srv),
// 		decodeValidateRequest,
// 		encodeResponse,
// 	))
// 	addr := ":8080"
// 	log.Println("serving through", addr)
// 	log.Fatal(http.ListenAndServe(addr, nil))
// }

// func decodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error) {
// 	var req hashRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return req, nil
// }

// func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
// 	var req validateRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return req, nil
// }

// func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
// 	return json.NewEncoder(w).Encode(response)
// }
