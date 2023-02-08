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

// WaypointGetRelease2Reader is a Reader for the WaypointGetRelease2 structure.
type WaypointGetRelease2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointGetRelease2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointGetRelease2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointGetRelease2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointGetRelease2OK creates a WaypointGetRelease2OK with default headers values
func NewWaypointGetRelease2OK() *WaypointGetRelease2OK {
	return &WaypointGetRelease2OK{}
}

/*
WaypointGetRelease2OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointGetRelease2OK struct {
	Payload *models.HashicorpWaypointRelease
}

// IsSuccess returns true when this waypoint get release2 o k response has a 2xx status code
func (o *WaypointGetRelease2OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint get release2 o k response has a 3xx status code
func (o *WaypointGetRelease2OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint get release2 o k response has a 4xx status code
func (o *WaypointGetRelease2OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint get release2 o k response has a 5xx status code
func (o *WaypointGetRelease2OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint get release2 o k response a status code equal to that given
func (o *WaypointGetRelease2OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointGetRelease2OK) Error() string {
	return fmt.Sprintf("[GET /project/{ref.sequence.application.project}/application/{ref.sequence.application.application}/release/{ref.sequence.number}][%d] waypointGetRelease2OK  %+v", 200, o.Payload)
}

func (o *WaypointGetRelease2OK) String() string {
	return fmt.Sprintf("[GET /project/{ref.sequence.application.project}/application/{ref.sequence.application.application}/release/{ref.sequence.number}][%d] waypointGetRelease2OK  %+v", 200, o.Payload)
}

func (o *WaypointGetRelease2OK) GetPayload() *models.HashicorpWaypointRelease {
	return o.Payload
}

func (o *WaypointGetRelease2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointRelease)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointGetRelease2Default creates a WaypointGetRelease2Default with default headers values
func NewWaypointGetRelease2Default(code int) *WaypointGetRelease2Default {
	return &WaypointGetRelease2Default{
		_statusCode: code,
	}
}

/*
WaypointGetRelease2Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointGetRelease2Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint get release2 default response
func (o *WaypointGetRelease2Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint get release2 default response has a 2xx status code
func (o *WaypointGetRelease2Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint get release2 default response has a 3xx status code
func (o *WaypointGetRelease2Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint get release2 default response has a 4xx status code
func (o *WaypointGetRelease2Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint get release2 default response has a 5xx status code
func (o *WaypointGetRelease2Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint get release2 default response a status code equal to that given
func (o *WaypointGetRelease2Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointGetRelease2Default) Error() string {
	return fmt.Sprintf("[GET /project/{ref.sequence.application.project}/application/{ref.sequence.application.application}/release/{ref.sequence.number}][%d] Waypoint_GetRelease2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetRelease2Default) String() string {
	return fmt.Sprintf("[GET /project/{ref.sequence.application.project}/application/{ref.sequence.application.application}/release/{ref.sequence.number}][%d] Waypoint_GetRelease2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetRelease2Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointGetRelease2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
