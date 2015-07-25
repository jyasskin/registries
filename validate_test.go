package validate

import (
	"io/ioutil"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ValidateSuite struct{}

var _ = Suite(&ValidateSuite{})

func (s *ValidateSuite) TestValidUUID(c *C) {
	c.Check(ValidUUID("00000000-0000-0000-0000-000000000000"), Equals, true)
	c.Check(ValidUUID("01234567-89ab-cdef-0123-456789abcdef"), Equals, true)
	c.Check(!ValidUUID("01234567-89AB-CDEF-0123-456789ABCDEF"), Equals, true)
	c.Check(!ValidUUID("g1234567-89ab-cdef-0123-456789abcdef"), Equals, true)
	c.Check(!ValidUUID("01234567-89ab-cdef-0123-456789abcdef0"), Equals, true)
	c.Check(!ValidUUID("0123456789abcdef0123456789abcdef"), Equals, true)
	c.Check(!ValidUUID("01234567089ab0cdef001230456789abcdef"), Equals, true)
}

func (s *ValidateSuite) TestBlacklistFile(c *C) {
	content, err := ioutil.ReadFile("gatt_blacklist.txt")
	c.Assert(err, IsNil)
	c.Assert(ValidateBlacklist(string(content)), IsNil)
}

func (s *ValidateSuite) TestValidateBlacklist(c *C) {
	c.Check(ValidateBlacklist(""), IsNil)

	// Lines are terminated by \n, not \r\n.
	c.Check(ValidateBlacklist("\r\n"), ErrorMatches, "line 1: '\r' is not a valid UUID")

	c.Check(ValidateBlacklist("# comment"), IsNil)
	c.Check(ValidateBlacklist("# comment\n" +
		"00001812-0000-1000-8000-00805f9b34fb"),
		IsNil)
	// No extraneous spaces.
	c.Check(ValidateBlacklist("  # comment\n" +
		"  00001812-0000-1000-8000-00805f9b34fb"),
		ErrorMatches, "line 1: Too many tokens")
	c.Check(ValidateBlacklist("00001812-0000-1000-8000-00805f9b34fb # not a comment"),
		ErrorMatches, "line 1: Too many tokens")
	c.Check(ValidateBlacklist("X0001812-0000-1000-8000-00805f9b34fb"),
		ErrorMatches, "line 1: 'X0001812-0000-1000-8000-00805f9b34fb' is not a valid UUID")
	c.Check(ValidateBlacklist("00001812-0000-1000-8000-00805f9b34fb exclude-reads"),
		IsNil)
	c.Check(ValidateBlacklist("00001812-0000-1000-8000-00805f9b34fb exclude-writes"),
		IsNil)
	c.Check(ValidateBlacklist("00001812-0000-1000-8000-00805f9b34fb\u00A0exclude-reads"),
		ErrorMatches, "line 1: '00001812-0000-1000-8000-00805f9b34fb\u00a0exclude-reads' is not a valid UUID")
	c.Check(ValidateBlacklist("X0001812-0000-1000-8000-00805f9b34fb exclude-reads"),
		ErrorMatches, "line 1: 'X0001812-0000-1000-8000-00805f9b34fb' is not a valid UUID")
	c.Check(ValidateBlacklist("00001812-0000-1000-8000-00805f9b34fb exclude"),
		ErrorMatches, "line 1: 'exclude' should be 'exclude', 'exclude-reads', or 'exclude-writes'")
	c.Check(ValidateBlacklist("00001812-0000-1000-8000-00805f9b34fb token token"),
		ErrorMatches, "line 1: Too many tokens")
}
