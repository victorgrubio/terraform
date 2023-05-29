// go:build tools

// The go generate command uses //go:generate comment directives to execute commands. 
// Calling the tfplugindocs command under that framework ensures that Go can 
// automatically generate consistent documentation for your provider, 
// without external developer tooling such as make.

package tools

import (
	// Ensure documentation generator is not removed from go.mod.
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
