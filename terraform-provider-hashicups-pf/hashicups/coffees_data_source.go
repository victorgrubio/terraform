package hashicups

import (
	"context"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &coffeesDataSource{}
	_ datasource.DataSourceWithConfigure = &coffeesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation
func NewCoffeesDataSource() datasource.DataSource {
	return &coffeesDataSource{}
}

// coffeesDataSource is the data source implementation
// Data sources use the optional Configure method to fetch configured clients from the provider.
// The provider configures the HashiCups client and the data source can save a reference to that client for its operations.
// Allow your data source type to store a reference to the HashiCups clien
type coffeesDataSource struct {
	client *hashicups.Client
}

// coffeesDataSourceModel maps the data source schema data
type coffeesDataSourceModel struct {
    ID      types.String   `tfsdk:"id"`
	Coffees []coffeesModel `tfsdk:"coffees"`
}

// coffeesModel maps coffees schema data.
type coffeesModel struct {
	ID          types.Int64               `tfsdk:"id"`
	Name        types.String              `tfsdk:"name"`
	Teaser      types.String              `tfsdk:"teaser"`
	Description types.String              `tfsdk:"description"`
	Price       types.Float64             `tfsdk:"price"`
	Image       types.String              `tfsdk:"image"`
	Ingredients []coffeesIngredientsModel `tfsdk:"ingredients"`
}

// coffeesIngredientsModel
type coffeesIngredientsModel struct {
	ID types.Int64 `tfsdk:"id"`
}

// Metadata returns the data source type name
func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

// Schema defines the schema for the data source
func (d *coffeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"coffees": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"teaser": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"price": schema.Float64Attribute{
							Computed: true,
						},
						"image": schema.StringAttribute{
							Computed: true,
						},
						"ingredients": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// The data source uses the Read method to refresh the Terraform state based on the schema data.
// The hashicups_coffees data source will use the configured HashiCups client to call the HashiCups API coffee listing endpoint
// and save this data to the Terraform state.

// 1. Reads coffees list. The method invokes the API client's GetCoffees method.
// 2. Maps response body to schema attributes. After the method reads the coffees, it maps the []hashicups.Coffee response to coffeesModel so the data source can set the Terraform state.
// 3. Sets state with coffees list.
// Read refreshes the Terraform state with the latest data.
func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state coffeesDataSourceModel

	coffees, err := d.client.GetCoffees()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read HashiCups Coffees",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, coffee := range coffees {
		coffeeState := coffeesModel{
			ID:          types.Int64Value((int64(coffee.ID))),
			Name:        types.StringValue(coffee.Name),
			Teaser:      types.StringValue(coffee.Teaser),
			Description: types.StringValue(coffee.Description),
			Price:       types.Float64Value((float64(coffee.Price))),
			Image:       types.StringValue(coffee.Image),
		}

		for _, ingredient := range coffee.Ingredient {
			coffeeState.Ingredients = append(coffeeState.Ingredients, coffeesIngredientsModel{
				ID: types.Int64Value(int64(ingredient.ID)),
			})
		}
		state.Coffees = append(state.Coffees, coffeeState)
	}

	state.ID = types.StringValue("placeholder")

	// Set state
    diags := resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
	

}

// Fetch the HashiCups client from the provider by adding the Configure method to your data source type.
func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*hashicups.Client)
}
