package gapi

import (
	"context"
	"database/sql"

	db "github.com/corey888773/golang-course/db/sqlc"
	"github.com/corey888773/golang-course/pb"
	"github.com/corey888773/golang-course/util"
	"github.com/corey888773/golang-course/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	err := invalidArgumentError(violations)
	if err != nil {
		return nil, err
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no matching user found: %s", err)
		}
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password: %s", err)
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshTokenPaylaod, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokeDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refrsh token: %s", err)
	}

	metadata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPaylaod.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPaylaod.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	res := &pb.LoginUserResponse{
		User:                  convertUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPaylaod.ExpiredAt),
		SessionId:             session.ID.String(),
	}

	return res, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return nil
}
