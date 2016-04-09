package parse

// STATUS: 0%

/*
using System;
using System.Collections.Generic;
using System.IO;

namespace MetaEmu
{
    public class TiffReader : SpecificFormatReader
    {
        public TiffReader(FileStream fs) : base(fs)
        {
            name = "TIFF";
            extensions = ".tif; .tiff";
        }

        override public bool IsRecognized()
        {
            BaseStream.Position = 0;

            // XXX dont know magic numbers just guessing
            if (ReadByte() != 'I' || ReadByte() != 'I' || ReadByte() != '*' || ReadByte() != 0)
                return false;

            return true;
        }

        override public List<Chunk> GetFileStructure()
        {
            if (!IsRecognized())
                throw new Exception("not a tiff");

            List<Chunk> res = new List<Chunk>();

            var header = new Chunk();
            header.offset = 0;
            header.length = 4;
            header.Text = "TIFF identifier";
            res.Add(header);

            return res;
        }
    }
}
*/
