// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: api/proto/admin_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AdminService_CreateCategory_FullMethodName  = "/pb.AdminService/CreateCategory"
	AdminService_GetCategories_FullMethodName   = "/pb.AdminService/GetCategories"
	AdminService_GetUsers_FullMethodName        = "/pb.AdminService/GetUsers"
	AdminService_CreatePromocode_FullMethodName = "/pb.AdminService/CreatePromocode"
	AdminService_GetPromocodes_FullMethodName   = "/pb.AdminService/GetPromocodes"
)

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminServiceClient interface {
	CreateCategory(ctx context.Context, in *AddCategoryRequest, opts ...grpc.CallOption) (*Response, error)
	GetCategories(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetCategoryResponse, error)
	GetUsers(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetUsersResponse, error)
	CreatePromocode(ctx context.Context, in *CreatePromocodeRequest, opts ...grpc.CallOption) (*Response, error)
	GetPromocodes(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetPromocodeResponse, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) CreateCategory(ctx context.Context, in *AddCategoryRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, AdminService_CreateCategory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetCategories(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetCategoryResponse, error) {
	out := new(GetCategoryResponse)
	err := c.cc.Invoke(ctx, AdminService_GetCategories_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetUsers(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetUsersResponse, error) {
	out := new(GetUsersResponse)
	err := c.cc.Invoke(ctx, AdminService_GetUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) CreatePromocode(ctx context.Context, in *CreatePromocodeRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, AdminService_CreatePromocode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetPromocodes(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GetPromocodeResponse, error) {
	out := new(GetPromocodeResponse)
	err := c.cc.Invoke(ctx, AdminService_GetPromocodes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
// All implementations must embed UnimplementedAdminServiceServer
// for forward compatibility
type AdminServiceServer interface {
	CreateCategory(context.Context, *AddCategoryRequest) (*Response, error)
	GetCategories(context.Context, *Request) (*GetCategoryResponse, error)
	GetUsers(context.Context, *Request) (*GetUsersResponse, error)
	CreatePromocode(context.Context, *CreatePromocodeRequest) (*Response, error)
	GetPromocodes(context.Context, *Request) (*GetPromocodeResponse, error)
	mustEmbedUnimplementedAdminServiceServer()
}

// UnimplementedAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (UnimplementedAdminServiceServer) CreateCategory(context.Context, *AddCategoryRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCategory not implemented")
}
func (UnimplementedAdminServiceServer) GetCategories(context.Context, *Request) (*GetCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategories not implemented")
}
func (UnimplementedAdminServiceServer) GetUsers(context.Context, *Request) (*GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedAdminServiceServer) CreatePromocode(context.Context, *CreatePromocodeRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePromocode not implemented")
}
func (UnimplementedAdminServiceServer) GetPromocodes(context.Context, *Request) (*GetPromocodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPromocodes not implemented")
}
func (UnimplementedAdminServiceServer) mustEmbedUnimplementedAdminServiceServer() {}

// UnsafeAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServiceServer will
// result in compilation errors.
type UnsafeAdminServiceServer interface {
	mustEmbedUnimplementedAdminServiceServer()
}

func RegisterAdminServiceServer(s grpc.ServiceRegistrar, srv AdminServiceServer) {
	s.RegisterService(&AdminService_ServiceDesc, srv)
}

func _AdminService_CreateCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).CreateCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_CreateCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).CreateCategory(ctx, req.(*AddCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_GetCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetCategories(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_GetUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetUsers(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_CreatePromocode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePromocodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).CreatePromocode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_CreatePromocode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).CreatePromocode(ctx, req.(*CreatePromocodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetPromocodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetPromocodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_GetPromocodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetPromocodes(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminService_ServiceDesc is the grpc.ServiceDesc for AdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCategory",
			Handler:    _AdminService_CreateCategory_Handler,
		},
		{
			MethodName: "GetCategories",
			Handler:    _AdminService_GetCategories_Handler,
		},
		{
			MethodName: "GetUsers",
			Handler:    _AdminService_GetUsers_Handler,
		},
		{
			MethodName: "CreatePromocode",
			Handler:    _AdminService_CreatePromocode_Handler,
		},
		{
			MethodName: "GetPromocodes",
			Handler:    _AdminService_GetPromocodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/admin_service.proto",
}
