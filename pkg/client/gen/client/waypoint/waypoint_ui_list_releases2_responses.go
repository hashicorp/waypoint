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

// WaypointUIListReleases2Reader is a Reader for the WaypointUIListReleases2 structure.
type WaypointUIListReleases2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointUIListReleases2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointUIListReleases2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointUIListReleases2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointUIListReleases2OK creates a WaypointUIListReleases2OK with default headers values
func NewWaypointUIListReleases2OK() *WaypointUIListReleases2OK {
	return &WaypointUIListReleases2OK{}
}

/*
WaypointUIListReleases2OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointUIListReleases2OK struct {
	Payload *models.HashicorpWaypointUIListReleasesResponse
}

// IsSuccess returns true when this waypoint Ui list releases2 o k response has a 2xx status code
func (o *WaypointUIListReleases2OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint Ui list releases2 o k response has a 3xx status code
func (o *WaypointUIListReleases2OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint Ui list releases2 o k response has a 4xx status code
func (o *WaypointUIListReleases2OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint Ui list releases2 o k response has a 5xx status code
func (o *WaypointUIListReleases2OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint Ui list releases2 o k response a status code equal to that given
func (o *WaypointUIListReleases2OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointUIListReleases2OK) Error() string {
	return fmt.Sprintf("[GET /ui/releases/application/{application.application}][%d] waypointUiListReleases2OK  %+v", 200, o.Payload)
}

func (o *WaypointUIListReleases2OK) String() string {
	return fmt.Sprintf("[GET /ui/releases/application/{application.application}][%d] waypointUiListReleases2OK  %+v", 200, o.Payload)
}

func (o *WaypointUIListReleases2OK) GetPayload() *models.HashicorpWaypointUIListReleasesResponse {
	return o.Payload
}

func (o *WaypointUIListReleases2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointUIListReleasesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointUIListReleases2Default creates a WaypointUIListReleases2Default with default headers values
func NewWaypointUIListReleases2Default(code int) *WaypointUIListReleases2Default {
	return &WaypointUIListReleases2Default{
		_statusCode: code,
	}
}

/*
WaypointUIListReleases2Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointUIListReleases2Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint UI list releases2 default response
func (o *WaypointUIListReleases2Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint UI list releases2 default response has a 2xx status code
func (o *WaypointUIListReleases2Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint UI list releases2 default response has a 3xx status code
func (o *WaypointUIListReleases2Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint UI list releases2 default response has a 4xx status code
func (o *WaypointUIListReleases2Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint UI list releases2 default response has a 5xx status code
func (o *WaypointUIListReleases2Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint UI list releases2 default response a status code equal to that given
func (o *WaypointUIListReleases2Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointUIListReleases2Default) Error() string {
	return fmt.Sprintf("[GET /ui/releases/application/{application.application}][%d] Waypoint_UI_ListReleases2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointUIListReleases2Default) String() string {
	return fmt.Sprintf("[GET /ui/releases/application/{application.application}][%d] Waypoint_UI_ListReleases2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointUIListReleases2Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointUIListReleases2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
