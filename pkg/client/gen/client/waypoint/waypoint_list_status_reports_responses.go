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

// WaypointListStatusReportsReader is a Reader for the WaypointListStatusReports structure.
type WaypointListStatusReportsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointListStatusReportsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointListStatusReportsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointListStatusReportsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointListStatusReportsOK creates a WaypointListStatusReportsOK with default headers values
func NewWaypointListStatusReportsOK() *WaypointListStatusReportsOK {
	return &WaypointListStatusReportsOK{}
}

/*
WaypointListStatusReportsOK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointListStatusReportsOK struct {
	Payload *models.HashicorpWaypointListStatusReportsResponse
}

// IsSuccess returns true when this waypoint list status reports o k response has a 2xx status code
func (o *WaypointListStatusReportsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint list status reports o k response has a 3xx status code
func (o *WaypointListStatusReportsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint list status reports o k response has a 4xx status code
func (o *WaypointListStatusReportsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint list status reports o k response has a 5xx status code
func (o *WaypointListStatusReportsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint list status reports o k response a status code equal to that given
func (o *WaypointListStatusReportsOK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointListStatusReportsOK) Error() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/status-reports][%d] waypointListStatusReportsOK  %+v", 200, o.Payload)
}

func (o *WaypointListStatusReportsOK) String() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/status-reports][%d] waypointListStatusReportsOK  %+v", 200, o.Payload)
}

func (o *WaypointListStatusReportsOK) GetPayload() *models.HashicorpWaypointListStatusReportsResponse {
	return o.Payload
}

func (o *WaypointListStatusReportsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointListStatusReportsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointListStatusReportsDefault creates a WaypointListStatusReportsDefault with default headers values
func NewWaypointListStatusReportsDefault(code int) *WaypointListStatusReportsDefault {
	return &WaypointListStatusReportsDefault{
		_statusCode: code,
	}
}

/*
WaypointListStatusReportsDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointListStatusReportsDefault struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint list status reports default response
func (o *WaypointListStatusReportsDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint list status reports default response has a 2xx status code
func (o *WaypointListStatusReportsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint list status reports default response has a 3xx status code
func (o *WaypointListStatusReportsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint list status reports default response has a 4xx status code
func (o *WaypointListStatusReportsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint list status reports default response has a 5xx status code
func (o *WaypointListStatusReportsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint list status reports default response a status code equal to that given
func (o *WaypointListStatusReportsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointListStatusReportsDefault) Error() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/status-reports][%d] Waypoint_ListStatusReports default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListStatusReportsDefault) String() string {
	return fmt.Sprintf("[GET /project/{application.project}/application/{application.application}/status-reports][%d] Waypoint_ListStatusReports default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListStatusReportsDefault) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointListStatusReportsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
