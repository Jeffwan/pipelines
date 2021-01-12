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
	"github.com/golang/glog"
	"github.com/kubeflow/pipelines/backend/src/apiserver/common"
	"github.com/kubeflow/pipelines/backend/src/common/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const BYTEDANCE_JWT_PUBLIC_KEY_ENDPOINT = "https://cloud.bytedance.net/auth/api/v1/public_key"

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

	// Compatibility - for non-experiment requests. Frontend only append token to experiment calls.
	// and we don't want to block other calls like listPipelines, listRuns
	if len(jwtToken) == 0 {
		return ctx, nil
	}
	// We don't want to check header until we append it to all backend requests
	//if len(jwtToken) == 0 {
	//	return nil, status.Errorf(codes.InvalidArgument, "missing 'x-jwt-token' header")
	//}

	if strings.TrimSpace(jwtToken[0]) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty 'x-jwt-token' header")
	}

	jwtAuthenticator, err := NewJwtAuthenticator(BYTEDANCE_JWT_PUBLIC_KEY_ENDPOINT)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Fail to initialize jwt authenticator: %v", err)
	}
	username, err := jwtAuthenticator.AuthenticateToken(jwtToken[0])
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "'x-jwt-token' not valid: %v", err)
	}

	//jwtToken = []string{"jiaxin.wang"}

	// TODO: if we use interceptor, then all other calls with jwt-token will be failed
	// Share we consider to proxy all request from frontend to backend with jwt-token?


	// Get user claim from jwt token and store it in the metadata context for future usage
	// https://chromium.googlesource.com/external/github.com/grpc/grpc-go/+show/refs/heads/master/Documentation/grpc-metadata.md
	newMd := metadata.Pairs("user", username)
	newCtx := metadata.NewIncomingContext(ctx, metadata.Join(md, newMd))
	return newCtx, nil
}
