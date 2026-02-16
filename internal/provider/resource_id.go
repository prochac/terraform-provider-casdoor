// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// importStateOwnerName parses a "owner/name" import ID and sets the owner,
// name, and id attributes in the resource state.
func importStateOwnerName(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	owner, name, ok := strings.Cut(req.ID, "/")
	if !ok || owner == "" || name == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'owner/name', got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), owner)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
