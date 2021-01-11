package parser

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

// Context wrangles all of the captured data about your input service declaration file. It tracks
// the module/package information, the service(s) that were defined, the request/response structs
// that were defined in the file, etc. It's the output of Parse() and is the input value when we
// evaluate Go templates to generate other source files based on this service definition info.
type Context struct {
	// File is the entire syntax tree from when we parsed your input file.
	File *ast.File
	// Path is the relative path to the service definition file we're parsing.
	Path string
	// Package contains information about the package where the service definition resides.
	Package *PackageDeclaration
	// OutputPackage contains information about the package where the generated code will go.
	OutputPackage *PackageDeclaration
	// Module contains info from "go.mod" about the entire module where the service/package is defined.
	Module *ModuleDeclaration
	// Services encapsulates snapshot info for all service interfaces that were defined in the input file.
	Services []*ServiceDeclaration
	// Models encapsulates snapshot info for all service request/response structs that were defined in the input file.
	Models []*ServiceModelDeclaration

	// currentService is used internally when processing method info to know what service you're "inside" of.
	currentService *ServiceDeclaration
	// currentMethod is used internally when processing field info to know what service operation you're "inside" of.
	currentMethod *ServiceMethodDeclaration
}

// AddService appends a service definition to this parsing context.
func (ctx *Context) AddService(service *ServiceDeclaration) {
	ctx.Services = append(ctx.Services, service)
}

// AddModel appends a request/response struct definition to this parsing context.
func (ctx *Context) AddModel(model *ServiceModelDeclaration) {
	ctx.Models = append(ctx.Models, model)
}

// ModelByName looks through "Models" to find the one whose method/function name matches 'name'.
func (ctx Context) ModelByName(name string) *ServiceModelDeclaration {
	for _, model := range ctx.Models {
		if model.Name == name {
			return model
		}
	}
	return nil
}

// ServiceDeclaration wrangles all of the information we could grab about the service from the
// interface that defined it.
type ServiceDeclaration struct {
	// Name is the name of the service/interface.
	Name string
	// Methods are all of the functions explicitly defined on this service.
	Methods []*ServiceMethodDeclaration
	// Node is the syntax tree object for the interface that described this service.
	Node *ast.Object
}

// AddMethod appends a method definition to this service, indicating that the service contains this operation.
func (service *ServiceDeclaration) AddMethod(method *ServiceMethodDeclaration) {
	service.Methods = append(service.Methods, method)
}

// MethodByName fetches the service operation with the given function name. This returns nil when there
// are no functions in this interface/service by that name.
func (service ServiceDeclaration) MethodByName(name string) *ServiceMethodDeclaration {
	for _, m := range service.Methods {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// ServiceMethodDeclaration defines a single operation/function within a service (one of the interface functions).
type ServiceMethodDeclaration struct {
	// Name is the name of the function defined in the service interface (the function name to call this operation).
	Name string
	// Request contains the details about the model/type/struct for this operation's input/request value.
	Request *ServiceModelDeclaration
	// Response contains the details about the model/type/struct for this operation's output/response value.
	Response *ServiceModelDeclaration
	// HTTPMethod indicates if the RPC gateway should use a GET, POST, etc when exposing this operation via HTTP.
	HTTPMethod string
	// HTTPPath defines the URL pattern to provide to the gateway's router/mux to access this operation.
	HTTPPath string
	// HTTPStatus indicates what success status code the gateway should use when responding via HTTP (e.g. 200, 202, etc)
	HTTPStatus int
	// Node is the syntax tree object that defined this function within the service interface.
	Node *ast.Field
}

// String returns the method signature for this operation.
func (method ServiceMethodDeclaration) String() string {
	return fmt.Sprintf("%s(context.Context, %v) (%v, error)",
		method.Name,
		method.Request,
		method.Response,
	)
}

// ServiceModelDeclaration contains information about request/response structs defined in your declaration file.
type ServiceModelDeclaration struct {
	// Name is the name of the type/struct used when defining the request/response value.
	Name string
	// Node is the syntax tree object that defined this type/struct.
	Node *ast.Object
}

// String just returns the model type's name.
func (model ServiceModelDeclaration) String() string {
	return model.Name
}

// ModuleDeclaration contains information about the Go module that the service belongs
// to. This is information scraped from project's "go.mod" file.
type ModuleDeclaration struct {
	// Name is the fully qualified module name (e.g. "github.com/someuser/modulename")
	Name string
	// Directory is the absolute path to the root directory of the module (where go.mod resides).
	Directory string
}

// GoMod returns the absolute path to the "go.mod" file for this module on the system running frodoc.
func (module ModuleDeclaration) GoMod() string {
	return filepath.Join(module.Directory, "go.mod")
}

// PackageDeclaration defines the subpackage that the service resides in.
type PackageDeclaration struct {
	// Name is just the raw package name (no path info)
	Name string
	// Import is the fully qualified package name (e.g. "github.com/someuser/modulename/foo/bar/baz")
	Import string
	// Directory is the absolute path to the package.
	Directory string
}