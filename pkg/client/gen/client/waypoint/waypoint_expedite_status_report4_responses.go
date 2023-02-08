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

// WaypointExpediteStatusReport4Reader is a Reader for the WaypointExpediteStatusReport4 structure.
type WaypointExpediteStatusReport4Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointExpediteStatusReport4Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointExpediteStatusReport4OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointExpediteStatusReport4Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointExpediteStatusReport4OK creates a WaypointExpediteStatusReport4OK with default headers values
func NewWaypointExpediteStatusReport4OK() *WaypointExpediteStatusReport4OK {
	return &WaypointExpediteStatusReport4OK{}
}

/*
WaypointExpediteStatusReport4OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointExpediteStatusReport4OK struct {
	Payload *models.HashicorpWaypointExpediteStatusReportResponse
}

// IsSuccess returns true when this waypoint expedite status report4 o k response has a 2xx status code
func (o *WaypointExpediteStatusReport4OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint expedite status report4 o k response has a 3xx status code
func (o *WaypointExpediteStatusReport4OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint expedite status report4 o k response has a 4xx status code
func (o *WaypointExpediteStatusReport4OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint expedite status report4 o k response has a 5xx status code
func (o *WaypointExpediteStatusReport4OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint expedite status report4 o k response a status code equal to that given
func (o *WaypointExpediteStatusReport4OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointExpediteStatusReport4OK) Error() string {
	return fmt.Sprintf("[PUT /project/{release.sequence.application.project}/application/{release.sequence.application.application}/release/{release.sequence.number}/status-report][%d] waypointExpediteStatusReport4OK  %+v", 200, o.Payload)
}

func (o *WaypointExpediteStatusReport4OK) String() string {
	return fmt.Sprintf("[PUT /project/{release.sequence.application.project}/application/{release.sequence.application.application}/release/{release.sequence.number}/status-report][%d] waypointExpediteStatusReport4OK  %+v", 200, o.Payload)
}

func (o *WaypointExpediteStatusReport4OK) GetPayload() *models.HashicorpWaypointExpediteStatusReportResponse {
	return o.Payload
}

func (o *WaypointExpediteStatusReport4OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointExpediteStatusReportResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointExpediteStatusReport4Default creates a WaypointExpediteStatusReport4Default with default headers values
func NewWaypointExpediteStatusReport4Default(code int) *WaypointExpediteStatusReport4Default {
	return &WaypointExpediteStatusReport4Default{
		_statusCode: code,
	}
}

/*
WaypointExpediteStatusReport4Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointExpediteStatusReport4Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint expedite status report4 default response
func (o *WaypointExpediteStatusReport4Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint expedite status report4 default response has a 2xx status code
func (o *WaypointExpediteStatusReport4Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint expedite status report4 default response has a 3xx status code
func (o *WaypointExpediteStatusReport4Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint expedite status report4 default response has a 4xx status code
func (o *WaypointExpediteStatusReport4Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint expedite status report4 default response has a 5xx status code
func (o *WaypointExpediteStatusReport4Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint expedite status report4 default response a status code equal to that given
func (o *WaypointExpediteStatusReport4Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointExpediteStatusReport4Default) Error() string {
	return fmt.Sprintf("[PUT /project/{release.sequence.application.project}/application/{release.sequence.application.application}/release/{release.sequence.number}/status-report][%d] Waypoint_ExpediteStatusReport4 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointExpediteStatusReport4Default) String() string {
	return fmt.Sprintf("[PUT /project/{release.sequence.application.project}/application/{release.sequence.application.application}/release/{release.sequence.number}/status-report][%d] Waypoint_ExpediteStatusReport4 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointExpediteStatusReport4Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointExpediteStatusReport4Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
