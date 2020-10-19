package dynatrace

import (
	"fmt"
	"github.com/Kissy/go-dynatrace/dynatrace"
	"github.com/go-openapi/runtime"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func isErrorCode(err error, code int) bool {
	if e, ok := err.(*runtime.APIError); ok && e.Code == code {
		return true
	}

	return false
}

func handleNotFoundError(err error, d *schema.ResourceData) error {
	if isErrorCode(err, 404) {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("error reading: %s: %s", d.Id(), err)
}

func handleBadRequestPayload(payload *dynatrace.ErrorEnvelope) error {
	constraintViolations := ""
	for _, constraintViolation := range payload.Error.ConstraintViolations {
		constraintViolations += fmt.Sprintf("%s: %s\n", constraintViolation.Path, constraintViolation.Message)
	}
	return fmt.Errorf("bad request: %d — %s\n %s", payload.Error.Code, payload.Error.Message, constraintViolations)
}