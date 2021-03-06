package semgrepapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRulesets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesetsRead,
		Schema: map[string]*schema.Schema{
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ruleset_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_paths": &schema.Schema{
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRulesetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*AppContext)
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/registry/ruleset_rule_paths", "https://semgrep.dev/api"), nil)
	if c.isAuthenticated {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	rules := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&rules)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rules", rules); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
