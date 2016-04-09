package parse

// STATUS: 0%

/*
using System;
using System.Collections.Generic;
using System.IO;

namespace MetaEmu
{
    public class JpegReader : SpecificFormatReader
    {
        public JpegReader(FileStream fs) : base(fs)
        {
            name = "JPEG";
            extensions = ".jpg; .jpeg";
            mimetype = "image/jpeg";
        }

        override public bool IsRecognized()
        {
            BaseStream.Position = 0;

            if (ReadByte() != 0xFF || ReadByte() != 0xD8)
                return false;

            // skip 4 bytes in data stream
            ReadInt32();

            if (ReadByte() != 'J' || ReadByte() != 'F' || ReadByte() != 'I' || ReadByte() != 'F' || ReadByte() != 0)
                return false;

            return true;
        }

        override public List<Chunk> GetFileStructure()
        {
            int count = 0;
            List<Chunk> res = new List<Chunk>();

            if (!IsRecognized())
                throw new Exception("not a jpeg");

            BaseStream.Position = 0;

            do {
                count++;

                var marker = new Chunk();
                marker.offset = BaseStream.Position;

                byte marker0 = ReadByte();
                byte marker1 = ReadByte();

                if (marker0 != 0xFF) {
                    Log("parse error, found 0x" + marker0.ToString("x2") + " at offset " + (BaseStream.Position - 2).ToString("x4"));
                    break;
                }

                marker.Text = "Type 0x" + marker1.ToString("x4");
                Log("Found marker " + marker1.ToString("x2"));

                var type = new BigEndian16BitChunk("Type");
                type.offset = marker.offset;
                marker.Nodes.Add(type);

                if (marker1 == 0xD8) {
                    // NOTE: this marker dont have any content
                    marker.Text = "SOI - Start of Image";
                    marker.length = 2;
                    res.Add(marker);
                    continue;
                }
                if (marker1 == 0xD9) {
                    marker.Text = "EOI - End of Image";
                    marker.length = 2;
                    res.Add(marker);

                    Log("Ending parser since EOI marker was detected");
                    break;
                }

                var length = type.RelativeToBigEndian16("Length");
                marker.Nodes.Add(length);

                uint lenghtValue = (uint)ReadInt16BE(length.offset);

                marker.length = 2 + lenghtValue;
                Log("len = " + marker.length);


                var data = length.RelativeTo("Data", lenghtValue - 2);
                marker.Nodes.Add(data);

                switch (marker1) {
                case 0xC0:
                    marker.Text = "SOF0 - Baseline DCT";
                    break;
                case 0xC1:
                    marker.Text = "SOF1 - Extended sequential DCT";
                    break;
                case 0xC2:
                    marker.Text = "SOF2 - Progressive DCT";
                    break;
                case 0xC3:
                    marker.Text = "SOF3 - Lossless (sequential)";
                    break;
                case 0xC4:
                    marker.Text = "DHT - Huffman table";
                    break;

                case 0xDA:
                    marker.Text = "SOS - Start of scan";

                    int component_count = ReadByte(data.offset);

                    var Components = length.RelativeToByte("Color components");
                    data.Nodes.Add(Components);

                    for (int i = 0; i < component_count; i++) {
                        //        For each component
                        //        An ID
                        //        An AC table # (Low Nibble)
                        //        An DC table # (High Nibble)

                        var kex = new BigEndian16BitChunk("Color");
                        kex.offset = data.offset + 1 + (i * kex.length);
                        data.Nodes.Add(kex);
                    }

                    var Unknown = new Chunk();
                    Unknown.offset = data.offset + 1 + (2 * component_count);
                    Unknown.length = 3;
                    Unknown.Text = "Unknown";
                    data.Nodes.Add(Unknown);

                    // TODO: now follows compressed image data, how to get length of data?

                    break;

                case 0xDB:
                    marker.Text = "DQT - Quantization table";
                    break;

                case 0xE0:
                    marker.Text = "APP0 - Application Use";

                    var Identifier = length.RelativeToZeroTerminatedString("Id String", 5);
                    marker.Nodes.Add(Identifier);

                    var Version = Identifier.RelativeToVersionMajorMinor16("Revision");
                    marker.Nodes.Add(Version);

                    // Units used for Resolution
                    var Units = Version.RelativeToByte("Units used");
                    marker.Nodes.Add(Units);

                    var Width = Units.RelativeToBigEndian16("Width");
                    marker.Nodes.Add(Width);

                    var Height = Width.RelativeToBigEndian16("Height");
                    marker.Nodes.Add(Height);

                    var XThumbnail = Height.RelativeToByte("Horizontal Pixels");
                    marker.Nodes.Add(XThumbnail);

                    var YThumbnail = XThumbnail.RelativeToByte("Vertical Pixels");
                    marker.Nodes.Add(YThumbnail);
                    break;

                case 0xE1:
                    marker.Text = "APP1 - Application Use";
                    // XXXX
                    break;

                case 0xFE:
                    marker.Text = "COM - Comment";
                    break;

                default:
                    throw new Exception("TODO " + marker1.ToString("x2"));
                }

                BaseStream.Position += lenghtValue - 2;
                res.Add(marker);

            } while (count < 19);

            return res;
        }

        public static void Log(string s)
        {
            Console.WriteLine(s);
        }
    }
}
*/
