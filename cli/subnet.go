package cli

import (
	"fmt"
	"strings"
	"net"
	"strconv"
	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/subnets"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
	"github.com/go-openapi/strfmt"
)

type SubnetOps struct{}

func (be SubnetOps) GetType() interface{} {
	return &models.Subnet{}
}

func (be SubnetOps) GetId(obj interface{}) (string, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to subnet create")
	}
	return *subnet.Name, nil
}

func (be SubnetOps) GetIndexes() map[string]string {
	b := &backend.Subnet{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be SubnetOps) List(parms map[string]string) (interface{}, error) {
	params := subnets.NewListSubnetsParams()
	if listLimit != -1 {
		t1 := int64(listLimit)
		params = params.WithLimit(&t1)
	}
	if listOffset != -1 {
		t1 := int64(listOffset)
		params = params.WithOffset(&t1)
	}
	for k, v := range parms {
		switch k {
		case "Name":
			params = params.WithName(&v)
		case "Subnet":
			params = params.WithSubnet(&v)
		case "Strategy":
			params = params.WithStrategy(&v)
		case "NextServer":
			params = params.WithNextServer(&v)
		}
	}

	d, e := session.Subnets.ListSubnets(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Get(id string) (interface{}, error) {
	d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Create(obj interface{}) (interface{}, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet create")
	}
	d, e := session.Subnets.CreateSubnet(subnets.NewCreateSubnetParams().WithBody(subnet), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet patch")
	}
	d, e := session.Subnets.PatchSubnet(subnets.NewPatchSubnetParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Delete(id string) (interface{}, error) {
	d, e := session.Subnets.DeleteSubnet(subnets.NewDeleteSubnetParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addSubnetCommands()
	App.AddCommand(tree)
}

func addSubnetCommands() (res *cobra.Command) {
	singularName := "subnet"
	name := "subnets"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	op := &SubnetOps{}
	commands := commonOps(singularName, name, op)

	commands = append(commands, &cobra.Command{
		Use:   "range [subnetName] [startIP] [endIP]",
		Short: fmt.Sprintf("set the range of a subnet"),
		Long:  `Helper function to set the range of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage=false
			subName :=args[0]
			StartAddr := args[1]
			EndAddr := args[2]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload

			var IPfirst strfmt.IPv4
			e=IPfirst.Scan(EndAddr)
			if e != nil {
				return fmt.Errorf("Error: %v is not a valid IPv4", StartAddr)
			}

			var IPlast strfmt.IPv4
			e=IPlast.Scan(EndAddr)
			if e != nil {
				return fmt.Errorf("Error: %v is not a valid IPv4", EndAddr)
			}

			sub.ActiveStart = &IPfirst
			sub.ActiveEnd = &IPlast
			fmt.Printf("startIP: %v\nendIP: %v\n",*sub.ActiveStart,*sub.ActiveEnd)
			return nil
			
		},
	})


	commands = append(commands, &cobra.Command{
		Use:   "subnet [subnetName] [subnet CIDR]",
		Short: fmt.Sprintf("Set the CIDR network address"),
		Long:  `Helper function to set the CIDR of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			CIDR:= args[1]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload

			_, _, e2:=net.ParseCIDR(CIDR)
			if e2!=nil{
				return	fmt.Errorf("Error: %v is not a valid subnet CIDR", CIDR)

			}
			sub.Subnet = &CIDR
			fmt.Printf("%v\n",*sub.Subnet)
			return nil
			
			
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "strategy [subnetName] [MAC]",
		Short: fmt.Sprintf("Set Subnet strategy"),
		Long:  `Helper function to get the range of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			MACAddress:=args[1]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload

			sub.Strategy = &MACAddress
			fmt.Printf("%v\n",*sub.Strategy)
			return nil
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "pickers [subnetName] [list]",
		Short: fmt.Sprintf("assigns IP allocation methods to a subnet"),
		Long:  `Helper function that accepts a string of methods to allocate IP addresses separated by commas`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			pickerString:=args[1]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload
			
			sub.Pickers = strings.Split(pickerString,",")
			fmt.Printf(strings.Join(sub.Pickers,", "))
			return nil
			
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "nextserver [subnetName] [IP]",
		Short: fmt.Sprintf("Set next non-reserved IP"),
		Long:  `Helper function to set the first non-reserved IP of a subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			IPAddr:= args[1]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload

			var nextIP strfmt.IPv4
			e=nextIP.Scan(IPAddr)
			if e != nil {
				return fmt.Errorf("Error: %v is not a valid IPv4", IPAddr)
			}

			sub.NextServer=&nextIP
			fmt.Printf("%v\n",*sub.NextServer)
			return nil
			
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "leasetimes [subnetName] [active] [reserved]",
		Short: fmt.Sprintf("Set the leasetimes of a subnet"),
		Long:  `Helper function to get the range of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			activeTimeString:=args[1]
			reservedTimeString:=args[2]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload

			activeTime64Int,e:=strconv.ParseInt(activeTimeString,10,32)
			if e!=nil{
				return fmt.Errorf("Error: %v could not be read as a number",activeTimeString)

			}
			activeTime:=int32(activeTime64Int)

			reservedTime64Int,e:=strconv.ParseInt(reservedTimeString,10,32)
			if e!=nil{
				return fmt.Errorf("Error: %v could not be read as a number",reservedTimeString)

			}
			reservedTime:=int32(reservedTime64Int)

			sub.ActiveLeaseTime=&activeTime
			sub.ReservedLeaseTime=&reservedTime
			fmt.Printf("Active Lease Times=%v\nReserved Lease Times=%v\n",*sub.ActiveLeaseTime,*sub.ReservedLeaseTime)
			return nil



		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "set [subnetName] option [number] to [value]",
		Short: fmt.Sprintf("Set the given subnet's dhcpOption to a value"),
		Long:  `Helper function that sets the specified dhcpOption from a given subnet to a value. You cannot add an option.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			ChangedVal:=args[2] 
			newVal:=args[4]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub:=d.Payload

			changeVal, err:=strconv.Atoi(ChangedVal)
			if err!=nil{
				return fmt.Errorf("Error: %v could not be read as a number",ChangedVal)

			}

			

			if changeVal>=len(sub.Options){
				return fmt.Errorf("Error: option %v does not exist",changeVal)
			}

			sub.Options[changeVal]=&models.DhcpOption{
				Value: &newVal,
			}

			fmt.Printf("%v\n",*sub.Options[changeVal])
			return nil
		},
	})


	commands = append(commands, &cobra.Command{
		Use:   "get [subnetName] option [number]",
		Short: fmt.Sprintf("Get dhcpOption [number]"),
		Long:  `Helper function that gets the specified dhcpOption from a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage=false
			subName:=args[0]
			gettingVal:=args[2]

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}

			getVal, err:=strconv.Atoi(gettingVal)
			if err!=nil{
				return fmt.Errorf("Error: %v could not be read as a number",gettingVal)
			}

			sub:=d.Payload
			if len(sub.Options)<=getVal{
				return fmt.Errorf("Error: option %v does not exist",getVal)
			}
			fmt.Printf("Option %v: %v",getVal,*sub.Options[getVal])
			return nil

		},
	})
	res.AddCommand(commands...)
	return res
}
