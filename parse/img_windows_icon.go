package parse

// STATUS xxx

/*
using System;
using System.Collections.Generic;
using System.IO;

namespace MetaEmu
{
    public class IconReader : SpecificFormatReader
    {
        public IconReader(FileStream fs) : base(fs)
        {
            name = "Windows Icon";
            extensions = ".ico";
        }

        override public bool IsRecognized()
        {
            BaseStream.Position = 0;

            if (ReadByte() != 0 || ReadByte() != 0)
                return false;

            byte type = ReadByte();
            // 1 = icon, 2 = cursor
            if (type != 1 && type != 2)
                return false;

            if (ReadByte() != 0)
                return false;

            return true;
        }

        override public List<Chunk> GetFileStructure()
        {
            BaseStream.Position = 0;

            if (!IsRecognized())
                throw new Exception("not a ico");

            BaseStream.Position = 0;

            byte type = ReadByte(2);

            List<Chunk> res = new List<Chunk>();

            var header = new Chunk();
            header.offset = 0;
            header.Text = "ICON header";

            var identifier = new Chunk();
            identifier.offset = 0;
            identifier.length = 4;
            if (type == 1) {
                identifier.Text = "ICO identifier";
            } else if (type == 2) {
                identifier.Text = "CUR identifier";
            } else {
                throw new Exception("unknown " + type);
            }
            header.Nodes.Add(identifier);

            var numIcons = identifier.RelativeToLittleEndian16("Resources");
            int numIconsValue = ReadInt16(numIcons.offset);
            header.Nodes.Add(numIcons);

            res.Add(header);

            uint iconEntryLength = 16;

            Log("parsing " + numIconsValue + " resources");

            for (int i = 0; i < numIconsValue; i++) {
                var iconEntry = new Chunk("Resource entry #" + (i + 1));
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
            }

            header.length = (uint)(6 + (numIconsValue * iconEntryLength));

            return res;
        }

        public static void Log(string s)
        {
            Console.WriteLine("[Icon] " + s);
        }
    }
}

*/
