package parse

/*

private Chunk ParsePEHeader()
{
    Log("PE header found");

    this.name = "PE executable (Win32)";

    var peHead = new Chunk();
    peHead.offset = ExtendedHeaderOffset;
    peHead.Text = "PE header";
    //                    peHead.childs = new List<Chunk>();

    var peIdentifier = new LittleEndian32BitChunk("PE identifier");
    peIdentifier.offset = ExtendedHeaderOffset;
    peHead.Nodes.Add(peIdentifier);

    // COFF header
    var Machine = peIdentifier.RelativeToLittleEndian16("Machine = ");

    BaseStream.Position = Machine.offset;
    int MachineValue = ReadInt16();
    switch (MachineValue) {
    case 0x14c:
        Machine.Text += "Intel 386";
        break;

    default:
        Machine.Text += "Unknown";
        throw new Exception("Unrecognized machine: " + MachineValue.ToString("x4"));
    }
    peHead.Nodes.Add(Machine);

    // NumberOfSections
    var NumberOfSections = Machine.RelativeToLittleEndian16("Sections");
    peHead.Nodes.Add(NumberOfSections);

    BaseStream.Position = NumberOfSections.offset;
    var NumberOfSectionsValue = (uint)ReadInt16();

    var TimeDateStamp = NumberOfSections.RelativeToLittleEndianDateStamp("TimeDateStamp");
    peHead.Nodes.Add(TimeDateStamp);

    var PointerToSymbolTable = TimeDateStamp.RelativeToLittleEndian32("Pointer To Symbol Table");
    peHead.Nodes.Add(PointerToSymbolTable);

    // NumberOfSymbols
    var NumberOfSymbols = PointerToSymbolTable.RelativeToLittleEndian32("Symbols");
    peHead.Nodes.Add(NumberOfSymbols);

    var SizeOfOptionalHeader = NumberOfSymbols.RelativeToLittleEndian16("Size Of Optional Header");
    peHead.Nodes.Add(SizeOfOptionalHeader);

    BaseStream.Position = SizeOfOptionalHeader.offset;
    var SizeOfOptionalHeaderValue = (uint)ReadInt16();

    var Characteristics = SizeOfOptionalHeader.RelativeToLittleEndian16("Characteristics");
    peHead.Nodes.Add(Characteristics);

    var optHead = new Chunk("Optional header");
    optHead.length = SizeOfOptionalHeaderValue;
    optHead.offset = Characteristics.offset + Characteristics.length;

    var optSignature = new Chunk("Optional signature");
    optSignature.offset = optHead.offset;
    optSignature.length = 2;
    optHead.Nodes.Add(optSignature);

    BaseStream.Position = optSignature.offset;

    if (ReadByte() != 0x0b || ReadByte() != 0x01)
        throw new Exception("Unrecognized optional header value found");

    var LinkerVersion = optSignature.RelativeToVersionMajorMinor16("Linker version");
    LinkerVersion.Text = "Linker version";
    optHead.Nodes.Add(LinkerVersion);

    var SizeOfCode = LinkerVersion.RelativeToLittleEndian32("Size of Code");
    optHead.Nodes.Add(SizeOfCode);

    var SizeOfInitializedData = SizeOfCode.RelativeToLittleEndian32("Size Of Initialized Data");
    optHead.Nodes.Add(SizeOfInitializedData);

    var SizeOfUninitializedData = SizeOfInitializedData.RelativeToLittleEndian32("Size Of Uninitialized Data");
    optHead.Nodes.Add(SizeOfUninitializedData);

    //The RVA of the code entry point
    var AddressOfEntryPoint = SizeOfUninitializedData.RelativeToLittleEndian32("Address Of Entry Point");
    BaseStream.Position = AddressOfEntryPoint.offset;
    int AddressOfEntryPointValue = ReadInt32();
    Log("XXXXX entry offset = " + AddressOfEntryPointValue.ToString("x6"));
    this.EntryPoint = AddressOfEntryPointValue;

    optHead.Nodes.Add(AddressOfEntryPoint);

    var BaseOfCode = AddressOfEntryPoint.RelativeToLittleEndian32("Base of Code");
    optHead.Nodes.Add(BaseOfCode);

    var BaseOfData = BaseOfCode.RelativeToLittleEndian32("Base of Data");
    optHead.Nodes.Add(BaseOfData);

    var ImageBase = BaseOfData.RelativeToLittleEndian32("Image Base");
    optHead.Nodes.Add(ImageBase);

    var SectionAlignment = ImageBase.RelativeToLittleEndian32("Section Alignment");
    optHead.Nodes.Add(SectionAlignment);

    var FileAlignment = SectionAlignment.RelativeToLittleEndian32("File Alignment");
    optHead.Nodes.Add(FileAlignment);

    var OSVersion = FileAlignment.RelativeToVersionMajorMinor32("OS Version");
    optHead.Nodes.Add(OSVersion);

    var ImageVersion = OSVersion.RelativeToVersionMajorMinor32("Image Version");
    optHead.Nodes.Add(ImageVersion);

    var SubsystemVersion = ImageVersion.RelativeToVersionMajorMinor32("Subsystem Version");
    optHead.Nodes.Add(SubsystemVersion);

    var Reserved = SubsystemVersion.RelativeToLittleEndian32("Reserved");
    optHead.Nodes.Add(Reserved);

    var SizeOfImage = Reserved.RelativeToLittleEndian32("Size of Image");
    optHead.Nodes.Add(SizeOfImage);

    var SizeOfHeaders = SizeOfImage.RelativeToLittleEndian32("Size of Headers");
    optHead.Nodes.Add(SizeOfHeaders);


    var Checksum = SizeOfHeaders.RelativeToLittleEndian32("Checksum");
    optHead.Nodes.Add(Checksum);

    var Subsystem = Checksum.RelativeToLittleEndian16("Subsystem = ");
    optHead.Nodes.Add(Subsystem);

    BaseStream.Position = Subsystem.offset;
    int SubsystemValue = ReadInt16();

    switch (SubsystemValue) {
    case 0x0001:
        Subsystem.Text += "Native";
        break;
    case 0x0002:
        Subsystem.Text += "GUI";
        break;
    case 0x0003:
        Subsystem.Text += "Console";
        break;
    case 0x0005:
        Subsystem.Text += "OS/2";
        break;
    case 0x0007:
        Subsystem.Text += "POSIX";
        break;
    default:
        Log("Unknown " + SubsystemValue.ToString("x4"));
        break;
    }

    var DLLCharacteristics = Subsystem.RelativeToLittleEndian16("DLL Characteristics");
    optHead.Nodes.Add(DLLCharacteristics);

    var SizeOfStackReserve = DLLCharacteristics.RelativeToLittleEndian32("Size Of Stack Reserve");
    optHead.Nodes.Add(SizeOfStackReserve);

    var SizeOfStackCommit = SizeOfStackReserve.RelativeToLittleEndian32("Size Of Stack Commit");
    optHead.Nodes.Add(SizeOfStackCommit);

    var SizeOfHeapReserve = SizeOfStackCommit.RelativeToLittleEndian32("Size Of Heap Reserve");
    optHead.Nodes.Add(SizeOfHeapReserve);

    var SizeOfHeapCommit = SizeOfHeapReserve.RelativeToLittleEndian32("Size Of Heap Commit");
    optHead.Nodes.Add(SizeOfHeapCommit);

    var LoaderFlags = SizeOfHeapCommit.RelativeToLittleEndian32("Loader Flags");
    optHead.Nodes.Add(LoaderFlags);

    var NumberOfRvaAndSizes = LoaderFlags.RelativeToLittleEndian32("Number of RVA");
    optHead.Nodes.Add(NumberOfRvaAndSizes);

    BaseStream.Position = NumberOfRvaAndSizes.offset;
    var NumberOfRvaAndSizesValue = (uint)ReadInt32();
    Log("NumberOfRvaAndSizesValue = " + NumberOfRvaAndSizesValue);

    if (NumberOfRvaAndSizesValue != 16)
        throw new Exception("odd number of RVA:s = " + NumberOfRvaAndSizesValue);

    var DataDirectory = new Chunk("Data Directory");
    DataDirectory.offset = NumberOfRvaAndSizes.offset + NumberOfRvaAndSizes.length;
    DataDirectory.length = NumberOfRvaAndSizesValue * 8;

    for (int i = 0; i < NumberOfRvaAndSizesValue; i++) {
        var RVAChunk = new Chunk();
        RVAChunk.length = 8;
        RVAChunk.offset = DataDirectory.offset + (i * RVAChunk.length);
        switch (i) {
        case 0:
            RVAChunk.Text = "Exports";
            break;
        case 1:
            RVAChunk.Text = "Imports";
            break;
        case 2:
            RVAChunk.Text = "Resources";
            break;
        case 5:
            RVAChunk.Text = "Base reolcations";
            break;
        case 9:
            RVAChunk.Text = "Thread Local Storage";
            break;
        case 12:
            RVAChunk.Text = "Import Address Table";
            break;
        case 14:
            RVAChunk.Text = "CLR Header";
            break;
        default:
            RVAChunk.Text = "Data Directory # " + i;
            break;
        }

        var VirtualAddress = new Chunk();
        VirtualAddress.length = 4;
        VirtualAddress.offset = RVAChunk.offset;
        BaseStream.Position = VirtualAddress.offset;
        int VirtualAddressValue = ReadInt32();
        VirtualAddress.Text = "Virtual Address = " + VirtualAddressValue.ToString("x8");
        RVAChunk.Nodes.Add(VirtualAddress);

        var RVASize = new Chunk();
        RVASize.length = 4;
        RVASize.offset = VirtualAddress.offset + VirtualAddress.length;
        BaseStream.Position = RVASize.offset;
        var RVASizeValue = (uint)ReadInt32();

        if (RVASizeValue == 0)
            RVAChunk.Text = "Empty";

        RVASize.Text = "Size = " + RVASizeValue;
        RVAChunk.Nodes.Add(RVASize);

        // TODO: create list of DataDirectory entries. later; calculate their physical offsets from virtual offsets

        var dd = new SectionPointer();
        dd.Text = RVAChunk.Text;
        dd.virtualOffset = VirtualAddressValue;
        dd.length = RVASizeValue;
        dd.realOffset = 0; // TODO calc later

        if (RVASizeValue > 0) {
            dataDirectory.Add(dd);
            DataDirectory.Nodes.Add(RVAChunk);
        }
    }

    optHead.Nodes.Add(DataDirectory);

    peHead.length = (uint)((optHead.offset + optHead.length) - peHead.offset);

    var SectionsOverview = new Chunk();
    SectionsOverview.length = NumberOfSectionsValue * 40;
    SectionsOverview.offset = peHead.offset + peHead.length;
    SectionsOverview.Text = "Sections";

    if (SizeOfOptionalHeaderValue > 0)
        peHead.Nodes.Add(optHead);

    if (NumberOfSectionsValue > 0)
        peHead.Nodes.Add(SectionsOverview);

    for (int i = 0; i < NumberOfSectionsValue; i++) {
        var SectionChunk = new Chunk();
        SectionChunk.length = 40;
        SectionChunk.offset = peHead.offset + peHead.length + (i * SectionChunk.length);

        // Section Name - common names are .text .data .bss
        var SectionName = new ZeroTerminatedStringChunk("Section name", 8);
        SectionName.offset = SectionChunk.offset;
        SectionChunk.Nodes.Add(SectionName);

        var SectioNameValue = "XXX FIXME AGAIN"; // SectionName.GetString(d);

        SectionChunk.Text = "Section " + SectioNameValue;

        // Size of the section once it is loaded to memory
        var SectionSize = SectionName.RelativeToLittleEndian32("Section size loaded");
        SectionChunk.Nodes.Add(SectionSize);

        // RVA (location) of section once it is loaded to memory
        var RVALocation = SectionSize.RelativeToLittleEndian32("RVA location of section");
        BaseStream.Position = RVALocation.offset;
        int RVALocationValue = ReadInt32();
        SectionChunk.Nodes.Add(RVALocation);

        // Physical size of section on disk
        var PhysSize = RVALocation.RelativeToLittleEndian32("Physical size of section");
        BaseStream.Position = PhysSize.offset;
        var PhysSizeValue = (uint)ReadInt32();
        SectionChunk.Nodes.Add(PhysSize);

        // Physical location of section on disk (from start of disk image)
        var PhysOffset = PhysSize.RelativeToLittleEndian32("Physical offset of section");
        BaseStream.Position = PhysOffset.offset;
        int PhysOffsetValue = ReadInt32();
        SectionChunk.Nodes.Add(PhysOffset);

        // Reserved (usually zero) (used in object formats)
        var Reserved12 = PhysOffset.RelativeTo("Reserved", 12);
        SectionChunk.Nodes.Add(Reserved12);

        // Section flags
        var SectionFlags = Reserved12.RelativeToLittleEndian32("Section flags");
        BaseStream.Position = SectionFlags.offset;
        int SectionFlagsValue = ReadInt32();

        if ((SectionFlagsValue & 0x00000020) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000020 = Code";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x00000040) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000040 = Initialized data";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x00000080) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000080 = Uninitialized data";
            SectionFlags.Nodes.Add(note);
        }

        if ((SectionFlagsValue & 0x00000200) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x00000200 = Info";
            SectionFlags.Nodes.Add(note);
        }

        if ((SectionFlagsValue & 0x02000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x02000000 = Discardable";
            SectionFlags.Nodes.Add(note);
        }

        if ((SectionFlagsValue & 0x10000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x10000000 = Shared";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x20000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x20000000 = Executable";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x40000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x40000000 = Readable";
            SectionFlags.Nodes.Add(note);
        }
        if ((SectionFlagsValue & 0x80000000) != 0) {
            var note = new Chunk();
            note.offset = SectionFlags.offset;
            note.length = SectionFlags.length;
            note.Text = "0x80000000 = Writeable";
            SectionFlags.Nodes.Add(note);
        }

        SectionChunk.Nodes.Add(SectionFlags);


        var sectionPointer = new SectionPointer();
        sectionPointer.length = PhysSizeValue;
        sectionPointer.realOffset = PhysOffsetValue;
        sectionPointer.virtualOffset = RVALocationValue;
        sectionPointer.Text = SectioNameValue;
        if (sectionPointer.length > 0)
            sections.Add(sectionPointer);

        SectionsOverview.Nodes.Add(SectionChunk);
    }

    peHead.length += SectionsOverview.length;

    return peHead;
}
*/
