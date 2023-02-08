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

// WaypointGetServerConfigReader is a Reader for the WaypointGetServerConfig structure.
type WaypointGetServerConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointGetServerConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointGetServerConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointGetServerConfigDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointGetServerConfigOK creates a WaypointGetServerConfigOK with default headers values
func NewWaypointGetServerConfigOK() *WaypointGetServerConfigOK {
	return &WaypointGetServerConfigOK{}
}

/*
WaypointGetServerConfigOK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointGetServerConfigOK struct {
	Payload *models.HashicorpWaypointGetServerConfigResponse
}

// IsSuccess returns true when this waypoint get server config o k response has a 2xx status code
func (o *WaypointGetServerConfigOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint get server config o k response has a 3xx status code
func (o *WaypointGetServerConfigOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint get server config o k response has a 4xx status code
func (o *WaypointGetServerConfigOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint get server config o k response has a 5xx status code
func (o *WaypointGetServerConfigOK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint get server config o k response a status code equal to that given
func (o *WaypointGetServerConfigOK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointGetServerConfigOK) Error() string {
	return fmt.Sprintf("[GET /server/config][%d] waypointGetServerConfigOK  %+v", 200, o.Payload)
}

func (o *WaypointGetServerConfigOK) String() string {
	return fmt.Sprintf("[GET /server/config][%d] waypointGetServerConfigOK  %+v", 200, o.Payload)
}

func (o *WaypointGetServerConfigOK) GetPayload() *models.HashicorpWaypointGetServerConfigResponse {
	return o.Payload
}

func (o *WaypointGetServerConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointGetServerConfigResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointGetServerConfigDefault creates a WaypointGetServerConfigDefault with default headers values
func NewWaypointGetServerConfigDefault(code int) *WaypointGetServerConfigDefault {
	return &WaypointGetServerConfigDefault{
		_statusCode: code,
	}
}

/*
WaypointGetServerConfigDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointGetServerConfigDefault struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint get server config default response
func (o *WaypointGetServerConfigDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint get server config default response has a 2xx status code
func (o *WaypointGetServerConfigDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint get server config default response has a 3xx status code
func (o *WaypointGetServerConfigDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint get server config default response has a 4xx status code
func (o *WaypointGetServerConfigDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint get server config default response has a 5xx status code
func (o *WaypointGetServerConfigDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint get server config default response a status code equal to that given
func (o *WaypointGetServerConfigDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointGetServerConfigDefault) Error() string {
	return fmt.Sprintf("[GET /server/config][%d] Waypoint_GetServerConfig default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetServerConfigDefault) String() string {
	return fmt.Sprintf("[GET /server/config][%d] Waypoint_GetServerConfig default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetServerConfigDefault) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointGetServerConfigDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
