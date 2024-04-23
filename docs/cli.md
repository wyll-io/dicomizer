# CLI

## Commands

### `start`

Starts the web server and the cron job to fetch automatically DICOM data from a PACS server.

Example:

```bash
# Start the server on port 80 by default, listening on localhost by default and fetch data every day at midnight
$ dicomizer start "0 0 * * *" --pacs "10.0.0.56:104" --aet "MY_AET" --aec "THEIR_AEC" --aem "MY_AEM"
```

### `anonymize`

Anonymize a DICOM file.

Example:

```bash
# Anonymize a DICOM file located at /path/to/file.dcm and output it in `./anonymized.dcm`
$ dicomizer anonymize /path/to/file.dcm
# Anonymize a DICOM file located at /path/to/file.dcm and output it in `/another/path/ano.dcm`
$ dicomizer anonymize /path/to/file.dcm /another/path/ano.dcm
```
