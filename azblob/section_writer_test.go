package azblob

import (
	"bytes"
	"io"

	chk "gopkg.in/check.v1"
)

func (s *aztestsSuite) TestSectionWriter(c *chk.C) {
	b := [10]byte{}
	buffer := newBytesWriter(b[:])

	section := newSectionWriter(buffer, 0, 5)
	c.Assert(section.count, chk.Equals, int64(5))
	c.Assert(section.offset, chk.Equals, int64(0))
	c.Assert(section.position, chk.Equals, int64(0))

	count, err := section.Write([]byte{1, 2, 3})
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, 3)
	c.Assert(section.position, chk.Equals, int64(3))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{4, 5, 6})
	c.Assert(err, chk.ErrorMatches, "Not enough space for all bytes")
	c.Assert(count, chk.Equals, 2)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{6, 7, 8})
	c.Assert(err, chk.ErrorMatches, "End of section reached")
	c.Assert(count, chk.Equals, 0)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	section = newSectionWriter(buffer, 5, 6)
	c.Assert(section.count, chk.Equals, int64(6))
	c.Assert(section.offset, chk.Equals, int64(5))
	c.Assert(section.position, chk.Equals, int64(0))

	count, err = section.Write([]byte{6, 7, 8})
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, 3)
	c.Assert(section.position, chk.Equals, int64(3))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0})

	count, err = section.Write([]byte{9, 10, 11})
	c.Assert(err, chk.ErrorMatches, "Not enough space for all bytes")
	c.Assert(count, chk.Equals, 2)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	count, err = section.Write([]byte{11, 12, 13})
	c.Assert(err, chk.ErrorMatches, "Offset value is out of range")
	c.Assert(count, chk.Equals, 0)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

func (s *aztestsSuite) TestSectionWriterCopySrcDestEmpty(c *chk.C) {
	input := make([]byte, 0)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := newBytesWriter(output)
	section := newSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, int64(0))
}

func (s *aztestsSuite) TestSectionWriterCopyDestEmpty(c *chk.C) {
	input := make([]byte, 10)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := newBytesWriter(output)
	section := newSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	c.Assert(err, chk.ErrorMatches, "End of section reached")
	c.Assert(count, chk.Equals, int64(0))
}
