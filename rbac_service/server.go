package main

import (
	"context"
	"log"

	pb "github.com/krishnachowdaryvanam/authboard/rbac_service/rbacpb"
)

type RBACServer struct {
	pb.UnimplementedRbacServiceServer
}

func (s *RBACServer) CheckAccess(ctx context.Context, req *pb.CheckAccessRequest) (*pb.CheckAccessResponse, error) {
	allowed, err := CheckUserAccess(req.UserId, req.Resource)
	if err != nil {
		return &pb.CheckAccessResponse{Allowed: false}, err
	}
	return &pb.CheckAccessResponse{Allowed: allowed}, nil
}

func (s *RBACServer) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	err := InsertUserRole(req.UserId, req.Role)
	if err != nil {
		return &pb.AssignRoleResponse{Success: false}, err
	}
	log.Printf("Assigned role '%s' to user '%s'", req.Role, req.UserId)
	return &pb.AssignRoleResponse{Success: true}, nil
}

func (s *RBACServer) RemoveRole(ctx context.Context, req *pb.RemoveRoleRequest) (*pb.RemoveRoleResponse, error) {
	err := RemoveUserRole(req.UserId, req.Role)
	if err != nil {
		return &pb.RemoveRoleResponse{Success: false}, err
	}
	log.Printf("Removed role '%s' from user '%s'", req.Role, req.UserId)
	return &pb.RemoveRoleResponse{Success: true}, nil
}
