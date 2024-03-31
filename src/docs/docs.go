// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "get all users with their roles",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all users",
                "operationId": "get-all-users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    }
                }
            }
        },
        "/users/delete/{id}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "delete a user from the database",
                "produces": [
                    "application/json"
                ],
                "summary": "Delete a user",
                "operationId": "delete-user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User deleted"
                    }
                }
            }
        },
        "/users/insert": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "insert a new user into the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Insert a new user",
                "operationId": "insert-user",
                "parameters": [
                    {
                        "description": "user to insert",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            }
        },
        "/users/update": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "update an existing user in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update an existing user",
                "operationId": "update-user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "user to update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Policy": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.Role": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "rolePolicies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.RolePolicy"
                    }
                }
            }
        },
        "models.RolePolicy": {
            "type": "object",
            "properties": {
                "policy": {
                    "$ref": "#/definitions/models.Policy"
                },
                "policyId": {
                    "type": "integer"
                },
                "read": {
                    "type": "boolean"
                },
                "role": {
                    "$ref": "#/definitions/models.Role"
                },
                "roleId": {
                    "type": "integer"
                },
                "write": {
                    "type": "boolean"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "role": {
                    "$ref": "#/definitions/models.Role"
                },
                "roleId": {
                    "type": "integer"
                },
                "tenantId": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server for using Swagger with Gin.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}