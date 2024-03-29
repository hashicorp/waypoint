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

// WaypointListBuildsReader is a Reader for the WaypointListBuilds structure.
type WaypointListBuildsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointListBuildsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointListBuildsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointListBuildsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointListBuildsOK creates a WaypointListBuildsOK with default headers values
func NewWaypointListBuildsOK() *WaypointListBuildsOK {
	return &WaypointListBuildsOK{}
}

/*
WaypointListBuildsOK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointListBuildsOK struct {
	Payload *models.HashicorpWaypointListBuildsResponse
}

// IsSuccess returns true when this waypoint list builds o k response has a 2xx status code
func (o *WaypointListBuildsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint list builds o k response has a 3xx status code
func (o *WaypointListBuildsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint list builds o k response has a 4xx status code
func (o *WaypointListBuildsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint list builds o k response has a 5xx status code
func (o *WaypointListBuildsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint list builds o k response a status code equal to that given
func (o *WaypointListBuildsOK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointListBuildsOK) Error() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/builds][%d] waypointListBuildsOK  %+v", 200, o.Payload)
}

func (o *WaypointListBuildsOK) String() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/builds][%d] waypointListBuildsOK  %+v", 200, o.Payload)
}

func (o *WaypointListBuildsOK) GetPayload() *models.HashicorpWaypointListBuildsResponse {
	return o.Payload
}

func (o *WaypointListBuildsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointListBuildsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointListBuildsDefault creates a WaypointListBuildsDefault with default headers values
func NewWaypointListBuildsDefault(code int) *WaypointListBuildsDefault {
	return &WaypointListBuildsDefault{
		_statusCode: code,
	}
}

/*
WaypointListBuildsDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointListBuildsDefault struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint list builds default response
func (o *WaypointListBuildsDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint list builds default response has a 2xx status code
func (o *WaypointListBuildsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint list builds default response has a 3xx status code
func (o *WaypointListBuildsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint list builds default response has a 4xx status code
func (o *WaypointListBuildsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint list builds default response has a 5xx status code
func (o *WaypointListBuildsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint list builds default response a status code equal to that given
func (o *WaypointListBuildsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointListBuildsDefault) Error() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/builds][%d] Waypoint_ListBuilds default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListBuildsDefault) String() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/builds][%d] Waypoint_ListBuilds default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListBuildsDefault) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointListBuildsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
