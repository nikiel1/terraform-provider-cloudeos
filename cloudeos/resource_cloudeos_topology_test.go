// Copyright (c) 2020 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package cloudeos

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestResourceTopology(t *testing.T) {
	r.Test(t, r.TestCase{
		Providers:    testProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testResourceTopologyDestroy,
		Steps: []r.TestStep{
			{
				Config:      testResourceDuplicateTopology,
				ExpectError: regexp.MustCompile("cloudeos_topology topo-test1 already exists"),
			},
			{
				Config: testResourceInitialTopologyConfig,
				Check:  testResourceInitialTopologyCheck,
			},
			{
				Config: testResourceUpdatedTopologyConfig,
				Check:  testResourceUpdatedTopologyCheck,
			},
		},
	})
}

var testResourceInitialTopologyConfig = fmt.Sprintf(`
provider "cloudeos" {
  cvaas_domain = "apiserver.cv-play.corp.arista.io"
  cvaas_server = "www.cv-play.corp.arista.io"
  // clouddeploy token
  service_account_web_token = %q
}

resource "cloudeos_topology" "topology2" {
   topology_name = "topo-test2"
   bgp_asn = "65000-65100"
   vtep_ip_cidr = "1.0.0.0/16"
   terminattr_ip_cidr = "2.0.0.0/16"
   dps_controlplane_cidr = "3.0.0.0/16"
}

`, os.Getenv("token"))

var testResourceDuplicateTopology = fmt.Sprintf(`
provider "cloudeos" {
  cvaas_domain = "apiserver.cv-play.corp.arista.io"
  cvaas_server = "www.cv-play.corp.arista.io"
  // clouddeploy token
  service_account_web_token = %q
}
resource "cloudeos_topology" "topology0" {
   topology_name = "topo-test1"
   bgp_asn = "65000-65100"
   vtep_ip_cidr = "4.0.0.0/16"
   terminattr_ip_cidr = "5.0.0.0/16"
   dps_controlplane_cidr = "6.0.0.0/16"
}

resource "cloudeos_topology" "topology1" {
   topology_name = "topo-test1"
   bgp_asn = "65000-65100"
   vtep_ip_cidr = "10.0.0.0/16"
   terminattr_ip_cidr = "11.0.0.0/16"
   dps_controlplane_cidr = "12.0.0.0/16"
   depends_on = [cloudeos_topology.topology0]
}
`, os.Getenv("token"))

var resourceTopoID = ""

func testResourceInitialTopologyCheck(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["cloudeos_topology.topology2"]
	if resourceState == nil {
		return fmt.Errorf("cloudeos_topology.topology resource not found in state")
	}

	topoState := resourceState.Primary
	if topoState == nil {
		return fmt.Errorf("cloudeos_topology.topology resource has no primary instance")
	}

	if topoState.ID == "" {
		return fmt.Errorf("cloudeos_topology.topology ID not assigned %s", topoState.ID)
	}
	resourceTopoID = topoState.ID // use this for update testing
	if got, want := topoState.Attributes["topology_name"], "topo-test2"; got != want {
		return fmt.Errorf("topology topology_name contains %s; want %s", got, want)
	}

	if got, want := topoState.Attributes["bgp_asn"], "65000-65100"; got != want {
		return fmt.Errorf("topology bgp_asn contains %s; want %s", got, want)
	}

	if got, want := topoState.Attributes["vtep_ip_cidr"], "1.0.0.0/16"; got != want {
		return fmt.Errorf("topology vtep_ip_cidr contains %s; want %s", got, want)
	}

	if got, want := topoState.Attributes["terminattr_ip_cidr"], "2.0.0.0/16"; got != want {
		return fmt.Errorf("topology terminattr_ip_cidr contains %s; want %s", got, want)
	}

	if got, want := topoState.Attributes["dps_controlplane_cidr"], "3.0.0.0/16"; got != want {
		return fmt.Errorf("topology dps_controlplane_cidr contains %s; want %s", got, want)
	}

	return nil
}

var testResourceUpdatedTopologyConfig = fmt.Sprintf(`
provider "cloudeos" {
  cvaas_domain = "apiserver.cv-play.corp.arista.io"
  cvaas_server = "www.cv-play.corp.arista.io"
  // clouddeploy token
  service_account_web_token = %q
}

resource "cloudeos_topology" "topology2" {
   topology_name = "topo-test2"
   bgp_asn = "65000-65500"
   vtep_ip_cidr = "1.0.0.0/16"
   terminattr_ip_cidr = "2.0.0.0/16"
   dps_controlplane_cidr = "3.0.0.0/16"
}
`, os.Getenv("token"))

func testResourceUpdatedTopologyCheck(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["cloudeos_topology.topology2"]
	topoState := resourceState.Primary
	if topoState.ID != resourceTopoID {
		return fmt.Errorf("cloudeos_topology.topology ID has changed during update %s to %s",
			resourceTopoID, topoState.ID)
	}

	if got, want := topoState.Attributes["bgp_asn"], "65000-65500"; got != want {
		return fmt.Errorf("topology topology_name contains %s; want %s", got, want)
	}
	return nil
}

func testResourceTopologyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudeos_topology" {
			continue
		}
		// TODO
		return nil
	}
	return nil
}