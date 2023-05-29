// Terraform providers are server processes that Terraform interacts with
// to handle each data source and resource operation, such as creating a resource on a remote system.
// Later in these tutorials, you will connect those Terraform operations to a locally running HashiCups API.

// Serving a provider follows these steps:

// Starts a provider server process. By implementing the main function,
// which is the code execution starting point for Go language programs,
// a long-running server will listen for Terraform requests.
// Framework provider servers also support optional functionality
// such as enabling support for debugging tools. You will not implement this functionality in these tutorials.

package main

import (
	"context"

	"terraform-provider-hashicups-pf/hashicups"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {

	providerserver.Serve(context.Background(), hashicups.New, providerserver.ServeOpts{
		Address: "hashicorp.com/edu/hashicups-pf",
	})
}
