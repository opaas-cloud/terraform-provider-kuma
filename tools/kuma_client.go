package tools

import "github.com/hashicorp/terraform-plugin-framework/types"

type KumaClient struct {
	Host     string
	Username string
	Password string
}

type KumaMonitorModel struct {
	ID  types.Number `tfsdk:"id"`
	URL types.String `tfsdk:"url"`
}

type KumaMonitorJsonModel struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
