package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krishnachowdaryvanam/authboard/user_service/userpb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

func (*UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	id := uuid.New().String()
	createdAt := time.Now()
	updatedAt := createdAt

	_, err := db.Exec("INSERT INTO users (id, tenant_id, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, req.TenantId, req.Email, req.Password, req.Role, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Id:        id,
		TenantId:  req.TenantId,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: createdAt.String(),
		UpdatedAt: updatedAt.String(),
	}, nil
}

func (*UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	row := db.QueryRow("SELECT id, tenant_id, email, password, created_at, updated_at FROM users WHERE id=$1", req.Id)

	var u userpb.UserResponse
	err := row.Scan(&u.Id, &u.TenantId, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (*UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	updatedAt := time.Now()

	_, err := db.Exec("UPDATE users SET email=$1, password=$2, role=$3, updated_at=$4 WHERE id=$5",
		req.Email, req.Password, req.Role, updatedAt, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Id:        req.Id,
		Email:     req.Email,
		Password:  req.Password,
		UpdatedAt: updatedAt.String(),
	}, nil
}

func (*UserServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, err
	}
	return &userpb.DeleteUserResponse{Success: true}, nil
}
