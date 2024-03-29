// Code generated by go-swagger; DO NOT EDIT.

package waypoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/waypoint/pkg/client/gen/models"
)

// WaypointListTriggers2Reader is a Reader for the WaypointListTriggers2 structure.
type WaypointListTriggers2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointListTriggers2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointListTriggers2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointListTriggers2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointListTriggers2OK creates a WaypointListTriggers2OK with default headers values
func NewWaypointListTriggers2OK() *WaypointListTriggers2OK {
	return &WaypointListTriggers2OK{}
}

/*
WaypointListTriggers2OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointListTriggers2OK struct {
	Payload *models.HashicorpWaypointListTriggerResponse
}

// IsSuccess returns true when this waypoint list triggers2 o k response has a 2xx status code
func (o *WaypointListTriggers2OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint list triggers2 o k response has a 3xx status code
func (o *WaypointListTriggers2OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint list triggers2 o k response has a 4xx status code
func (o *WaypointListTriggers2OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint list triggers2 o k response has a 5xx status code
func (o *WaypointListTriggers2OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint list triggers2 o k response a status code equal to that given
func (o *WaypointListTriggers2OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointListTriggers2OK) Error() string {
	return fmt.Sprintf("[GET /project/{project.project}/triggers][%d] waypointListTriggers2OK  %+v", 200, o.Payload)
}

func (o *WaypointListTriggers2OK) String() string {
	return fmt.Sprintf("[GET /project/{project.project}/triggers][%d] waypointListTriggers2OK  %+v", 200, o.Payload)
}

func (o *WaypointListTriggers2OK) GetPayload() *models.HashicorpWaypointListTriggerResponse {
	return o.Payload
}

func (o *WaypointListTriggers2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointListTriggerResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointListTriggers2Default creates a WaypointListTriggers2Default with default headers values
func NewWaypointListTriggers2Default(code int) *WaypointListTriggers2Default {
	return &WaypointListTriggers2Default{
		_statusCode: code,
	}
}

/*
WaypointListTriggers2Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointListTriggers2Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint list triggers2 default response
func (o *WaypointListTriggers2Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint list triggers2 default response has a 2xx status code
func (o *WaypointListTriggers2Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint list triggers2 default response has a 3xx status code
func (o *WaypointListTriggers2Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint list triggers2 default response has a 4xx status code
func (o *WaypointListTriggers2Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint list triggers2 default response has a 5xx status code
func (o *WaypointListTriggers2Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint list triggers2 default response a status code equal to that given
func (o *WaypointListTriggers2Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointListTriggers2Default) Error() string {
	return fmt.Sprintf("[GET /project/{project.project}/triggers][%d] Waypoint_ListTriggers2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListTriggers2Default) String() string {
	return fmt.Sprintf("[GET /project/{project.project}/triggers][%d] Waypoint_ListTriggers2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListTriggers2Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointListTriggers2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
