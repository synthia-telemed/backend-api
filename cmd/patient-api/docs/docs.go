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
        "/appointment": {
            "get": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Appointment"
                ],
                "summary": "Get list of appointment of the patient",
                "responses": {
                    "200": {
                        "description": "List of appointment group by status",
                        "schema": {
                            "$ref": "#/definitions/hospital.CategorizedAppointment"
                        }
                    },
                    "400": {
                        "description": "Patient not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/appointment/{appointmentID}": {
            "get": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Appointment"
                ],
                "summary": "Get an appointment detail by appointment ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the appointment",
                        "name": "appointmentID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "An appointment detail",
                        "schema": {
                            "$ref": "#/definitions/handler.GetAppointmentResponse"
                        }
                    },
                    "400": {
                        "description": "appointmentID is invalid",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "The patient doesn't own the appointment",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Appointment not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/appointment/{appointmentID}/roomID": {
            "get": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Appointment"
                ],
                "summary": "Get room ID of the appointment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the appointment",
                        "name": "appointmentID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Room ID for the appointment",
                        "schema": {
                            "$ref": "#/definitions/handler.GetAppointmentRoomIDResponse"
                        }
                    },
                    "400": {
                        "description": "appointmentID is invalid",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "The patient doesn't own the appointment",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "RoomID of the appointment not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
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
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Provided credential is not in the hospital system",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
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
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/measurement/blood-pressure": {
            "post": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Measurement"
                ],
                "summary": "Record blood pressure",
                "parameters": [
                    {
                        "description": "Blood pressure information",
                        "name": "BloodPressureRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.BloodPressureRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/datastore.BloodPressure"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/measurement/glucose": {
            "post": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Measurement"
                ],
                "summary": "Record glucose level",
                "parameters": [
                    {
                        "description": "Glucose level information",
                        "name": "GlucoseRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.GlucoseRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/datastore.Glucose"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
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
                    },
                    {
                        "JWSToken": []
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
                                "$ref": "#/definitions/datastore.CreditCard"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Payment"
                ],
                "summary": "Add new credit card",
                "parameters": [
                    {
                        "description": "Token from Omise and name of credit card",
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
                        "description": "Limited number of credit cards is reached",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/payment/credit-card/{cardID}": {
            "delete": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Payment"
                ],
                "summary": "Delete saved credit card",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the credit card to delete",
                        "name": "cardID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Invalid credit card ID",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Patient doesn't own the specified credit card",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Credit card not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Payment"
                ],
                "summary": "Set isDefault status of credit card",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the credit card to set isDefault",
                        "name": "cardID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Invalid credit card ID",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Patient doesn't own the specified credit card",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Credit card not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/payment/pay/{invoiceID}/credit-card/{cardID}": {
            "post": {
                "security": [
                    {
                        "UserID": []
                    },
                    {
                        "JWSToken": []
                    }
                ],
                "tags": [
                    "Payment"
                ],
                "summary": "Pay invoice with credit card method",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the credit card to be charged",
                        "name": "cardID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "ID of the invoice to pay",
                        "name": "invoiceID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Payment information",
                        "schema": {
                            "$ref": "#/definitions/handler.PayInvoiceWithCreditCardResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid credit card ID or invoice ID",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Patient doesn't own the specified credit card or invoice",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Credit card or invoice not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "datastore.BloodPressure": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "date_time": {
                    "type": "string"
                },
                "deletedAt": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "diastolic": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "patient_id": {
                    "type": "integer"
                },
                "pulse": {
                    "type": "integer"
                },
                "systolic": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "datastore.CreditCard": {
            "type": "object",
            "properties": {
                "brand": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "expiry": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_default": {
                    "type": "boolean"
                },
                "last_4_digits": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "patient_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "datastore.Glucose": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "date_time": {
                    "type": "string"
                },
                "deletedAt": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "id": {
                    "type": "integer"
                },
                "is_before_meal": {
                    "type": "boolean"
                },
                "patient_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "datastore.Payment": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "credit_card": {
                    "$ref": "#/definitions/datastore.CreditCard"
                },
                "credit_card_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "invoice_id": {
                    "type": "integer"
                },
                "method": {
                    "type": "string"
                },
                "paid_at": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "gorm.DeletedAt": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        },
        "handler.AddCreditCardRequest": {
            "type": "object",
            "required": [
                "card_token"
            ],
            "properties": {
                "card_token": {
                    "type": "string"
                },
                "is_default": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.BloodPressureRequest": {
            "type": "object",
            "required": [
                "date_time",
                "diastolic",
                "pulse",
                "systolic"
            ],
            "properties": {
                "date_time": {
                    "type": "string"
                },
                "diastolic": {
                    "type": "integer"
                },
                "pulse": {
                    "type": "integer"
                },
                "systolic": {
                    "type": "integer"
                }
            }
        },
        "handler.GetAppointmentResponse": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string"
                },
                "doctor": {
                    "$ref": "#/definitions/hospital.DoctorOverview"
                },
                "duration": {
                    "type": "number"
                },
                "end_date_time": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "invoice": {
                    "$ref": "#/definitions/hospital.Invoice"
                },
                "next_appointment": {
                    "type": "string"
                },
                "patient_id": {
                    "type": "string"
                },
                "payment": {
                    "$ref": "#/definitions/datastore.Payment"
                },
                "prescriptions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/hospital.Prescription"
                    }
                },
                "start_date_time": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "handler.GetAppointmentRoomIDResponse": {
            "type": "object",
            "properties": {
                "room_id": {
                    "type": "string"
                }
            }
        },
        "handler.GlucoseRequest": {
            "type": "object",
            "required": [
                "date_time",
                "is_before_meal",
                "value"
            ],
            "properties": {
                "date_time": {
                    "type": "string"
                },
                "is_before_meal": {
                    "type": "boolean"
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "handler.PayInvoiceWithCreditCardResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "credit_card": {
                    "$ref": "#/definitions/datastore.CreditCard"
                },
                "credit_card_id": {
                    "type": "integer"
                },
                "failure_message": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "invoice_id": {
                    "type": "integer"
                },
                "method": {
                    "type": "string"
                },
                "paid_at": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
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
                "expired_at": {
                    "type": "string"
                },
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
        "hospital.AppointmentOverview": {
            "type": "object",
            "properties": {
                "doctor": {
                    "$ref": "#/definitions/hospital.DoctorOverview"
                },
                "end_date_time": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "patient": {
                    "$ref": "#/definitions/hospital.PatientOverview"
                },
                "start_date_time": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "hospital.CategorizedAppointment": {
            "type": "object",
            "properties": {
                "cancelled": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/hospital.AppointmentOverview"
                    }
                },
                "completed": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/hospital.AppointmentOverview"
                    }
                },
                "scheduled": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/hospital.AppointmentOverview"
                    }
                }
            }
        },
        "hospital.DoctorOverview": {
            "type": "object",
            "properties": {
                "full_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "position": {
                    "type": "string"
                },
                "profile_pic_url": {
                    "type": "string"
                }
            }
        },
        "hospital.Invoice": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "invoice_discounts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/hospital.InvoiceDiscount"
                    }
                },
                "invoice_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/hospital.InvoiceItem"
                    }
                },
                "paid": {
                    "type": "boolean"
                },
                "total": {
                    "type": "number"
                }
            }
        },
        "hospital.InvoiceDiscount": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "hospital.InvoiceItem": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "hospital.PatientOverview": {
            "type": "object",
            "properties": {
                "full_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "hospital.Prescription": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "picture_url": {
                    "type": "string"
                }
            }
        },
        "server.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWSToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
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
	BasePath:         "/patient/api",
	Schemes:          []string{},
	Title:            "Synthia Patient Backend API",
	Description:      "This is a Synthia patient backend API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
