package urlfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNetworkRuleText(t *testing.T) {
	pattern, options, whitelist, err := parseRuleText("||example.org^")
	assert.Equal(t, "||example.org^", pattern)
	assert.Equal(t, "", options)
	assert.Equal(t, false, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("||example.org^$third-party")
	assert.Equal(t, "||example.org^", pattern)
	assert.Equal(t, "third-party", options)
	assert.Equal(t, false, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("@@||example.org^$third-party")
	assert.Equal(t, "||example.org^", pattern)
	assert.Equal(t, "third-party", options)
	assert.Equal(t, true, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("@@||example.org/this$is$path$third-party")
	assert.Equal(t, "||example.org/this$is$path", pattern)
	assert.Equal(t, "third-party", options)
	assert.Equal(t, true, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("||example.org/this$is$path$third-party")
	assert.Equal(t, "||example.org/this$is$path", pattern)
	assert.Equal(t, "third-party", options)
	assert.Equal(t, false, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("/regex/")
	assert.Equal(t, "/regex/", pattern)
	assert.Equal(t, "", options)
	assert.Equal(t, false, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("@@/regex/")
	assert.Equal(t, "/regex/", pattern)
	assert.Equal(t, "", options)
	assert.Equal(t, true, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("@@/regex/$replace=/test/test2/")
	assert.Equal(t, "/regex/", pattern)
	assert.Equal(t, "replace=/test/test2/", options)
	assert.Equal(t, true, whitelist)
	assert.Nil(t, err)

	pattern, options, whitelist, err = parseRuleText("/regex/$replace=/test/test2/")
	assert.Equal(t, "/regex/", pattern)
	assert.Equal(t, "replace=/test/test2/", options)
	assert.Equal(t, false, whitelist)
	assert.Nil(t, err)

	_, _, _, err = parseRuleText("@@")
	assert.NotNil(t, err)
}

func checkModifier(t *testing.T, name string, option NetworkRuleOption, enabled bool) {
	ruleText := "||example.org$" + name
	if (option & OptionWhitelistOnly) == option {
		ruleText = "@@" + ruleText
	}

	f, err := NewNetworkRule(ruleText, 0)
	assert.Nil(t, err)
	assert.NotNil(t, f)

	if enabled {
		assert.True(t, f.IsOptionEnabled(option))
	} else {
		assert.True(t, f.IsOptionDisabled(option))
	}
}

func TestParseModifiers(t *testing.T) {
	checkModifier(t, "important", OptionImportant, true)
	checkModifier(t, "third-party", OptionThirdParty, true)
	checkModifier(t, "~first-party", OptionThirdParty, true)
	checkModifier(t, "first-party", OptionThirdParty, false)
	checkModifier(t, "~third-party", OptionThirdParty, false)
	checkModifier(t, "match-case", OptionMatchCase, true)
	checkModifier(t, "~match-case", OptionMatchCase, false)

	checkModifier(t, "elemhide", OptionElemhide, true)
	checkModifier(t, "generichide", OptionGenerichide, true)
	checkModifier(t, "genericblock", OptionGenericblock, true)
	checkModifier(t, "jsinject", OptionJsinject, true)
	checkModifier(t, "urlblock", OptionUrlblock, true)
	checkModifier(t, "content", OptionContent, true)
	checkModifier(t, "extension", OptionExtension, true)

	checkModifier(t, "document", OptionElemhide, true)
	checkModifier(t, "document", OptionJsinject, true)
	checkModifier(t, "document", OptionUrlblock, true)
	checkModifier(t, "document", OptionContent, true)
	checkModifier(t, "document", OptionExtension, true)

	checkModifier(t, "stealth", OptionStealth, true)

	checkModifier(t, "popup", OptionPopup, true)
	checkModifier(t, "empty", OptionEmpty, true)
	checkModifier(t, "mp4", OptionMp4, true)
}

func TestDisablingExtensionModifier(t *testing.T) {
	ruleText := "@@||example.org$document,~extension"

	f, err := NewNetworkRule(ruleText, 0)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.False(t, f.IsOptionEnabled(OptionExtension))
	assert.False(t, f.IsOptionDisabled(OptionExtension))
}

func checkRequestType(t *testing.T, name string, requestType RequestType, permitted bool) {
	f, err := NewNetworkRule("||example.org^$"+name, 0)
	assert.Nil(t, err)
	assert.NotNil(t, f)

	if permitted {
		assert.Equal(t, f.permittedRequestTypes, requestType)
	} else {
		assert.Equal(t, f.restrictedRequestTypes, requestType)
	}
}

func TestParseRequestTypeModifiers(t *testing.T) {
	checkRequestType(t, "script", TypeScript, true)
	checkRequestType(t, "~script", TypeScript, false)

	checkRequestType(t, "stylesheet", TypeStylesheet, true)
	checkRequestType(t, "~stylesheet", TypeStylesheet, false)

	checkRequestType(t, "subdocument", TypeSubdocument, true)
	checkRequestType(t, "~subdocument", TypeSubdocument, false)

	checkRequestType(t, "object", TypeObject, true)
	checkRequestType(t, "~object", TypeObject, false)

	checkRequestType(t, "object", TypeObject, true)
	checkRequestType(t, "~object", TypeObject, false)

	checkRequestType(t, "image", TypeImage, true)
	checkRequestType(t, "~image", TypeImage, false)

	checkRequestType(t, "xmlhttprequest", TypeXmlhttprequest, true)
	checkRequestType(t, "~xmlhttprequest", TypeXmlhttprequest, false)

	checkRequestType(t, "object-subrequest", TypeObjectSubrequest, true)
	checkRequestType(t, "~object-subrequest", TypeObjectSubrequest, false)

	checkRequestType(t, "media", TypeMedia, true)
	checkRequestType(t, "~media", TypeMedia, false)

	checkRequestType(t, "font", TypeFont, true)
	checkRequestType(t, "~font", TypeFont, false)

	checkRequestType(t, "websocket", TypeWebsocket, true)
	checkRequestType(t, "~websocket", TypeWebsocket, false)

	checkRequestType(t, "other", TypeOther, true)
	checkRequestType(t, "~other", TypeOther, false)
}

func TestFindShortcut(t *testing.T) {
	shortcut := findShortcut("||example.org^")
	assert.Equal(t, "example.org", shortcut)

	shortcut = findShortcut("|https://*examp")
	assert.Equal(t, "https://", shortcut)

	shortcut = findRegexpShortcut("/example/")
	assert.Equal(t, "example", shortcut)

	shortcut = findRegexpShortcut("/^http:\\/\\/example/")
	assert.Equal(t, "/example", shortcut)

	shortcut = findRegexpShortcut("/^http:\\/\\/[a-z]+\\.example/")
	assert.Equal(t, "example", shortcut)

	shortcut = findRegexpShortcut("//")
	assert.Equal(t, "", shortcut)

	shortcut = findRegexpShortcut("/^http:\\/\\/(?!test.)example.org/")
	assert.Equal(t, "", shortcut)
}

func TestSimpleBasicRules(t *testing.T) {
	// Simple matching rule
	f, err := NewNetworkRule("||example.org^", 0)
	r := NewRequest("https://example.org/", "", TypeOther)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	// Simple regex rule
	f, err = NewNetworkRule("/example\\.org/", 0)
	r = NewRequest("https://example.org/", "", TypeOther)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))
}

func TestInvalidModifiers(t *testing.T) {
	_, err := NewNetworkRule("||example.org^$unknown", 0)
	assert.NotNil(t, err)

	// Whitelist-only modifier
	_, err = NewNetworkRule("||example.org^$elemhide", 0)
	assert.NotNil(t, err)

	// Blacklist-only modifier
	_, err = NewNetworkRule("@@||example.org^$popup", 0)
	assert.NotNil(t, err)
}

func TestMatchCase(t *testing.T) {
	f, err := NewNetworkRule("||example.org^$match-case", 0)
	r := NewRequest("https://example.org/", "", TypeOther)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://EXAMPLE.org/", "", TypeOther)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))
}

func TestThirdParty(t *testing.T) {
	f, err := NewNetworkRule("||example.org^$third-party", 0)

	// First-party 1
	r := NewRequest("https://example.org/", "", TypeOther)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	// First-party 2
	r = NewRequest("https://sub.example.org/", "https://example.org/", TypeOther)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	// Third-party
	r = NewRequest("https://example.org/", "https://example.com", TypeOther)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	f, err = NewNetworkRule("||example.org^$first-party", 0)

	// First-party 1
	r = NewRequest("https://example.org/", "", TypeOther)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	// First-party
	r = NewRequest("https://sub.example.org/", "https://example.org/", TypeOther)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	// Third-party
	r = NewRequest("https://example.org/", "https://example.com", TypeOther)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))
}

func TestContentType(t *testing.T) {
	// $script
	f, err := NewNetworkRule("||example.org^$script", 0)
	r := NewRequest("https://example.org/", "", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://example.org/", "", TypeDocument)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	// $script and $stylesheet
	f, err = NewNetworkRule("||example.org^$script,stylesheet", 0)
	r = NewRequest("https://example.org/", "", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://example.org/", "", TypeStylesheet)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://example.org/", "", TypeDocument)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	// Everything except $script and $stylesheet
	f, err = NewNetworkRule("@@||example.org^$~script,~stylesheet", 0)
	r = NewRequest("https://example.org/", "", TypeScript)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	r = NewRequest("https://example.org/", "", TypeStylesheet)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	r = NewRequest("https://example.org/", "", TypeDocument)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))
}

func TestDomainRestrictions(t *testing.T) {
	// Just one permitted domain
	f, err := NewNetworkRule("||example.org^$domain=example.org", 0)
	r := NewRequest("https://example.org/", "", TypeScript)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	r = NewRequest("https://example.org/", "https://example.org/", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://example.org/", "https://subdomain.example.org/", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	// One permitted, subdomain restricted
	f, err = NewNetworkRule("||example.org^$domain=example.org|~subdomain.example.org", 0)
	r = NewRequest("https://example.org/", "", TypeScript)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	r = NewRequest("https://example.org/", "https://example.org/", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://example.org/", "https://subdomain.example.org/", TypeScript)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	// One restricted
	f, err = NewNetworkRule("||example.org^$domain=~example.org", 0)
	r = NewRequest("https://example.org/", "", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))

	r = NewRequest("https://example.org/", "https://example.org/", TypeScript)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	r = NewRequest("https://example.org/", "https://subdomain.example.org/", TypeScript)
	assert.Nil(t, err)
	assert.False(t, f.Match(r))

	// Wide restricted
	f, err = NewNetworkRule("$domain=example.org", 0)
	r = NewRequest("https://example.com/", "https://example.org/", TypeScript)
	assert.Nil(t, err)
	assert.True(t, f.Match(r))
}

func TestInvalidDomainRestrictions(t *testing.T) {
	_, err := NewNetworkRule("||example.org^$domain=", 0)
	assert.NotNil(t, err)

	_, err = NewNetworkRule("||example.org^$domain=|example.com", 0)
	assert.NotNil(t, err)
}

func TestNetworkRulePriority(t *testing.T) {
	compareRulesPriority(t, "@@||example.org$important", "@@||example.org$important", false)
	compareRulesPriority(t, "@@||example.org$important", "||example.org$important", true)
	compareRulesPriority(t, "@@||example.org$important", "@@||example.org", true)
	compareRulesPriority(t, "@@||example.org$important", "||example.org", true)

	compareRulesPriority(t, "||example.org$important", "@@||example.org$important", false)
	compareRulesPriority(t, "||example.org$important", "||example.org$important", false)
	compareRulesPriority(t, "||example.org$important", "@@||example.org", true)
	compareRulesPriority(t, "||example.org$important", "||example.org", true)

	compareRulesPriority(t, "@@||example.org", "@@||example.org$important", false)
	compareRulesPriority(t, "@@||example.org", "||example.org$important", false)
	compareRulesPriority(t, "@@||example.org", "@@||example.org", false)
	compareRulesPriority(t, "@@||example.org", "||example.org", true)

	compareRulesPriority(t, "||example.org", "@@||example.org$important", false)
	compareRulesPriority(t, "||example.org", "||example.org$important", false)
	compareRulesPriority(t, "||example.org", "@@||example.org", false)
	compareRulesPriority(t, "||example.org", "||example.org", false)
}

func TestInvalidRule(t *testing.T) {
	r, err := NewNetworkRule("*$third-party", -1)
	assert.Nil(t, r)
	assert.NotNil(t, err)

	r, err = NewNetworkRule("$third-party", -1)
	assert.Nil(t, r)
	assert.NotNil(t, err)

	r, err = NewNetworkRule("ad$third-party", -1)
	assert.Nil(t, r)
	assert.NotNil(t, err)

	// This one is valid because it has domain restriction
	r, err = NewNetworkRule("$domain=ya.ru", -1)
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func compareRulesPriority(t *testing.T, left string, right string, expected bool) {
	l, err := NewNetworkRule(left, -1)
	assert.Nil(t, err)
	r, err := NewNetworkRule(right, -1)
	assert.Nil(t, err)
	assert.Equal(t, expected, l.isHigherPriority(r))
}
