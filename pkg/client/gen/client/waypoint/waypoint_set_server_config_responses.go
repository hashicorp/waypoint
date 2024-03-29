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

// WaypointSetServerConfigReader is a Reader for the WaypointSetServerConfig structure.
type WaypointSetServerConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointSetServerConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointSetServerConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointSetServerConfigDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointSetServerConfigOK creates a WaypointSetServerConfigOK with default headers values
func NewWaypointSetServerConfigOK() *WaypointSetServerConfigOK {
	return &WaypointSetServerConfigOK{}
}

/*
WaypointSetServerConfigOK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointSetServerConfigOK struct {
	Payload interface{}
}

// IsSuccess returns true when this waypoint set server config o k response has a 2xx status code
func (o *WaypointSetServerConfigOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint set server config o k response has a 3xx status code
func (o *WaypointSetServerConfigOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint set server config o k response has a 4xx status code
func (o *WaypointSetServerConfigOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint set server config o k response has a 5xx status code
func (o *WaypointSetServerConfigOK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint set server config o k response a status code equal to that given
func (o *WaypointSetServerConfigOK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointSetServerConfigOK) Error() string {
	return fmt.Sprintf("[POST /server/config][%d] waypointSetServerConfigOK  %+v", 200, o.Payload)
}

func (o *WaypointSetServerConfigOK) String() string {
	return fmt.Sprintf("[POST /server/config][%d] waypointSetServerConfigOK  %+v", 200, o.Payload)
}

func (o *WaypointSetServerConfigOK) GetPayload() interface{} {
	return o.Payload
}

func (o *WaypointSetServerConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointSetServerConfigDefault creates a WaypointSetServerConfigDefault with default headers values
func NewWaypointSetServerConfigDefault(code int) *WaypointSetServerConfigDefault {
	return &WaypointSetServerConfigDefault{
		_statusCode: code,
	}
}

/*
WaypointSetServerConfigDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointSetServerConfigDefault struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint set server config default response
func (o *WaypointSetServerConfigDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint set server config default response has a 2xx status code
func (o *WaypointSetServerConfigDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint set server config default response has a 3xx status code
func (o *WaypointSetServerConfigDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint set server config default response has a 4xx status code
func (o *WaypointSetServerConfigDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint set server config default response has a 5xx status code
func (o *WaypointSetServerConfigDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint set server config default response a status code equal to that given
func (o *WaypointSetServerConfigDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointSetServerConfigDefault) Error() string {
	return fmt.Sprintf("[POST /server/config][%d] Waypoint_SetServerConfig default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointSetServerConfigDefault) String() string {
	return fmt.Sprintf("[POST /server/config][%d] Waypoint_SetServerConfig default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointSetServerConfigDefault) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointSetServerConfigDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
