package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/config"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/db"
	generated "github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/controller"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/usecase/pr-review"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/usecase/repository"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(logger *zap.Logger, cfg *config.Config) {
	const GracefulShutdownTimeout = time.Second * 3
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.PG.URL)

	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return
	}

	defer dbPool.Close()

	// migrations
	db.SetupPostgres(dbPool, logger)

	repo := repository.NewPostgresRepository(logger, dbPool)
	useCases := pr_review.New(logger, repo, repo, repo)

	ctrl := controller.New(logger, useCases, useCases, useCases)

	go runRest(ctx, cfg, logger)
	go runGrpc(cfg, logger, ctrl)

	<-ctx.Done()
	time.Sleep(GracefulShutdownTimeout)
}

func runRest(ctx context.Context, cfg *config.Config, logger *zap.Logger) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	address := "localhost:" + cfg.GRPC.Port
	err := generated.RegisterPRReviewerServiceHandlerFromEndpoint(ctx, mux, address, opts)

	if err != nil {
		logger.Error("can not register grpc gateway", zap.Error(err))
		os.Exit(-1)
	}

	gatewayPort := ":" + cfg.GatewayPort
	logger.Info("gateway listening at port", zap.String("port", gatewayPort))

	if err = http.ListenAndServe(gatewayPort, mux); err != nil {
		logger.Error("gateway listen error", zap.Error(err))
	}
}

func runGrpc(cfg *config.Config, logger *zap.Logger, libraryService generated.PRReviewerServiceServer) {
	port := ":" + cfg.GRPC.Port
	lis, err := net.Listen("tcp", port)

	if err != nil {
		logger.Error("can not open tcp socket", zap.Error(err))
		os.Exit(-1)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	generated.RegisterPRReviewerServiceServer(s, libraryService)

	logger.Info("grpc server listening at port", zap.String("port", port))

	if err = s.Serve(lis); err != nil {
		logger.Error("grpc server listen error", zap.Error(err))
	}
}
