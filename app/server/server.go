package server

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/joho/godotenv"
	email "github.com/softcorp-io/hqs-email-service/email"
	handler "github.com/softcorp-io/hqs-email-service/handler"
	proto "github.com/softcorp-io/hqs_proto/go_hqs/hqs_email_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Init - initialize .env variables.
func Init(zapLog *zap.Logger) {
	if err := godotenv.Load("hqs.env"); err != nil {
		zapLog.Error(fmt.Sprintf("Could not load hqs.env with err %v", err))
	}
}

// load environment to send emails and create email struct
func loadEmailAuth(zapLog *zap.Logger) *email.AuthEmail {
	emailFrom, ok := os.LookupEnv("EMAIL_FROM")
	if !ok {
		zapLog.Fatal("Could not load emailFrom")
	}
	emailPassword, ok := os.LookupEnv("EMAIL_PASSOWRD")
	if !ok {
		zapLog.Fatal("Could not load password for email")
	}
	smtpHost, ok := os.LookupEnv("SMPT_HOST")
	if !ok {
		zapLog.Fatal("Could not load smptHost for email")
	}
	smptPort, ok := os.LookupEnv("SMPT_PORT")
	if !ok {
		zapLog.Fatal("Could not load smptPort for email")
	}
	return email.NewAuthEmail(emailFrom, emailPassword, smtpHost, smptPort, zapLog)
}

// Run - runs a go microservice. Uses zap for logging and a waitGroup for async testing.
func Run(zapLog *zap.Logger, wg *sync.WaitGroup) {
	// create email service
	authEmail := loadEmailAuth(zapLog)

	handle := handler.NewHandler(zapLog, authEmail)

	// create the service and run the service
	port, ok := os.LookupEnv("SERVICE_PORT")
	if !ok {
		zapLog.Fatal("Could not get service port")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zapLog.Fatal(fmt.Sprintf("Failed to listen with err %v", err))
	}
	defer lis.Close()

	zapLog.Info(fmt.Sprintf("Service running on port: %s", port))

	// setup grpc
	grpcServer := grpc.NewServer()

	// register handler
	proto.RegisterEmailServiceServer(grpcServer, handle)

	// run the server
	wg.Done()
	if err := grpcServer.Serve(lis); err != nil {
		zapLog.Fatal(fmt.Sprintf("Failed to serve with err %v", err))
	}
}
