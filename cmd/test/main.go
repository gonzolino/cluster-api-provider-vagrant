package main

import (
	"fmt"

	vagrant "github.com/gonzolino/cluster-api-provider-vagrant/cloud/vagrant/client"
)

func main() {
	v := vagrant.NewVagrantfile("Vagrantfile", "2")
	v.SetMachine("test", &vagrant.Machine{
		Box:      "ubuntu/trusty64",
		Provider: "virtualbox",
		Cpus:     "1",
		Memory:   "1024",
	})
	err := v.Write()
	if err != nil {
		fmt.Printf("Failed to write Vagrantfile to %s: %v", v.Path, err)
	}
}
