package main

import (
	middleware "API-GATEWAY/Middleware"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	pb "github.com/Y-Fleet/Grpc-Api/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /login")
		conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		log.Println("Login Api Call")
		client := pb.NewAuthServiceClient(conn)
		login := r.FormValue("login")
		password := r.FormValue("password")

		req := &pb.CheckUsersRequest{Login: login, Password: password}
		res, err := client.CheckUsers(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(map[string]string{"message": res.Message, "token": res.Token, "refresh_token": res.TokenRefresh})
		w.Write(jsonResponse)
		log.Println("Sending request for /login")
	}).Methods("POST")

	r.HandleFunc("/RefreshToken", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Refresh Api Call")
		conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		client := pb.NewAuthServiceClient(conn)
		refreshToken := r.FormValue("refreshToken")
		req := &pb.RefreshTokenRequest{RefreshToken: refreshToken}
		res, err := client.RefreshToken(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(map[string]string{"token": res.Token})
		w.Write(jsonResponse)
	}).Methods("POST")
	r.Handle("/WarehouseDash", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("WarehouseDash Api Call")
		conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		client := pb.NewWarehouseManagementServiceClient(conn)
		req := &pb.RenderWarehouseRequest{}
		res, err := client.RenderWarehouse(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(res.Warehouse)
		w.Write(jsonResponse)
	}))).Methods("GET")

	r.Handle("/AddItem", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AddItem Api Call")
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Error reading request body: %s", err)
		}
		fmt.Println("Request body:", string(requestBody))
		defer conn.Close()
		client := pb.NewInventoryServiceClient(conn)
		req := &pb.AddItemRequest{}
		_ = json.Unmarshal(requestBody, req)
		res, err := client.AddItem(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(res.Message)
		w.Write(jsonResponse)
	}))).Methods("PUT")

	r.Handle("/GetItem", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetItem Api Call")
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		client := pb.NewInventoryServiceClient(conn)
		req := &pb.GetItemRequest{}
		res, err := client.GetItem(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(res.Items)
		w.Write(jsonResponse)
	}))).Methods("GET")

	r.Handle("/DelItem", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("DelItem Api Call")
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		client := pb.NewInventoryServiceClient(conn)

		var requestBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ID := requestBody["ID"]

		req := &pb.DelItemRequest{ID: ID}
		res, err := client.DelItem(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(res.Message)
		w.Write(jsonResponse)
	}))).Methods("DELETE")

	r.Handle("/InfoWarehouse", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("InfoWarehouse Api Call")
		warehouseId := r.URL.Query().Get("id")
		if warehouseId == "" {
			http.Error(w, "Missing warehouse ID parameter", http.StatusBadRequest)
			return
		}
		conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		client := pb.NewWarehouseManagementServiceClient(conn)
		req := &pb.InfoWarehouseRequest{Id: warehouseId}
		res, err := client.InfoWarehouse(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(res)
		w.Write(jsonResponse)
	}))).Methods("GET")

	r.Handle("/InfoWarehouse/{warehouse_id}/addStock", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("test API Call")
		vars := mux.Vars(r)
		warehouseID := vars["warehouse_id"]
		conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		var addStockReq struct {
			ItemID   string `json:"item_id"`
			Quantity int32  `json:"quantity"`
		}
		err = json.NewDecoder(r.Body).Decode(&addStockReq)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client := pb.NewWarehouseManagementServiceClient(conn)
		req := &pb.AddStockRequest{IDWharehouse: warehouseID, IDItems: addStockReq.ItemID, Stock: addStockReq.Quantity}
		res, err := client.AddStock(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse, _ := json.Marshal(res)
		w.Write(jsonResponse)
	}))).Methods("POST")

	r.Handle("/GetFleet", middleware.JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Fleet Api Call")
		conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		client := pb.NewFleetServiceClient(conn)
		req := &pb.Empty{}
		res, err := client.GetFleet(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(res.Vehicles)
		w.Write(jsonResponse)
	}))).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	log.Println("Starting API gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
