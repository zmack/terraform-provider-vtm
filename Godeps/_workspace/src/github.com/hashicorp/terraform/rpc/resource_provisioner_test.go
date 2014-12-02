package rpc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/whitepages/terraform-provider-stingray/Godeps/_workspace/src/github.com/hashicorp/terraform/terraform"
)

func TestResourceProvisioner_impl(t *testing.T) {
	var _ terraform.ResourceProvisioner = new(ResourceProvisioner)
}

func TestResourceProvisioner_apply(t *testing.T) {
	client, server := testNewClientServer(t)
	defer client.Close()

	p := server.ProvisionerFunc().(*terraform.MockResourceProvisioner)

	provisioner, err := client.ResourceProvisioner()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Apply
	output := &terraform.MockUIOutput{}
	state := &terraform.InstanceState{}
	conf := &terraform.ResourceConfig{}
	err = provisioner.Apply(output, state, conf)
	if !p.ApplyCalled {
		t.Fatal("apply should be called")
	}
	if !reflect.DeepEqual(p.ApplyConfig, conf) {
		t.Fatalf("bad: %#v", p.ApplyConfig)
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
}

func TestResourceProvisioner_validate(t *testing.T) {
	p := new(terraform.MockResourceProvisioner)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provisioner := &ResourceProvisioner{Client: client, Name: name}

	// Configure
	config := &terraform.ResourceConfig{
		Raw: map[string]interface{}{"foo": "bar"},
	}
	w, e := provisioner.Validate(config)
	if !p.ValidateCalled {
		t.Fatal("configure should be called")
	}
	if !reflect.DeepEqual(p.ValidateConfig, config) {
		t.Fatalf("bad: %#v", p.ValidateConfig)
	}
	if w != nil {
		t.Fatalf("bad: %#v", w)
	}
	if e != nil {
		t.Fatalf("bad: %#v", e)
	}
}

func TestResourceProvisioner_validate_errors(t *testing.T) {
	p := new(terraform.MockResourceProvisioner)
	p.ValidateReturnErrors = []error{errors.New("foo")}

	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provisioner := &ResourceProvisioner{Client: client, Name: name}

	// Configure
	config := &terraform.ResourceConfig{
		Raw: map[string]interface{}{"foo": "bar"},
	}
	w, e := provisioner.Validate(config)
	if !p.ValidateCalled {
		t.Fatal("configure should be called")
	}
	if !reflect.DeepEqual(p.ValidateConfig, config) {
		t.Fatalf("bad: %#v", p.ValidateConfig)
	}
	if w != nil {
		t.Fatalf("bad: %#v", w)
	}

	if len(e) != 1 {
		t.Fatalf("bad: %#v", e)
	}
	if e[0].Error() != "foo" {
		t.Fatalf("bad: %#v", e)
	}
}

func TestResourceProvisioner_validate_warns(t *testing.T) {
	p := new(terraform.MockResourceProvisioner)
	p.ValidateReturnWarns = []string{"foo"}

	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provisioner := &ResourceProvisioner{Client: client, Name: name}

	// Configure
	config := &terraform.ResourceConfig{
		Raw: map[string]interface{}{"foo": "bar"},
	}
	w, e := provisioner.Validate(config)
	if !p.ValidateCalled {
		t.Fatal("configure should be called")
	}
	if !reflect.DeepEqual(p.ValidateConfig, config) {
		t.Fatalf("bad: %#v", p.ValidateConfig)
	}
	if e != nil {
		t.Fatalf("bad: %#v", e)
	}

	expected := []string{"foo"}
	if !reflect.DeepEqual(w, expected) {
		t.Fatalf("bad: %#v", w)
	}
}
