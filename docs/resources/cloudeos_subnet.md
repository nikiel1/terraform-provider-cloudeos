# cloudeos_subnet

The `cloudeos_subnet` resource sends the information about the subnet that has been deployed to CVaaS.
## Example Usage

```hcl
resource "cloudeos_topology" "topology" {
   topology_name = "topo-test"
   bgp_asn = "65000-65100"
   vtep_ip_cidr = "1.0.0.0/16"
   terminattr_ip_cidr = "4.0.0.0/16"
   dps_controlplane_cidr = "3.0.0.0/16"
}

resource "cloudeos_clos" "clos" {
   name = "clos-test"
   topology_name = cloudeos_topology.topology.topology_name
   cv_container_name = "CloudLeaf"
}

resource "cloudeos_wan" "wan" {
   name = "wan-test"
   topology_name = cloudeos_topology.topology.topology_name
   cv_container_name = "CloudEdge"
}

resource "cloudeos_vpc_config" "vpc" {
  cloud_provider = "aws"
  topology_name = cloudeos_topology.topology.topology_name
  clos_name = cloudeos_clos.clos.name
  wan_name = cloudeos_wan.wan.name
  role = "CloudEdge"
  cnps = "Dev"
  tags = {
       Name = "edgeVpc"
       Cnps = "Dev"
  }
  region = "us-west-1"
}

resource "aws_vpc" "vpc" {
    cidr_block = "100.0.0.0/16"
}

resource "aws_security_group" "sg" {
  name = "example-sg"
}

resource "cloudeos_vpc_status" "vpc" {
  cloud_provider = cloudeos_vpc_config.vpc.cloud_provider
  vpc_id = aws_vpc.vpc.id
  security_group_id = aws_security_group.sg.id
  cidr_block = aws_vpc.vpc.cidr_block
  igw = aws_security_group.sg.name
  role = cloudeos_vpc_config.vpc.role
  topology_name = cloudeos_topology.topology.topology_name
  tags = cloudeos_vpc_config.vpc.tags
  clos_name = cloudeos_clos.clos.name
  wan_name = cloudeos_wan.wan.name
  cnps = cloudeos_vpc_config.vpc.cnps
  region = cloudeos_vpc_config.vpc.region
  account = "dummy_aws_account"
  tf_id = cloudeos_vpc_config.vpc.tf_id
}

resource "cloudeos_subnet" "subnet" {
  cloud_provider = cloudeos_vpc_status.vpc.cloud_provider
  vpc_id = cloudeos_vpc_status.vpc.vpc_id
  availability_zone = "us-west-1b"
  subnet_id = "subnet-id"
  cidr_block = "15.0.0.0/24"
  subnet_name = "edgeSubnet"
}
```

## Argument Reference

* `cloud_provider` - (Required) Cloud Provider in which the subnet is being deployed. Supported: aws or azure.
* `vpc_id` - (Required) VPC ID in which this subnet is created, equivalent to rg_name in Azure.
* `subnet_id` - (Required) ID of subnet deployed in AWS/Azure.
* `cidr_block` - (Required) CIDR of the subnet.
* `subnet_name` - (Required) Name of the subnet.
* `vnet_name` - (Optional) VNET name, only needed in Azure.
* `availability_zone` - (Optional) Availability zone.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `ID` - The ID of the cloudeos_subnet resource.