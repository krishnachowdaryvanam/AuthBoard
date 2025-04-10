package main

import (
	"context"
	"log"

	"github.com/krishnachowdaryvanam/authboard/tenant_service/tenantpb"
)

type TenantServer struct {
	tenantpb.UnimplementedTenantServiceServer
}

func (s *TenantServer) CreateTenant(ctx context.Context, req *tenantpb.CreateTenantRequest) (*tenantpb.TenantResponse, error) {
	t, err := CreateTenant(req.Name)
	if err != nil {
		return nil, err
	}
	return &tenantpb.TenantResponse{
		Id:   t.ID,
		Name: t.Name,
	}, nil

}

func (s *TenantServer) GetTenant(ctx context.Context, req *tenantpb.GetTenantRequest) (*tenantpb.TenantResponse, error) {
	t, err := GetTenant(req.Id)
	if err != nil {
		log.Printf("GetTenant error: %v", err)
		return nil, err
	}
	return &tenantpb.TenantResponse{
		Id:   t.ID,
		Name: t.Name,
	}, nil
}

func (s *TenantServer) UpdateTenant(ctx context.Context, req *tenantpb.UpdateTenantRequest) (*tenantpb.TenantResponse, error) {
	t, err := UpdateTenant(req.Id, req.Name)
	if err != nil {
		log.Printf("UpdateTenant error: %v", err)
		return nil, err
	}
	return &tenantpb.TenantResponse{
		Id:   t.ID,
		Name: t.Name,
	}, nil
}

func (s *TenantServer) DeleteTenant(ctx context.Context, req *tenantpb.DeleteTenantRequest) (*tenantpb.DeleteTenantResponse, error) {
	success, err := DeleteTenant(req.Id)
	if err != nil {
		log.Printf("DeleteTenant error: %v", err)
		return nil, err
	}

	return &tenantpb.DeleteTenantResponse{Success: success}, nil
}
