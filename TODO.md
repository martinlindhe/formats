# TODO cmd/formats


ELF exec
MSI exec
PDB debug info

VBE "VBScript Encoded Script File" (need samples)

RAW image


# TODO zip-based formats:

DOCX, PPTX, XLSX
    https://en.wikipedia.org/wiki/Office_Open_XML

ODT, ODP, ODS
    https://en.wikipedia.org/wiki/OpenDocument

APK (android application package, JAR-based)
    https://en.wikipedia.org/wiki/Android_application_package

JAR (java archive)
    https://en.wikipedia.org/wiki/JAR_%28file_format%29


# TODO
    .doc, .pps, .ppt, .xls: all are detected as "word" documents now.
        improve parser to distinguish between them


# TODO cmd/prober

  improve --short output:
  in ShortPrint():
      append some relevant info, depending on kind of file:
        - IMAGES: width + height + bpp
        - ARCHIVES: files in archive, total expanded size

    1st: mark all parsers of format group:
