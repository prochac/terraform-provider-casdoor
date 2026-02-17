// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sdkError checks the (bool, error) result from Casdoor SDK calls and adds
// appropriate diagnostics. Returns true if an error was recorded (caller
// should return early).
func sdkError(diags *diag.Diagnostics, ok bool, err error, msg string) bool {
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error %s", msg),
			err.Error(),
		)
		return true
	}

	if !ok {
		diags.AddError(
			fmt.Sprintf("Error %s", msg),
			fmt.Sprintf("Casdoor returned failure when %s", msg),
		)
		return true
	}

	return false
}

// stringListToSDK extracts a []string from a types.List. Returns nil if the
// list is null or unknown.
func stringListToSDK(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var result []string
	diags := list.ElementsAs(ctx, &result, false)

	return result, diags
}

// stringListFromSDK converts a []string to a types.List. Returns
// types.ListNull if the slice is nil or empty.
func stringListFromSDK(ctx context.Context, slice []string) (types.List, diag.Diagnostics) {
	if len(slice) == 0 {
		return types.ListNull(types.StringType), nil
	}

	return types.ListValueFrom(ctx, types.StringType, slice)
}
