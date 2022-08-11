package gapi

import (
	"context"
	"database/sql"

	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/pb"
	"github.com/burakkarasel/Bank-App/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LoginUser handles login gRPC requests
func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		// if err is sql.ErrNoRows username is invalid
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	// then we check for the password for given username
	err = util.CheckPassword(req.Password, user.HashedPassword)

	// if an error occurs it means that password is invalid
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "password is incorrect: %s", err)
	}

	// then we create an access token for this logged in user
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
	}

	// and then we create refresh token for this logged in user
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %s", err)
	}

	// here we create a new session
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.GetUsername(),
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot create session: %s", err)
	}

	// and then we send user and access token as a response
	resp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return resp, nil
}
