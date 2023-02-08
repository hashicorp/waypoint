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

// WaypointListJobs3Reader is a Reader for the WaypointListJobs3 structure.
type WaypointListJobs3Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointListJobs3Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointListJobs3OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointListJobs3Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointListJobs3OK creates a WaypointListJobs3OK with default headers values
func NewWaypointListJobs3OK() *WaypointListJobs3OK {
	return &WaypointListJobs3OK{}
}

/*
WaypointListJobs3OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointListJobs3OK struct {
	Payload *models.HashicorpWaypointListJobsResponse
}

// IsSuccess returns true when this waypoint list jobs3 o k response has a 2xx status code
func (o *WaypointListJobs3OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint list jobs3 o k response has a 3xx status code
func (o *WaypointListJobs3OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint list jobs3 o k response has a 4xx status code
func (o *WaypointListJobs3OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint list jobs3 o k response has a 5xx status code
func (o *WaypointListJobs3OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint list jobs3 o k response a status code equal to that given
func (o *WaypointListJobs3OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointListJobs3OK) Error() string {
	return fmt.Sprintf("[GET /jobs/project/{project.project}][%d] waypointListJobs3OK  %+v", 200, o.Payload)
}

func (o *WaypointListJobs3OK) String() string {
	return fmt.Sprintf("[GET /jobs/project/{project.project}][%d] waypointListJobs3OK  %+v", 200, o.Payload)
}

func (o *WaypointListJobs3OK) GetPayload() *models.HashicorpWaypointListJobsResponse {
	return o.Payload
}

func (o *WaypointListJobs3OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointListJobsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointListJobs3Default creates a WaypointListJobs3Default with default headers values
func NewWaypointListJobs3Default(code int) *WaypointListJobs3Default {
	return &WaypointListJobs3Default{
		_statusCode: code,
	}
}

/*
WaypointListJobs3Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointListJobs3Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint list jobs3 default response
func (o *WaypointListJobs3Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint list jobs3 default response has a 2xx status code
func (o *WaypointListJobs3Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint list jobs3 default response has a 3xx status code
func (o *WaypointListJobs3Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint list jobs3 default response has a 4xx status code
func (o *WaypointListJobs3Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint list jobs3 default response has a 5xx status code
func (o *WaypointListJobs3Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint list jobs3 default response a status code equal to that given
func (o *WaypointListJobs3Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointListJobs3Default) Error() string {
	return fmt.Sprintf("[GET /jobs/project/{project.project}][%d] Waypoint_ListJobs3 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListJobs3Default) String() string {
	return fmt.Sprintf("[GET /jobs/project/{project.project}][%d] Waypoint_ListJobs3 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointListJobs3Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointListJobs3Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
