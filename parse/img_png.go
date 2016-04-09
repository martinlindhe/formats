package parse

// STATUS: 0%

/*
using System;
using System.Collections.Generic;
using System.IO;

namespace MetaEmu
{
    public class PngReader : SpecificFormatReader
    {
        public PngReader(FileStream fs) : base(fs)
        {
            extensions = ".png; .mng";
        }

        private bool HasPngHeader()
        {
            BaseStream.Position = 0;

            if (ReadByte() != 0x89 || ReadByte() != 'P' || ReadByte() != 'N' || ReadByte() != 'G' ||
                ReadByte() != 0x0D || ReadByte() != 0x0A || ReadByte() != 0x1A || ReadByte() != 0x0A)
                return false;

            return true;
        }

        private bool HasMngHeader()
        {
            BaseStream.Position = 0;

            if (ReadByte() != 0x8A || ReadByte() != 'M' || ReadByte() != 'N' || ReadByte() != 'G' ||
                ReadByte() != 0x0D || ReadByte() != 0x0A || ReadByte() != 0x1A || ReadByte() != 0x0A)
                return false;

            return true;
        }

        override public bool IsRecognized()
        {
            if (HasPngHeader() || HasMngHeader())
                return true;

            return false;
        }

        private string ReadString(long offset, int length)
        {
            BaseStream.Position = offset;

            string res = "";

            for (int i = 0; i < length; i++) {
                res += ReadChar();
            }

            return res;
        }

        override public List<Chunk> GetFileStructure()
        {
            List<Chunk> res = new List<Chunk>();

            var header = new Chunk();
            if (HasPngHeader()) {
                header.Text = "PNG identifier";
                name = "PNG image";
                mimetype = "image/png";
            } else if (HasMngHeader()) {
                header.Text = "MNG identifier";
                name = "MNG image";
                mimetype = "video/x-mng";
            } else
                throw new Exception("not a png");

            header.offset = 0;
            header.length = 8;
            res.Add(header);

            long offset = header.offset + header.length;

            do {
                var chunk = new Chunk();
                chunk.offset = offset;

                var length = new BigEndian32BitChunk("Length");
                length.offset = offset;
                chunk.Nodes.Add(length);

                BaseStream.Position = length.offset;
                uint lengthVal = ReadUInt32BE();

                var type = length.RelativeToZeroTerminatedString("Type", 4);
                chunk.Nodes.Add(type);

                string typeStr = ReadString(type.offset, 4);

                chunk.Text = "Chunk " + typeStr;
                chunk.length = lengthVal + 4 + 4 + 4;  // "length" (4 byte) + "type" (4 byte) + data + crc (4 byte)

                var data = type.RelativeTo("Data", lengthVal);

                if (lengthVal > 0) {
                    if (typeStr == "IHDR") {
                        var width = type.RelativeToBigEndian32("Width");
                        data.Nodes.Add(width);

                        var height = width.RelativeToBigEndian32("Height");
                        data.Nodes.Add(height);

                        var bd = height.RelativeToByte("Bit depth");
                        data.Nodes.Add(bd);

                        var color = bd.RelativeToByte("Color type");
                        data.Nodes.Add(color);

                        var compression = color.RelativeToByte("Compression method");
                        data.Nodes.Add(compression);

                        var filter = compression.RelativeToByte("Filter method");
                        data.Nodes.Add(filter);

                        var interlace = filter.RelativeToByte("Interlace method");
                        data.Nodes.Add(interlace);
                    }

                    chunk.Nodes.Add(data);
                }

                var crc = data.RelativeToBigEndian32("Crc");
                chunk.Nodes.Add(crc);

                offset += chunk.length;

                res.Add(chunk);

                if (typeStr == "IEND") {
                    Log("Stopped parser after IEND chunk");
                    break;
                }

            } while (offset < BaseStream.Length);

            return res;
        }

        public static void Log(string s)
        {
            Console.WriteLine("[Png] " + s);
        }
    }
}

*/
