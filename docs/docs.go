// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
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
                "description": "Initiate auth process with government credential which will sent OTP to patient's phone number",
                "tags": [
                    "Auth"
                ],
                "summary": "Start signing-in with government credential",
                "parameters": [
                    {
                        "description": "Patient government credential (Passport ID or National ID)",
                        "name": "SigninRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SigninRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "OTP is sent to patient's phone number",
                        "schema": {
                            "$ref": "#/definitions/handler.SigninResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Provided credential is not in the hospital system",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify": {
            "post": {
                "description": "Complete auth process with OTP verification. It will return token if verification success.",
                "tags": [
                    "Auth"
                ],
                "summary": "Verify OTP and get token",
                "parameters": [
                    {
                        "description": "OTP that is sent to patient's phone number",
                        "name": "VerifyOTPRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.VerifyOTPRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "JWS Token for later use",
                        "schema": {
                            "$ref": "#/definitions/handler.VerifyOTPResponse"
                        }
                    },
                    "400": {
                        "description": "OTP is invalid or expired",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/payment/credit-card": {
            "get": {
                "security": [
                    {
                        "UserID": []
                    }
                ],
                "tags": [
                    "Payment"
                ],
                "summary": "Get lists of saved credit cards",
                "responses": {
                    "200": {
                        "description": "List of saved cards",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/payment.Card"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "UserID": []
                    }
                ],
                "tags": [
                    "Payment"
                ],
                "summary": "Add new credit card",
                "parameters": [
                    {
                        "description": "Token from Omise",
                        "name": "AddCreditCardRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.AddCreditCardRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Failed to add credit card",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.AddCreditCardRequest": {
            "type": "object",
            "required": [
                "card_token"
            ],
            "properties": {
                "card_token": {
                    "type": "string"
                }
            }
        },
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.SigninRequest": {
            "type": "object",
            "required": [
                "credential"
            ],
            "properties": {
                "credential": {
                    "type": "string"
                }
            }
        },
        "handler.SigninResponse": {
            "type": "object",
            "properties": {
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "handler.VerifyOTPRequest": {
            "type": "object",
            "required": [
                "otp"
            ],
            "properties": {
                "otp": {
                    "type": "string"
                }
            }
        },
        "handler.VerifyOTPResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "payment.Card": {
            "type": "object",
            "properties": {
                "brand": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_digits": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "UserID": {
            "type": "apiKey",
            "name": "X-USER-ID",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Synthia Patient Backend API",
	Description:      "This is a Synthia patient backend API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
