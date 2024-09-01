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
        "/api/categories/get": {
            "post": {
                "description": "gets categories",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Categories"
                ],
                "summary": "Получение категорий",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCategoriesDto"
                        }
                    }
                }
            }
        },
        "/api/products/get": {
            "post": {
                "description": "gets products with pagination and filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Products"
                ],
                "summary": "Получение товаров",
                "parameters": [
                    {
                        "description": "Получение товаров",
                        "name": "GetProductsRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetProductsDto"
                        }
                    }
                }
            }
        },
        "/api/users/login": {
            "post": {
                "description": "login the user and returns tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "Аутентификация пользователя",
                        "name": "LoginData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.LoginInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.TokenResponse"
                        }
                    }
                }
            }
        },
        "/api/users/reauth": {
            "get": {
                "description": "login the user and returns tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Обновление токена",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Рефреш",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ReauthResponse"
                        }
                    }
                }
            }
        },
        "/api/users/register": {
            "post": {
                "description": "register the user and returns tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Регистрация нового пользователя",
                "parameters": [
                    {
                        "description": "Регистрация нового пользователя",
                        "name": "RegisterData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.RegisterInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.TokenResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.GetCategoriesDto": {
            "type": "object",
            "properties": {
                "categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Category"
                    }
                }
            }
        },
        "dto.GetProductsDto": {
            "type": "object",
            "properties": {
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Product"
                    }
                }
            }
        },
        "dto.ReauthResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "dto.Request": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/filters.FilterBody"
                }
            }
        },
        "dto.TokenResponse": {
            "type": "object",
            "properties": {
                "refresh_token": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "filters.FilterBody": {
            "type": "object",
            "required": [
                "filter",
                "order",
                "skip",
                "take"
            ],
            "properties": {
                "filter": {},
                "order": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/filters.OrderBy"
                    }
                },
                "skip": {
                    "type": "integer"
                },
                "take": {
                    "type": "integer"
                }
            }
        },
        "filters.OrderBy": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "boolean"
                },
                "field": {
                    "type": "string"
                }
            }
        },
        "models.Category": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Product"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Product": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "imagesIds": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "maker": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "service.LoginInput": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 10
                }
            }
        },
        "service.RegisterInput": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Tag Api for shop",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
