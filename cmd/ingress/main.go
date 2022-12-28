package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bonsai-oss/webbase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.fsrv.services/fsrvcorp/analytics/tracking-server/pkg/proto"
	"gitlab.fsrv.services/fsrvcorp/analytics/tracking-server/pkg/static"
)

func handlerBuilder(ingestClient proto.IngestServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var headers []*proto.Header
		for k, v := range r.Header {
			headers = append(headers, &proto.Header{
				Key:   k,
				Value: strings.Join(v, " ; "),
			})
		}

		_, requestIngestError := ingestClient.IngestMessage(r.Context(), &proto.Request{
			Timestamp: timestamppb.New(time.Now()),
			Headers:   headers,
			Path:      r.URL.Path,
		})
		if requestIngestError != nil {
			log.Printf("can not ingest request %v", requestIngestError)
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Write(static.OneByOnePixelPng)
	}
}

func main() {
	databaseGrpcAddress := os.Getenv("DATABASE_GRPC_ADDRESS")
	if databaseGrpcAddress == "" {
		log.Fatal("missing DATABASE_GRPC_ADDRESS")
		os.Exit(1)
	}

	conn, err := grpc.Dial(databaseGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with client %v", err)
	}
	client := proto.NewIngestServiceClient(conn)
	fmt.Printf("connected to database at %s", databaseGrpcAddress)

	router := webbase.NewRouter()
	router.PathPrefix("/").HandlerFunc(handlerBuilder(client))
	webbase.ServeRouter("tracking-server", router, webbase.WithSentryDebug(true))
}
