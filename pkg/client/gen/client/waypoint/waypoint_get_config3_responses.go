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

// WaypointGetConfig3Reader is a Reader for the WaypointGetConfig3 structure.
type WaypointGetConfig3Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointGetConfig3Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointGetConfig3OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointGetConfig3Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointGetConfig3OK creates a WaypointGetConfig3OK with default headers values
func NewWaypointGetConfig3OK() *WaypointGetConfig3OK {
	return &WaypointGetConfig3OK{}
}

/*
WaypointGetConfig3OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointGetConfig3OK struct {
	Payload *models.HashicorpWaypointConfigGetResponse
}

// IsSuccess returns true when this waypoint get config3 o k response has a 2xx status code
func (o *WaypointGetConfig3OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint get config3 o k response has a 3xx status code
func (o *WaypointGetConfig3OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint get config3 o k response has a 4xx status code
func (o *WaypointGetConfig3OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint get config3 o k response has a 5xx status code
func (o *WaypointGetConfig3OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint get config3 o k response a status code equal to that given
func (o *WaypointGetConfig3OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointGetConfig3OK) Error() string {
	return fmt.Sprintf("[GET /project/{application.project}/{application.application}/config][%d] waypointGetConfig3OK  %+v", 200, o.Payload)
}

func (o *WaypointGetConfig3OK) String() string {
	return fmt.Sprintf("[GET /project/{application.project}/{application.application}/config][%d] waypointGetConfig3OK  %+v", 200, o.Payload)
}

func (o *WaypointGetConfig3OK) GetPayload() *models.HashicorpWaypointConfigGetResponse {
	return o.Payload
}

func (o *WaypointGetConfig3OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointConfigGetResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointGetConfig3Default creates a WaypointGetConfig3Default with default headers values
func NewWaypointGetConfig3Default(code int) *WaypointGetConfig3Default {
	return &WaypointGetConfig3Default{
		_statusCode: code,
	}
}

/*
WaypointGetConfig3Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointGetConfig3Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint get config3 default response
func (o *WaypointGetConfig3Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint get config3 default response has a 2xx status code
func (o *WaypointGetConfig3Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint get config3 default response has a 3xx status code
func (o *WaypointGetConfig3Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint get config3 default response has a 4xx status code
func (o *WaypointGetConfig3Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint get config3 default response has a 5xx status code
func (o *WaypointGetConfig3Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint get config3 default response a status code equal to that given
func (o *WaypointGetConfig3Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointGetConfig3Default) Error() string {
	return fmt.Sprintf("[GET /project/{application.project}/{application.application}/config][%d] Waypoint_GetConfig3 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetConfig3Default) String() string {
	return fmt.Sprintf("[GET /project/{application.project}/{application.application}/config][%d] Waypoint_GetConfig3 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetConfig3Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointGetConfig3Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
