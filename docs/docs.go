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
        "/auth/login": {
            "post": {
                "description": "Logs in a user with email and password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User logged in successfully, returns token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized - Invalid credentials",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Registers a new user with email, password, and name. Returns conflict error if email already exists.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registered",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": " Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/payments": {
            "post": {
                "description": "Marks a payment for a specific group and updates the group's payment status.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payments"
                ],
                "summary": "Marks a payment for a group.",
                "parameters": [
                    {
                        "description": "Mark Payment Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.MarkPaymentRequest"
                        }
                    },
                    {
                        "description": "groupId of the group for which the payment is marked",
                        "name": "groupId",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.MarkPaymentResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "404": {
                        "description": "Group not found or User not found",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "409": {
                        "description": "Payment already made",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/groups/": {
            "post": {
                "description": "Creates a group with the specified name and an associated bill, then adds the user as a member of the group.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Create a new group with an associated bill",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "CreateGroupWithBillRequest details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateGroupWithBillRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.CreateGroupWithBillResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/groups/member-groups": {
            "get": {
                "description": "Retrieves all groups associated with the authenticated user. Optionally filters the results by group status. If no status is provided, all groups will be returned.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "List groups the user belongs to",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The status of the groups to filter by. Valid values are 'PENDING' or 'DONE'",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with the list of groups.",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.ListMemberGroupsResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid status parameter provided.",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error.",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/groups/owned": {
            "get": {
                "description": "Fetches and returns a list of groups that are owned by the current user, including group members.",
                "tags": [
                    "groups"
                ],
                "summary": "List groups owned by the user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of groups owned by the user",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.ListOwnedGroupsResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/groups/{id}": {
            "delete": {
                "description": "Deletes a group identified by the specified group ID if the user is the creator of the group. (owner only can do this operation)",
                "tags": [
                    "groups"
                ],
                "summary": "Delete a group by ID (NOT NEEDED AS OF NOW)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the group to be deleted",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.DeleteGroupResponse"
                        }
                    },
                    "404": {
                        "description": "Group Not Found",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/groups/{id}/users": {
            "post": {
                "description": "Adds members identified by their email addresses to a group if the user is the creator of the group.",
                "tags": [
                    "groups"
                ],
                "summary": "Add members to a group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the group to which members will be added",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "List of user email IDs to add to the group",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.AddUsersToGroupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success message",
                        "schema": {
                            "$ref": "#/definitions/dto.AddUsersToGroupResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "404": {
                        "description": "Group Not Found or Users Not Found",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/payments/pending-payments": {
            "get": {
                "description": "Fetches all pending payments for the current user that have not been paid yet, including group ID, group name, bill ID, and amount owed.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payments"
                ],
                "summary": "Retrieve Pending Payments",
                "responses": {
                    "200": {
                        "description": "Successful response containing the list of pending payments and total amount.",
                        "schema": {
                            "$ref": "#/definitions/dto.PendingPaymentsWithTotalResponse"
                        }
                    },
                    "404": {
                        "description": "No pending payments found for the user.",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error occurred while fetching pending payments.",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/report": {
            "post": {
                "description": "Generates and downloads a PDF report for the groups created by the user within a specified date range.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/pdf"
                ],
                "tags": [
                    "reports"
                ],
                "summary": "Download PDF report of user's groups",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "from date in the format YYYY-MM-DD",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "to date in the format YYYY-MM-DD",
                        "name": "to",
                        "in": "query"
                    },
                    {
                        "description": "GetGroupReportRequest details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.GetGroupReportRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PDF report generated and downloaded",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        },
        "/v1/report/{id}": {
            "get": {
                "description": "Generates a detailed PDF report for the group specified by its ID. The report includes group details, associated bills, and member history.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/pdf"
                ],
                "tags": [
                    "reports"
                ],
                "summary": "Generate a PDF report for a specific group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Group ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PDF report generated successfully",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "404": {
                        "description": "Group not found",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.AddUsersToGroupRequest": {
            "type": "object",
            "required": [
                "userEmailIds"
            ],
            "properties": {
                "userEmailIds": {
                    "description": "UserIds      []uint   ` + "`" + `json:\"userIds\"` + "`" + `",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "dto.AddUsersToGroupResponse": {
            "description": "Response model for adding users to a group.",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.CreateGroupWithBillRequest": {
            "type": "object",
            "required": [
                "bill",
                "groupName"
            ],
            "properties": {
                "bill": {
                    "type": "object",
                    "required": [
                        "amount",
                        "name"
                    ],
                    "properties": {
                        "amount": {
                            "type": "number"
                        },
                        "name": {
                            "type": "string"
                        }
                    }
                },
                "groupName": {
                    "type": "string"
                }
            }
        },
        "dto.CreateGroupWithBillResponse": {
            "description": "Response model for the creation of a group with an associated bill.",
            "type": "object",
            "properties": {
                "billId": {
                    "type": "integer"
                },
                "groupId": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.DeleteGroupResponse": {
            "description": "Response model for deleting a group.",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.GetGroupReportRequest": {
            "description": "Request model for generating a group report based on date range.",
            "type": "object",
            "properties": {
                "from": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                }
            }
        },
        "dto.ListMemberGroupsResponse": {
            "description": "Response model for listing groups the user belongs to, including group details and member information.",
            "type": "object",
            "properties": {
                "group": {
                    "$ref": "#/definitions/models.Group"
                },
                "members": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GroupMember"
                    }
                }
            }
        },
        "dto.ListOwnedGroupsResponse": {
            "description": "ListOwnedGroupsResponse is the response model for listing owned groups.",
            "type": "object",
            "properties": {
                "group": {
                    "$ref": "#/definitions/models.Group"
                },
                "members": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GroupMember"
                    }
                }
            }
        },
        "dto.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "jis@jish.com"
                },
                "password": {
                    "type": "string",
                    "example": "Passw0rd@123"
                }
            }
        },
        "dto.MarkPaymentRequest": {
            "description": "Mark a payment for a specific group",
            "type": "object",
            "required": [
                "groupId"
            ],
            "properties": {
                "groupId": {
                    "type": "integer"
                },
                "remarks": {
                    "description": "Optional* remarks for the payment",
                    "type": "string"
                }
            }
        },
        "dto.MarkPaymentResponse": {
            "description": "Response for marking a payment",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.PendingPayments": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "billId": {
                    "type": "integer"
                },
                "groupId": {
                    "type": "integer"
                },
                "groupName": {
                    "type": "string"
                }
            }
        },
        "dto.PendingPaymentsWithTotalResponse": {
            "description": "Response model for listing pending payments with total amount.",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "pendingPayments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.PendingPayments"
                    }
                },
                "totalAmount": {
                    "type": "number"
                }
            }
        },
        "dto.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                }
            }
        },
        "errors.Error": {
            "description": "Error model for handling errors.",
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Bill": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "Total amount",
                    "type": "number"
                },
                "completed": {
                    "description": "Overall bill payment status",
                    "type": "boolean"
                },
                "groupId": {
                    "description": "Reference to the associated group",
                    "type": "integer"
                },
                "history": {
                    "description": "Bill payment history",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.BillHistory"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.BillHistory": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "Amount related to this history entry",
                    "type": "number"
                },
                "billId": {
                    "description": "Automatically inferred foreign key",
                    "type": "integer"
                },
                "createdAt": {
                    "description": "Auto-create timestamp",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "paidAt": {
                    "description": "Time of payment",
                    "type": "string"
                },
                "paidBy": {
                    "description": "User who made the payment",
                    "type": "string"
                }
            }
        },
        "models.Group": {
            "type": "object",
            "properties": {
                "bill": {
                    "$ref": "#/definitions/models.Bill"
                },
                "billId": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "createdBy": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "paidAmount": {
                    "type": "number"
                },
                "perUserSplitAmount": {
                    "type": "number"
                },
                "status": {
                    "type": "string"
                },
                "totalAmount": {
                    "type": "number"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "models.GroupMember": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "groupId": {
                    "type": "integer"
                },
                "hasPaid": {
                    "description": "Tracks if the member has paid",
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "remarks": {
                    "type": "string"
                },
                "splitAmount": {
                    "type": "number"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "SplitWise API",
	Description:      "This is an API for managing splits.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
