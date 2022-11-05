package openapi

//go:generate oapi-codegen -package server -generate types -o ../server/types.gen.go openapi.yml
//go:generate oapi-codegen -package server -generate spec -o ../server/spec.gen.go openapi.yml
//go:generate oapi-codegen -package server -generate server -o ../server/server.gen.go openapi.yml
