package couchdb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringsFromSet(d interface{}) []string {
	s := d.(*schema.Set)
	ret := []string{}
	for _, v := range s.List() {
		ret = append(ret, v.(string))
	}
	return ret
}

func AppendDiagnostic(diags diag.Diagnostics, err error, summary string) diag.Diagnostics {
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   err.Error(),
	})
	return diags
}