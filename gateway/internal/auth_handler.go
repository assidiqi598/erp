package internal

import (
	"context"

	authpb "github.com/assidiqi598/erp/services/auth/proto"
)

// Gateway handler for AuthService
type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authClient authpb.AuthServiceClient
}

// Constructor for the AuthHandler
func NewAuthHandler(authClient authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

// ChangeEmail implements __.AuthServiceServer.
func (a *AuthHandler) ChangeEmail(ctx context.Context, req *authpb.ChangeEmailRequest) (*authpb.ChangeEmailResponse, error) {
	resp, err := a.authClient.ChangeEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ChangePassword implements __.AuthServiceServer.
func (a *AuthHandler) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	resp, err := a.authClient.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// LoginWithEmailAndPass implements __.AuthServiceServer.
func (a *AuthHandler) LoginWithEmailAndPass(ctx context.Context, req *authpb.LoginWithEmailAndPassRequest) (*authpb.LoginResponse, error) {
	resp, err := a.authClient.LoginWithEmailAndPass(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Register implements __.AuthServiceServer.
func (a *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	resp, err := a.authClient.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RequestToChangePassword implements __.AuthServiceServer.
func (a *AuthHandler) RequestToChangePassword(ctx context.Context, req *authpb.RequestToChangePasswordRequest) (*authpb.RequestToChangePasswordResponse, error) {
	resp, err := a.authClient.RequestToChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ResendVerificationEmail implements __.AuthServiceServer.
func (a *AuthHandler) ResendVerificationEmail(ctx context.Context, req *authpb.ResendVerificationEmailRequest) (*authpb.ResendVerificationEmailResponse, error) {
	resp, err := a.authClient.ResendVerificationEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// VerifyEmail implements __.AuthServiceServer.
func (a *AuthHandler) VerifyEmail(ctx context.Context, req *authpb.VerifyEmailRequest) (*authpb.VerifyEmailResponse, error) {
	resp, err := a.authClient.VerifyEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
