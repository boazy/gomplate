package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNSIsIdempotent(t *testing.T) {
	assert.True(t, NS() == NS())
}
func TestFuncs(t *testing.T) {
	m := &Ec2Meta{nonAWS: true}
	i := &Ec2Info{
		metaClient: m,
		describer:  func() InstanceDescriber { return DummyInstanceDescriber{} },
	}
	af := &Funcs{meta: m, info: i}
	assert.Equal(t, "unknown", af.EC2Region())
	assert.Equal(t, "", af.EC2Meta("foo"))
	assert.Equal(t, "", af.EC2Tag("foo"))
	assert.Equal(t, "unknown", af.EC2Region())
}
