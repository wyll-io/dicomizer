# Dicomizer

Tool suite for DICOM file manipulation, anonymization and data exchange over HTTP(s).

## Package structure

- `pkg`: public packages available for public usage
  - `anonymize`: DICOM anonymization
  - `web`: DICOM over HTTP(s)
- `cmd`: command line tools
- `internal`: specific packages for internal usage

## Useful links

- [DICOM over HTTP(s)](https://dicom.nema.org/medical/dicom/current/output/pdf/part18.pdf)
- [DICOM Confidentiality Profiles](https://dicom.nema.org/medical/dicom/current/output/pdf/part15.pdf)
