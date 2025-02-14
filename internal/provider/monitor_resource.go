package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"math/big"
	"net/http"
	"terraform-provider-kuma/tools"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &monitorResource{}
	_ resource.ResourceWithConfigure = &monitorResource{}
)

// NewKumaResource is a helper function to simplify the provider implementation.
func NewMonitorResource() resource.Resource {
	return &monitorResource{}
}

// folderResource is the resource implementation.
type monitorResource struct {
	client *tools.KumaClient
}

// Configure adds the provider configured client to the resource.
func (r *monitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tools.KumaClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *tools.KumaClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *monitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Schema defines the schema for the resource.
func (r *monitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create a new resource.
func (r *monitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tools.KumaMonitorModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var monitor = tools.KumaMonitorJsonModel{
		Project: plan.Project.ValueString(),
	}

	out, err := json.Marshal(monitor)

	println(bytes.NewBuffer(out))

	if err != nil {
		resp.Diagnostics.AddError("Cannot send convert model to json", err.Error())
		return
	}

	request, err := http.NewRequest("POST", r.client.Host+"/add_monitor", bytes.NewBuffer(out))

	if err != nil {
		resp.Diagnostics.AddError("Cannot create request", err.Error())
		return
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send request", err.Error())
		return
	}

	if response.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(response.Body)
		resp.Diagnostics.AddError("Cannot send request", string(bodyBytes))
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	var data tools.KumaMonitorJsonModel

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	plan.ID = types.NumberValue(new(big.Float).SetInt64(int64(data.ID)))

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			resp.Diagnostics.AddError("Cannot close request", err.Error())
			return
		}
	}(response.Body)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *monitorResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *monitorResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tools.KumaMonitorModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var monitor = tools.KumaMonitorJsonModel{}

	out, err := json.Marshal(monitor)

	println(bytes.NewBuffer(out))

	if err != nil {
		resp.Diagnostics.AddError("Cannot send convert model to json", err.Error())
		return
	}

	request, err := http.NewRequest("DELETE", r.client.Host+"/add_monitor", bytes.NewBuffer(out))

	if err != nil {
		resp.Diagnostics.AddError("Cannot create request", err.Error())
		return
	}

	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot create delete request", err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			resp.Diagnostics.AddError("Cannot send delete request", err.Error())
			return
		}
	}(response.Body)
}
