package testclient

import (
	"context"

	proto "github.com/softcorp-io/hqs_proto/go_hqs/hqs_email_service"
)

// SendResetPasswordEmail - send a reset email to the user.
func SendResetPasswordEmail(cl proto.EmailServiceClient, name string, token string, to string) (*proto.Response, error) {
	// build context
	ctx := context.Background()

	emailResponse, err := cl.SendResetPasswordEmail(ctx, &proto.ResetPasswordEmail{
		Name:  name,
		To:    []string{to},
		Token: token,
	})
	if err != nil {
		return &proto.Response{}, err
	}
	return emailResponse, nil
}
