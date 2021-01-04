package testclient

import (
	"testing"

	proto "github.com/softcorp-io/hqs_proto/go_hqs/hqs_email_service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestSendResetEmail(t *testing.T) {
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	assert.Nil(t, err)
	client := proto.NewEmailServiceClient(conn)
	_, err = SendResetPasswordEmail(client, "Lucas E. B. Orellana", "asdz213k42", "oscar@orellana.dk")
	assert.Nil(t, err)
}
