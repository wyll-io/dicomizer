# CLI

## Commands

### `start`

Starts the web server and the cron job to fetch automatically DICOM data from a PACS server.

Example:

```bash
# Start the server on port 3000, listening on localhost and fetch data every day at midnight
$ dicomizer start localhost:3000 "0 0 * * *"
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

### `upload`

Upload DICOM files to configured S3 and add it to DynamoDB. The process also includes anonymization.

Example:

```bash
# Upload DICOM files located at /path/to/{file1,file2}.dcm
$ dicomizer upload /path/to/file1.dcm /path/to/file2.dcm
```
