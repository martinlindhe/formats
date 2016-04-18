# TODO cmd/formats

basic detection of the following formats:
  - MSI exec

# TODO
.doc, .pps, .ppt, .xls: all are detected as "word" documents now.
    improve parser to distinguish between them


# TODO zip-based formats:

DOCX, PPTX, XLSX
    https://en.wikipedia.org/wiki/Office_Open_XML

ODT, ODP, ODS
    https://en.wikipedia.org/wiki/OpenDocument

APK (android application package, JAR-based)
    https://en.wikipedia.org/wiki/Android_application_package

JAR (java archive)
    https://en.wikipedia.org/wiki/JAR_%28file_format%29

# TODO RAW images
    http://fileformats.archiveteam.org/wiki/Cameras_and_Digital_Image_Sensors
    https://www.rawsamples.ch/

# TODO cmd/prober

  improve --short output:
  in ShortPrint():
      append some relevant info, depending on kind of file:
        - IMAGES: width + height + bpp
        - ARCHIVES: files in archive, total expanded size
        - exec: show sections

    1st: mark all parsers of format group:
