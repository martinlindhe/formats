package parse

/*
public IntelMicrocode(FileStream fs) : base(fs)
{
    name = "Intel Microcode";
    extensions = ".dat";
}

override public bool IsRecognized()
{
    BaseStream.Position = 0;

    // version = 1
    var ver = ReadUInt32();
    if (ver != 1)
        return false;

    return true;
}

override public List<Chunk> GetFileStructure()
{
    List<Chunk> res = new List<Chunk>();
    // based on doc in Intel Software Developer Vol 3, 9.11.1

    var header = new Chunk("Intel Microcode header", 48);
    header.offset = 0;
    res.Add(header);

    var version = new Chunk("Version", 4);
    version.offset = 0;
    header.Nodes.Add(version);

    var update = version.RelativeTo("Update revision", 4);
    header.Nodes.Add(update);

    var date = update.RelativeTo("Date", 4);
    header.Nodes.Add(date);

    BaseStream.Position = date.offset;
    var dateString = ReadUInt32().ToString("x8");

    var monthValue = Convert.ToInt32(dateString.Substring(0, 2), 10);
    var dayValue = Convert.ToInt32(dateString.Substring(2, 2), 10);
    var yearValue = Convert.ToInt32(dateString.Substring(4, 4), 10);
    date.Text += " = " + yearValue.ToString("D4") + "-" + monthValue.ToString("D2") + "-" + dayValue.ToString("D2");

    var processorSignature = date.RelativeTo("Processor signature", 4);
    header.Nodes.Add(processorSignature);

    // TODO: Checksum is correct when the summation of all the DWORDs (including the extended Processor Signature Table) that comprise the microcode update result in 00000000H.
    var checksum = processorSignature.RelativeTo("Checksum", 4);
    header.Nodes.Add(checksum);

    var loaderRev = checksum.RelativeTo("Loader revision", 4);
    header.Nodes.Add(loaderRev);

    var processorFlags = loaderRev.RelativeTo("Processor flags", 4);
    header.Nodes.Add(processorFlags);


    var dataSize = processorFlags.RelativeTo("Data size", 4);
    header.Nodes.Add(dataSize);

    BaseStream.Position = dataSize.offset;
    var dataSizeValue = ReadUInt32();
    Console.WriteLine("datasize val = " + dataSizeValue.ToString("x8"));
    if (dataSizeValue == 0) {
        // If this value is 00000000H, then the microcode update encrypted data is 2000 bytes
        dataSizeValue = 0x2000;
    }

    var totalSize = dataSize.RelativeTo("Total size", 4);
    header.Nodes.Add(totalSize);

    var reserved = totalSize.RelativeTo("Reserved", 12);
    header.Nodes.Add(reserved);


    var updateData = reserved.RelativeTo("Update data", dataSizeValue);
    header.Nodes.Add(updateData);


    var extSignatureCount = updateData.RelativeTo("Extended Signature Count", 4);
    header.Nodes.Add(extSignatureCount);

    BaseStream.Position = extSignatureCount.offset;
    var extSignatureCountValue = ReadUInt32();
    Console.WriteLine("extSignatureCount = " + extSignatureCountValue);


    var extChecksum = extSignatureCount.RelativeTo("Extended Checksum", 4);
    header.Nodes.Add(extChecksum);

    var reserved2 = extChecksum.RelativeTo("Reserved (???)", 12);
    header.Nodes.Add(reserved2);

    for (int n = 0; n < extSignatureCountValue; n++) {
        var baseOffset = reserved2.offset + reserved2.length + (n * 12);

        var processorSignatureN = new Chunk("Processor signature #" + n, 4);
        processorSignatureN.offset = baseOffset;
        header.Nodes.Add(processorSignatureN);

        var processorFlagsN = processorSignatureN.RelativeTo("Processor flags #" + n, 4);
        header.Nodes.Add(processorFlagsN);

        var checksumN = processorFlagsN.RelativeTo("Checksum #" + n, 4);
        header.Nodes.Add(checksumN);
    }

    return res;
}
*/
