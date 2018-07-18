package vagrant

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Vagrantfile represents a Vagrantfile and offers read & write operations
type Vagrantfile struct {
	Path          string
	ConfigVersion string
	machines      map[string]*Machine
}

// Machine represents a vagrant VM
type Machine struct {
	Box      string
	Provider string
	Cpus     string
	Memory   string
}

// NewVagrantfile creates a new Vagrantfile object
func NewVagrantfile(path, configVersion string) *Vagrantfile {
	return &Vagrantfile{
		Path:          path,
		ConfigVersion: configVersion,
		machines:      make(map[string]*Machine),
	}
}

// Read creates a new Vagrantfile object by reading its content
// from the given path
func Read(path string) (*Vagrantfile, error) {
	//TODO: implement
	return nil, errors.New("Not implemented")
}

func (v *Vagrantfile) Write() error {
	file, err := os.Create(v.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	err = v.writeHead(w)
	if err != nil {
		return err
	}
	err = v.writeMachines(w)
	if err != nil {
		return err
	}
	err = v.writeEnd(w)
	if err != nil {
		return err
	}

	return w.Flush()
}

// SetMachine adds or replaces the machine with the given name
func (v *Vagrantfile) SetMachine(name string, machine *Machine) {
	v.machines[name] = machine
}

// RemoveMachine removes the machine with the given name
func (v *Vagrantfile) RemoveMachine(name string) {
	delete(v.machines, name)
}

// must is a helper function to reduce the two return values
// of bufio.Writer.WriteString to a single return value
func must(n int, err error) error {
	return err
}

func (v *Vagrantfile) writeHead(w *bufio.Writer) error {
	return must(w.WriteString(fmt.Sprintf("Vagrant.configure(\"%s\") do |config|\n", v.ConfigVersion)))
}

func (v *Vagrantfile) writeEnd(w *bufio.Writer) error {
	return must(w.WriteString("end\n"))
}

func (v *Vagrantfile) writeMachines(w *bufio.Writer) error {
	for name := range v.machines {
		err := v.writeMachine(w, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Vagrantfile) writeMachine(w *bufio.Writer, name string) error {
	if m, ok := v.machines[name]; ok {
		err := must(w.WriteString(fmt.Sprintf("  config.vm.define \"%s\" do |m|\n    m.vm.box = \"%s\"\n", name, m.Box)))
		if err != nil {
			return err
		}
		if p := m.Provider; p != "" {
			err := must(w.WriteString(fmt.Sprintf("    m.vm.provider \"%s\" do |p|\n", p)))
			if err != nil {
				return err
			}
			switch p {
			case "virtualbox":
				err = must(w.WriteString(fmt.Sprintf("      p.cpus = %s\n      p.memory = %s\n", m.Cpus, m.Memory)))
			default:
				err = must(w.WriteString(fmt.Sprintf("      # Unsupported provider '%s'\n", p)))
			}
			if err != nil {
				return err
			}
			err = must(w.WriteString("    end\n"))
			if err != nil {
				return err
			}
		}
		return must(w.WriteString("  end\n"))
	}
	return nil
}
