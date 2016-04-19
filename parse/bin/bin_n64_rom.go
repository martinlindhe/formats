package bin

// Nintendo 64 ROM image

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func N64ROM(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isN64ROM(&hdr) {
		return nil, nil
	}
	return parseN64ROM(file, pl)
}

func isN64ROM(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0x80 || b[1] != 0x37 || b[2] != 0x12 || b[3] != 0x40 {
		return false
	}
	return true
}

func parseN64ROM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le}, // XXX le/be
		}}}

	return &pl, nil
}

/*
override public List<Chunk> GetFileStructure()
{
    if (!IsRecognized())
        throw new Exception("not a z64");

    List<Chunk> res = new List<Chunk>();

    var header = new Chunk();
    header.offset = 0;
    header.length = 0x40;
    header.Text = "Z64 header";
    res.Add(header);

    var dom1LatReg = new ByteChunk("initial PI_BSB_DOM1_LAT_REG value");
    dom1LatReg.offset = 0;
    header.Nodes.Add(dom1LatReg);

    var dom1PgsReg = dom1LatReg.RelativeToByte("initial PI_BSB_DOM1_PGS_REG value");
    header.Nodes.Add(dom1PgsReg);

    var dom1PwdReg = dom1PgsReg.RelativeToByte("initial PI_BSB_DOM1_PWD_REG value");
    header.Nodes.Add(dom1PwdReg);

    var dom1PgsReg2 = dom1PwdReg.RelativeToByte("initial PI_BSB_DOM1_PGS_REG value"); // XXX TYPO!?!?
    header.Nodes.Add(dom1PgsReg2);

    var ClockRate = dom1PgsReg2.RelativeToLittleEndian32("Clock rate");
    header.Nodes.Add(ClockRate);
    var ProgramCounter = ClockRate.RelativeToLittleEndian32("Program Counter (PC)");
    header.Nodes.Add(ProgramCounter);
    var Release = ProgramCounter.RelativeToLittleEndian32("Release");
    header.Nodes.Add(Release);
    var Crc1 = Release.RelativeToLittleEndian32("CRC 1");
    header.Nodes.Add(Crc1);
    var Crc2 = Crc1.RelativeToLittleEndian32("CRC 2");
    header.Nodes.Add(Crc2);
    var Unknown1 = Crc2.RelativeToLittleEndian32("Reserved 1");
    header.Nodes.Add(Unknown1);
    var Unknown2 = Unknown1.RelativeToLittleEndian32("Reserved 2");
    header.Nodes.Add(Unknown2);

    var ImageName = Unknown2.RelativeTo("Image name", 20); // Padded with 0x00 or spaces (0x20)
    header.Nodes.Add(ImageName);

    var Unknown3 = ImageName.RelativeToLittleEndian32("Reserved 3");
    header.Nodes.Add(Unknown3);

    var Manufacturer = Unknown3.RelativeToLittleEndian32("Manufacturer ID");

    int ManufacturerValue = ReadInt32BE(Manufacturer.offset);
    switch (ManufacturerValue) {
    case 0x0000004E:
        Manufacturer.Text += " = Nintendo";
        break;
    default:
        Console.WriteLine("Unrecognized manufacturer id = 0x" + ManufacturerValue.ToString("x8"));
        break;
    }

    header.Nodes.Add(Manufacturer);

    var Cartridge = Manufacturer.RelativeToLittleEndian16("Cartridge ID");
    header.Nodes.Add(Cartridge);

    var Country = Cartridge.RelativeToByte("Country");
    var CountryValue = ReadByte(Country.offset);


    0x41 'A' (not documented, generic NTSC?)
    0x42 'B' "Brazilian"
    0x43 'C' "Chinese"
    0x44 'D' "German"
    0x45 'E' "North America"
    0x46 'F' "French"
    0x47 'G': Gateway 64 (NTSC)
    0x48 'H' "Dutch"
    0x49 'I' "Italian"
    0x4A 'J' "Japanese"
    0x4B 'K' "Korean"
    0x4C 'L': Gateway 64 (PAL)
    0x4E 'N' "Canadian"
    0x50 'P' "European (basic spec.)"
    0x53 'S' "Spanish"
    0x55 'U' "Australian"
    0x57 'W' "Scandinavian"
    0x58 'X' "Others"
    0x59 'Y' "Others"
    0x5A 'Z' "Others"

    string CountryName = "";
    switch (CountryValue) {
    case 0x44:
        CountryName = "Germany";
        break;
    case 0x45:
        CountryName = "USA";
        break;
    case 0x4A:
        CountryName = "Japan";
        break;
    case 0x50:
        CountryName = "Europe";
        break;
    case 0x55:
        CountryName = "Australia";
        break;

    default:
        Console.WriteLine("NOTICE unrecognized country code = " + CountryValue.ToString("x2"));
        break;
    }
    Country.Text += " = " + CountryName;
    header.Nodes.Add(Country);

    var MaskRomVersion = Country.RelativeToByte("Mask ROM version");
    header.Nodes.Add(MaskRomVersion);

    //0040h - 0FFFh (1008 dwords): Boot code
    var BootCode = MaskRomVersion.RelativeTo("Boot code", 0x0FC0);
    header.Nodes.Add(BootCode);

    return res;
}
*/
