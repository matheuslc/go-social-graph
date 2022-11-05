// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package server

// CreateUserIntent defines model for CreateUserIntent.
type CreateUserIntent struct {
	Username *string `json:"username,omitempty"`
}

// CreateUserFormdataRequestBody defines body for CreateUser for application/x-www-form-urlencoded ContentType.
type CreateUserFormdataRequestBody = CreateUserIntent