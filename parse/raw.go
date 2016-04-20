package parse

// for unrecognized files

func RAW(c *ParseChecker) (*ParsedLayout, error) {

	// TODO: make cmd/formats work without any Layout, to avoid a 0-length selected area
	c.ParsedLayout.FormatName = "raw"
	c.ParsedLayout.FileKind = Binary
	c.ParsedLayout.Layout = []Layout{{
		Offset: 0,
		Length: 0,
		Info:   "unrecognized data",
		Type:   Group,
		Childs: []Layout{
			{Offset: 0, Length: 0, Info: "data", Type: Bytes},
		}}}

	return &c.ParsedLayout, nil
}
