package cli

import (
	"os"

	"github.com/digitalrebar/provision/client/isos"
	"github.com/spf13/cobra"
)

type IsoOps struct{}

func (be IsoOps) GetIndexes() map[string]string {
	return map[string]string{}
}

func (be IsoOps) List(parms map[string]string) (interface{}, error) {
	d, e := session.Isos.ListIsos(isos.NewListIsosParams(), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be IsoOps) Upload(path string, f *os.File) (interface{}, error) {
	d, e := session.Isos.UploadIso(isos.NewUploadIsoParams().WithPath(path).WithBody(f), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be IsoOps) Delete(id string) (interface{}, error) {
	_, e := session.Isos.DeleteIso(isos.NewDeleteIsoParams().WithPath(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return "Good", nil
}

func init() {
	tree := addIsoCommands()
	App.AddCommand(tree)
}

func addIsoCommands() (res *cobra.Command) {
	singularName := "iso"
	name := "isos"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: "Commands to manage isos on the provisioner",
	}
	commands := commonOps(singularName, name, &IsoOps{})
	res.AddCommand(commands...)
	return res
}
