package tools

import "github.com/hashicorp/terraform-plugin-framework/types"

type KumaClient struct {
	Host     string
	Username string
	Password string
}

type KumaMonitorModel struct {
	ID      types.Number `tfsdk:"id"`
	Project types.String `tfsdk:"project"`
}

type KumaMonitorJsonModel struct {
	ID      int    `json:"id"`
	Project string `json:"project"`
}
