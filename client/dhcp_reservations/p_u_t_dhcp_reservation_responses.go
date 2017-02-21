package dhcp_reservations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/rackn/rocket-skates/models"
)

// PUTDhcpReservationReader is a Reader for the PUTDhcpReservation structure.
type PUTDhcpReservationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PUTDhcpReservationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPUTDhcpReservationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPUTDhcpReservationOK creates a PUTDhcpReservationOK with default headers values
func NewPUTDhcpReservationOK() *PUTDhcpReservationOK {
	return &PUTDhcpReservationOK{}
}

/*PUTDhcpReservationOK handles this case with default header values.

PUTDhcpReservationOK p u t dhcp reservation o k
*/
type PUTDhcpReservationOK struct {
	Payload *models.DhcpReservationInput
}

func (o *PUTDhcpReservationOK) Error() string {
	return fmt.Sprintf("[PUT /reservations/{id}][%d] pUTDhcpReservationOK  %+v", 200, o.Payload)
}

func (o *PUTDhcpReservationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DhcpReservationInput)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}