# Anonymizer

Anonymizer is a tool to anonymize data. It is used to replace sensitive data with a non-sensitive equivalent. It is used to protect the privacy of individuals and to comply with data protection laws.

[Reference](https://dicom.nema.org/dicom/2013/output/chtml/part15/chapter_E.html)

## How it works

It loads into memory the DICOM file and create a dataset. Then, it loops through the dataset and check for known tags that should be anonymized. If a tag is found, it is anonymized with the corresponding value.

### Anonymization content

NOTE: UID are replaced with the root UID of the project followed by 64 uinsigned characters.

Here's the detail of the anonymization content:

- UID: replaced with `1.2.826.0.1.3680043.10.1336.[64_UNSIGNED_CHARACTERS]`
- DA (date): replaced with `00010101`
- DT (date time): replaced with `00010101010101.000000+0000`
- TM (time): replaced with `000000.00`
- any other string value: replaced with `Anonymized`
- Some tags require to be deleted, they're removed from the dataset

### Tags to anonymize

You can find the list of tags to anonymize in the `pkg/anonymize/tags.go` file.
