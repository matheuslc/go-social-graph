// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package server

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new post for an user
	// (POST /api/post)
	PostHandler(ctx echo.Context) error
	// Repost an user post
	// (POST /api/post/{id}/repost)
	RepostHandler(ctx echo.Context, id openapi_types.UUID) error
	// Retrieve user profile information
	// (GET /api/profile/{user_id})
	ProfileHandler(ctx echo.Context, userId openapi_types.UUID) error
	// Create a new user
	// (POST /api/user)
	CreateUser(ctx echo.Context) error
	// Follows a user
	// (POST /api/user/{id}/follow/{from})
	FollowHandler(ctx echo.Context, id openapi_types.UUID, from openapi_types.UUID) error
	// Retrieve the timeline for a user
	// (GET /api/user/{id}/timeline)
	TimelineHandler(ctx echo.Context, id openapi_types.UUID) error
	// Unfollow a user
	// (POST /api/user/{id}/unfollow/{from})
	UnfollowHandler(ctx echo.Context, id openapi_types.UUID, from openapi_types.UUID) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostHandler converts echo context to params.
func (w *ServerInterfaceWrapper) PostHandler(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostHandler(ctx)
	return err
}

// RepostHandler converts echo context to params.
func (w *ServerInterfaceWrapper) RepostHandler(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RepostHandler(ctx, id)
	return err
}

// ProfileHandler converts echo context to params.
func (w *ServerInterfaceWrapper) ProfileHandler(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "user_id" -------------
	var userId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "user_id", runtime.ParamLocationPath, ctx.Param("user_id"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter user_id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ProfileHandler(ctx, userId)
	return err
}

// CreateUser converts echo context to params.
func (w *ServerInterfaceWrapper) CreateUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateUser(ctx)
	return err
}

// FollowHandler converts echo context to params.
func (w *ServerInterfaceWrapper) FollowHandler(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "from" -------------
	var from openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "from", runtime.ParamLocationPath, ctx.Param("from"), &from)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter from: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FollowHandler(ctx, id, from)
	return err
}

// TimelineHandler converts echo context to params.
func (w *ServerInterfaceWrapper) TimelineHandler(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.TimelineHandler(ctx, id)
	return err
}

// UnfollowHandler converts echo context to params.
func (w *ServerInterfaceWrapper) UnfollowHandler(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "from" -------------
	var from openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "from", runtime.ParamLocationPath, ctx.Param("from"), &from)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter from: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UnfollowHandler(ctx, id, from)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/api/post", wrapper.PostHandler)
	router.POST(baseURL+"/api/post/:id/repost", wrapper.RepostHandler)
	router.GET(baseURL+"/api/profile/:user_id", wrapper.ProfileHandler)
	router.POST(baseURL+"/api/user", wrapper.CreateUser)
	router.POST(baseURL+"/api/user/:id/follow/:from", wrapper.FollowHandler)
	router.GET(baseURL+"/api/user/:id/timeline", wrapper.TimelineHandler)
	router.POST(baseURL+"/api/user/:id/unfollow/:from", wrapper.UnfollowHandler)

}
