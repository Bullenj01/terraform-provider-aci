package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciFabricPathEndpoint() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciFabricPathEndpointRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func getRemoteFabricPathEndpoint(client *client.Client, dn string) (*models.FabricPathEndpoint, error) {
	fabricPathEpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricPathEp := models.FabricPathEndpointFromContainer(fabricPathEpCont)

	if fabricPathEp.DistinguishedName == "" {
		return nil, fmt.Errorf("FabricPathEndpoint %s not found", fabricPathEp.DistinguishedName)
	}

	return fabricPathEp, nil
}

func setFabricPathEndpointAttributes(fabricPathEp *models.FabricPathEndpoint, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fabricPathEp.DistinguishedName)
	d.Set("description", fabricPathEp.Description)

	fabricPathEpMap, _ := fabricPathEp.ToMap()
	d.Set("name", fabricPathEpMap["name"])

	return d
}

func dataSourceAciFabricPathEndpointRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)
	podID := d.Get("pod_id").(string)
	nodeID := d.Get("node_id").(string)

	pDN := fmt.Sprintf("topology/pod-%s/paths-%s", podID, nodeID)

	rn := fmt.Sprintf("pathep-[%s]", name)

	dn := fmt.Sprintf("%s/%s", pDN, rn)

	fabricPathEp, err := getRemoteFabricPathEndpoint(aciClient, dn)
	if err != nil {
		return err
	}

	setFabricPathEndpointAttributes(fabricPathEp, d)
	return nil
}
