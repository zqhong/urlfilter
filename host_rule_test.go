package urlfilter

import (
	"bytes"
	"log"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHostRuleText(t *testing.T) {
	rule, err := NewHostRule("127.0.1.1       thishost.mydomain.org  thishost", 1)
	assert.Nil(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, 1, rule.FilterListID)
	assert.Equal(t, net.IPv4(127, 0, 1, 1), rule.IP)
	assert.Equal(t, 2, len(rule.Hostnames))
	assert.Equal(t, "thishost.mydomain.org", rule.Hostnames[0])
	assert.Equal(t, "thishost", rule.Hostnames[1])

	rule, err = NewHostRule("209.237.226.90  www.opensource.org", 1)
	assert.Nil(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, 1, rule.FilterListID)
	assert.Equal(t, net.IPv4(209, 237, 226, 90), rule.IP)
	assert.Equal(t, 1, len(rule.Hostnames))
	assert.Equal(t, "www.opensource.org", rule.Hostnames[0])

	rule, err = NewHostRule("::1             localhost ip6-localhost ip6-loopback", 1)
	assert.Nil(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, 1, rule.FilterListID)
	assert.Equal(t, net.ParseIP("::1"), rule.IP)
	assert.Equal(t, 3, len(rule.Hostnames))
	assert.Equal(t, "localhost", rule.Hostnames[0])
	assert.Equal(t, "ip6-localhost", rule.Hostnames[1])
	assert.Equal(t, "ip6-loopback", rule.Hostnames[2])

	rule, err = NewHostRule("example.org", 1)
	assert.Nil(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, 1, rule.FilterListID)
	assert.Equal(t, net.IPv4(0, 0, 0, 0), rule.IP)
	assert.Equal(t, 1, len(rule.Hostnames))
	assert.Equal(t, "example.org", rule.Hostnames[0])

	rule, err = NewHostRule("#::1             localhost ip6-localhost ip6-loopback", 1)
	assert.NotNil(t, err)
	assert.Nil(t, rule)

	rule, err = NewHostRule("||example.org", 1)
	assert.NotNil(t, err)
	assert.Nil(t, rule)

	rule, err = NewHostRule("", 1)
	assert.NotNil(t, err)
	assert.Nil(t, rule)
}

func TestHostRuleMatch(t *testing.T) {
	rule, err := NewHostRule("127.0.1.1       thishost.mydomain.org  thishost", 1)
	assert.Nil(t, err)
	assert.True(t, rule.Match("thishost.mydomain.org"))
	assert.True(t, rule.Match("thishost"))
	assert.False(t, rule.Match("mydomain.org"))
	assert.False(t, rule.Match("example.org"))

	rule, err = NewHostRule("209.237.226.90  www.opensource.org", 1)
	assert.Nil(t, err)
	assert.True(t, rule.Match("www.opensource.org"))
	assert.False(t, rule.Match("opensource.org"))
}

func TestHostRuleSerialize(t *testing.T) {
	ruleText := "127.0.1.1       thishost.mydomain.org  thishost"
	rule, err := NewHostRule(ruleText, -1)
	assert.Nil(t, err)
	assert.NotNil(t, rule)

	b := bytes.Buffer{}
	length, err := SerializeRule(rule, &b)
	assert.Nil(t, err)
	assert.Equal(t, length, b.Len())

	log.Printf("Rule text length: %d", len(ruleText))
	log.Printf("Serialized length: %d", length)

	r, err := DeserializeRule(&b)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	deserializedRule := r.(*HostRule)
	assert.True(t, reflect.DeepEqual(rule, deserializedRule))
}
