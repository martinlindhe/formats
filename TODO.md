# TODO cmd/formats


hi:
DOCX, ELF,
APK, JAR, MSI, ODT, PDB
PPS, PPT, PPTX, XLS, XLSX

lo:
AXML, DEX, EOT, JSE,
PFB, RAW, T1, T2, TTC, VBE



# TODO cmd/prober

  improve --short output:
  in ShortPrint():
      append some relevant info, depending on kind of file:
        - IMAGES: width + height + bpp
        - ARCHIVES: files in archive, total expanded size

    1st: mark all parsers of format group:
