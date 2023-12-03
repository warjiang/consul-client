package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/warjiang/consul-client/pkg/consul"
)

type LookUpOptions struct {
	Name string
	Tag  string
}

func NewLookUpOptions() *LookUpOptions {
	return &LookUpOptions{}
}

func NewCmdLookUp() *cobra.Command {
	o := NewLookUpOptions()
	cmd := &cobra.Command{
		Use:   "lookup {name}",
		Short: "look up host and port of given name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.New("please input name, example: sd lookup name")
			}
			name := args[0]
			o.Name = name
			return o.Run()
		},
	}
	o.AddFlags(cmd)
	return cmd
}

func (o *LookUpOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Tag, "tag", "", "Tag to filter result of lookup name")
}

func (o *LookUpOptions) Run() error {
	endpoints, err := consul.Lookup(o.Name, consul.WithTag(o.Tag))
	if err != nil {
		return err
	}
	fmt.Printf("Service name: %s\n", o.Name)
	fmt.Printf("Data center: pri\n\n")
	fmt.Printf("IP\t\tPort\tTags\n")
	if len(endpoints) > 0 {
		for _, v := range endpoints {
			fmt.Printf("%s\t%d\t%s\n", v.Host, v.Port, v.Tags.ToString())
		}
	}
	return nil
}
