package openapi

//go:generate oapi-codegen -package openapi -generate types -o types.gen.go openapi.yml
//go:generate oapi-codegen -package openapi -generate spec -o spec.gen.go openapi.yml
