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
        "/auth/signin": {
            "post": {
                "description": "Sign in with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign in with email and password",
                "parameters": [
                    {
                        "description": "User sign in request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.EmailSignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User sign in successful",
                        "schema": {
                            "$ref": "#/definitions/models.SignInSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Malformed Request",
                        "schema": {
                            "$ref": "#/definitions/models.GenericErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid Credentials",
                        "schema": {
                            "$ref": "#/definitions/models.GenericErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.GenericErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/signout": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Invalidates the user's access and refresh tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign out user",
                "responses": {
                    "200": {
                        "description": "User sign out successful",
                        "schema": {
                            "$ref": "#/definitions/models.GenericSuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.GenericErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Signs up a new user with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "User sign up request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User creation successful",
                        "schema": {
                            "$ref": "#/definitions/models.GenericSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Malformed Request",
                        "schema": {
                            "$ref": "#/definitions/models.GenericErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.GenericErrorResponse"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Returns pong",
                "tags": [
                    "example"
                ],
                "summary": "Ping example",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.EmailSignInRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "Strongpassword123"
                }
            }
        },
        "models.ErrorResult": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Error message"
                },
                "status": {
                    "type": "string",
                    "example": "Error"
                }
            }
        },
        "models.GenericErrorResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/models.ErrorResult"
                }
            }
        },
        "models.GenericSuccessResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/models.SuccessResult"
                }
            }
        },
        "models.SignInSuccessResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/models.SignInSuccessResult"
                }
            }
        },
        "models.SignInSuccessResult": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.TokenResponse"
                },
                "message": {
                    "type": "string",
                    "example": "Success message"
                },
                "status": {
                    "type": "string",
                    "example": "Success"
                }
            }
        },
        "models.SignUpRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "Strongpassword123"
                }
            }
        },
        "models.SuccessResult": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Success message"
                },
                "status": {
                    "type": "string",
                    "example": "Success"
                }
            }
        },
        "models.TokenResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
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
	Host:             "localhost:5000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Blockstracker",
	Description:      "Blockstracker API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
