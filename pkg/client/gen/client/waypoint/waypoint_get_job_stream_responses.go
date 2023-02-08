// Code generated by go-swagger; DO NOT EDIT.

package waypoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hashicorp/waypoint/pkg/client/gen/models"
)

// WaypointGetJobStreamReader is a Reader for the WaypointGetJobStream structure.
type WaypointGetJobStreamReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointGetJobStreamReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointGetJobStreamOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointGetJobStreamDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointGetJobStreamOK creates a WaypointGetJobStreamOK with default headers values
func NewWaypointGetJobStreamOK() *WaypointGetJobStreamOK {
	return &WaypointGetJobStreamOK{}
}

/*
WaypointGetJobStreamOK describes a response with status code 200, with default header values.

A successful response.(streaming responses)
*/
type WaypointGetJobStreamOK struct {
	Payload *WaypointGetJobStreamOKBody
}

// IsSuccess returns true when this waypoint get job stream o k response has a 2xx status code
func (o *WaypointGetJobStreamOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint get job stream o k response has a 3xx status code
func (o *WaypointGetJobStreamOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint get job stream o k response has a 4xx status code
func (o *WaypointGetJobStreamOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint get job stream o k response has a 5xx status code
func (o *WaypointGetJobStreamOK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint get job stream o k response a status code equal to that given
func (o *WaypointGetJobStreamOK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointGetJobStreamOK) Error() string {
	return fmt.Sprintf("[GET /jobs/stream/{job_id}][%d] waypointGetJobStreamOK  %+v", 200, o.Payload)
}

func (o *WaypointGetJobStreamOK) String() string {
	return fmt.Sprintf("[GET /jobs/stream/{job_id}][%d] waypointGetJobStreamOK  %+v", 200, o.Payload)
}

func (o *WaypointGetJobStreamOK) GetPayload() *WaypointGetJobStreamOKBody {
	return o.Payload
}

func (o *WaypointGetJobStreamOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(WaypointGetJobStreamOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointGetJobStreamDefault creates a WaypointGetJobStreamDefault with default headers values
func NewWaypointGetJobStreamDefault(code int) *WaypointGetJobStreamDefault {
	return &WaypointGetJobStreamDefault{
		_statusCode: code,
	}
}

/*
WaypointGetJobStreamDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointGetJobStreamDefault struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint get job stream default response
func (o *WaypointGetJobStreamDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint get job stream default response has a 2xx status code
func (o *WaypointGetJobStreamDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint get job stream default response has a 3xx status code
func (o *WaypointGetJobStreamDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint get job stream default response has a 4xx status code
func (o *WaypointGetJobStreamDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint get job stream default response has a 5xx status code
func (o *WaypointGetJobStreamDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint get job stream default response a status code equal to that given
func (o *WaypointGetJobStreamDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointGetJobStreamDefault) Error() string {
	return fmt.Sprintf("[GET /jobs/stream/{job_id}][%d] Waypoint_GetJobStream default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetJobStreamDefault) String() string {
	return fmt.Sprintf("[GET /jobs/stream/{job_id}][%d] Waypoint_GetJobStream default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointGetJobStreamDefault) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointGetJobStreamDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
WaypointGetJobStreamOKBody Stream result of hashicorp.waypoint.GetJobStreamResponse
swagger:model WaypointGetJobStreamOKBody
*/
type WaypointGetJobStreamOKBody struct {

	// error
	Error *models.GrpcGatewayRuntimeStreamError `json:"error,omitempty"`

	// result
	Result *models.HashicorpWaypointGetJobStreamResponse `json:"result,omitempty"`
}

// Validate validates this waypoint get job stream o k body
func (o *WaypointGetJobStreamOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateError(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateResult(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *WaypointGetJobStreamOKBody) validateError(formats strfmt.Registry) error {
	if swag.IsZero(o.Error) { // not required
		return nil
	}

	if o.Error != nil {
		if err := o.Error.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("waypointGetJobStreamOK" + "." + "error")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("waypointGetJobStreamOK" + "." + "error")
			}
			return err
		}
	}

	return nil
}

func (o *WaypointGetJobStreamOKBody) validateResult(formats strfmt.Registry) error {
	if swag.IsZero(o.Result) { // not required
		return nil
	}

	if o.Result != nil {
		if err := o.Result.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("waypointGetJobStreamOK" + "." + "result")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("waypointGetJobStreamOK" + "." + "result")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this waypoint get job stream o k body based on the context it is used
func (o *WaypointGetJobStreamOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateError(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateResult(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *WaypointGetJobStreamOKBody) contextValidateError(ctx context.Context, formats strfmt.Registry) error {

	if o.Error != nil {
		if err := o.Error.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("waypointGetJobStreamOK" + "." + "error")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("waypointGetJobStreamOK" + "." + "error")
			}
			return err
		}
	}

	return nil
}

func (o *WaypointGetJobStreamOKBody) contextValidateResult(ctx context.Context, formats strfmt.Registry) error {

	if o.Result != nil {
		if err := o.Result.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("waypointGetJobStreamOK" + "." + "result")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("waypointGetJobStreamOK" + "." + "result")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *WaypointGetJobStreamOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *WaypointGetJobStreamOKBody) UnmarshalBinary(b []byte) error {
	var res WaypointGetJobStreamOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
