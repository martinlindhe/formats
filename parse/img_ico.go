package parse

// STATUS xxx

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Windows Icon / Cursor image resources

func ICO(file *os.File) (*ParsedLayout, error) {

	if !isICO(file) {
		return nil, nil
	}
	return parseICO(file)
}

func readIconHeader(file *os.File) ([3]uint16, error) {

	file.Seek(0, os.SEEK_SET)
	var b [3]uint16
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return b, err
	}
	return b, nil
}

func isICO(file *os.File) bool {

	b, _ := readIconHeader(file)
	if b[0] != 0 {
		return false
	}

	// 1 = icon, 2 = cursor
	if b[1] != 1 && b[1] != 2 {
		return false
	}

	return true
}

func parseICO(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}
	typeName := ""

	hdr, _ := readIconHeader(file)
	switch hdr[1] {
	case 1:
		typeName = "icon"
	case 2:
		typeName = "cursor"
	default:
		typeName = "unknown"
	}

	fileHeader := Layout{
		Offset: 0,
		Length: 6,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: Uint16le},
			Layout{Offset: 2, Length: 2, Info: typeName, Type: Uint16le},
			Layout{Offset: 4, Length: 2, Info: "number of resources", Type: Uint16le},
		},
	}

	res.Layout = append(res.Layout, fileHeader)

	// map up resources
	numIcons := hdr[2]
	//	iconEntryLength := 16

	fmt.Println("parsing ", numIcons, " resources")

	for i := 0; i < int(numIcons); i++ {
		/*
		   iconEntry = new Chunk("Resource entry #" + (i + 1));
		   iconEntry.length = iconEntryLength;
		   iconEntry.offset = numIcons.offset + numIcons.length + (i * iconEntry.length);
		   header.Nodes.Add(iconEntry);

		   var width = new ByteChunk("Width");
		   width.offset = iconEntry.offset;
		   iconEntry.Nodes.Add(width);

		   var height = width.RelativeToByte("Height");
		   iconEntry.Nodes.Add(height);

		   //  ColorCount Maximum number of colors
		   var ColorCount = height.RelativeToByte("Color count");
		   iconEntry.Nodes.Add(ColorCount);

		   //  Reserved (always 0)
		   var Reserved = ColorCount.RelativeToByte("Reserved");
		   iconEntry.Nodes.Add(Reserved);

		   // Planes (always 0 or 1)
		   var Planes = Reserved.RelativeToLittleEndian16("Planes");
		   iconEntry.Nodes.Add(Planes);

		   // BitCount (always 0)
		   var BitCount = Planes.RelativeToLittleEndian16("Bit count");
		   iconEntry.Nodes.Add(BitCount);

		   //  BytesInRes Length of icon bitmap in bytes
		   var DataSize = BitCount.RelativeToLittleEndian32("Data size");
		   var DataSizeValue = (uint)ReadInt32(DataSize.offset);
		   iconEntry.Nodes.Add(DataSize);

		   // ImageOffset Offset position of icon bitmap in file
		   var ImageOffset = DataSize.RelativeToLittleEndian32("Image offset");
		   var OffsetValue = ReadInt32(ImageOffset.offset);
		   iconEntry.Nodes.Add(ImageOffset);

		   var Data = new Chunk();
		   Data.Text = "Resource data # " + (i + 1);
		   Data.offset = OffsetValue;
		   Data.length = DataSizeValue;
		   res.Add(Data);
		*/
	}

	// XXX
	return &res, nil
}

/*


    header.length = (uint)(6 + (numIconsValue * iconEntryLength));

    return res;
}
*/
