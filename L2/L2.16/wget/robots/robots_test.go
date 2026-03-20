package robots

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	content := `
User-agent: *
Disallow: /admin
Disallow: /private/
Allow: /

User-agent: Googlebot
Disallow: /secret
`
	rules := parse(content)

	assert.False(t, rules.IsAllowed("https://example.com/admin"))
	assert.False(t, rules.IsAllowed("https://example.com/admin/users"))
	assert.False(t, rules.IsAllowed("https://example.com/private/data"))
	assert.True(t, rules.IsAllowed("https://example.com/about"))
	assert.True(t, rules.IsAllowed("https://example.com/"))
	// /secret запрещён только для Googlebot, не для *
	assert.True(t, rules.IsAllowed("https://example.com/secret"))
}

func TestParseEmpty(t *testing.T) {
	rules := parse("")
	assert.True(t, rules.IsAllowed("https://example.com/anything"))
}

func TestParseNoDisallow(t *testing.T) {
	content := `
User-agent: *
Allow: /
`
	rules := parse(content)
	assert.True(t, rules.IsAllowed("https://example.com/anything"))
}
