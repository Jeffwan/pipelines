// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"github.com/kubeflow/pipelines/backend/src/apiserver/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/pipelines/backend/src/common/util"
	"google.golang.org/grpc"
)

// apiServerInterceptor implements UnaryServerInterceptor that provides the common wrapping logic
// to be executed before and after all API handler calls, e.g. Logging, error handling.
// For more details, see https://github.com/grpc/grpc-go/blob/master/interceptor.go
func apiServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	glog.Infof("%v handler starting", info.FullMethod)
	newCtx := ctx
	if common.IsMultiUserInSingleNamespaceMode() {
		// validate X-Jwt-Token
		newCtx, err = validateJwtToken(ctx)
		if err != nil {
			return nil ,err
		}
	}

	resp, err = handler(newCtx, req)
	if err != nil {
		util.LogError(util.Wrapf(err, "%s call failed", info.FullMethod))
		// Convert error to gRPC errors
		err = util.ToGRPCError(err)
		return
	}
	glog.Infof("%v handler finished", info.FullMethod)
	return
}

func validateJwtToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get metadata")
	}

	jwtToken := md["x-jwt-token"]
	if len(jwtToken) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing 'x-jwt-token' header")
	}

	if strings.Trim(jwtToken[0], " ") == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty 'x-jwt-token' header")
	}

	// TODO: we need to parse JWT token and get user name out
	// if jwt token if invalid, 400. check generate
	// otherwise, save user into the context? Where to store?

	// Assume we decode JWT token here, we want to retrieve user and save into metadata
	// We want to store jwt token once to optimize the performance

	// https://chromium.googlesource.com/external/github.com/grpc/grpc-go/+show/refs/heads/master/Documentation/grpc-metadata.md
	newMd := metadata.Pairs("user", jwtToken[0])
	newCtx := metadata.NewIncomingContext(ctx, metadata.Join(md, newMd))
	return newCtx, nil
}
