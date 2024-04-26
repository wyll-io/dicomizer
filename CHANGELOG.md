<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.5"></a>
## [v1.0.5] - 2024-04-26
### Chore
- **docs:** add link to AWS CLI documentation
- **docs:** add cli & aws docs


<a name="v1.0.4"></a>
## [v1.0.4] - 2024-04-26
### Chore
- **changelog:** release v1.0.4

### Fix
- **s3:** fix examen date in folder name


<a name="v1.0.3"></a>
## [v1.0.3] - 2024-04-25
### Chore
- **changelog:** release v1.0.3
- **docker:** use debian:bookworm-slim
- **docs:** update cli docs for crontab flag
- **docs:** update docs according to new arguments

### Fix
- **anonymize:** add missing tag
- **cli:** set crontab as flag
- **cli:** update args
- **s3:** add subfolder for patient's dcm


<a name="v1.0.2"></a>
## [v1.0.2] - 2024-04-15
### Chore
- **changelog:** release v1.0.2

### Fix
- **S3:** fix data sent to S3


<a name="v1.0.1"></a>
## [v1.0.1] - 2024-04-15
### Chore
- **config:** move nosq-workbench schema
- **release:** release v1.0.1
- **tools:** update flake.nix, css and htmx version

### Doc
- **readme:** add doc link

### Feat
- **error:** update error message

### Fix
- **anonymize:** deep sequence anonymize
- **check:** read file for hash check
- **hash:** disable has verification for investigation


<a name="v1.0.0"></a>
## v1.0.0 - 2024-02-22
### Chore
- **ci:** add Dockerfile
- **ci:** add git-chglog & goreleaser
- **doc:** add flow and anonymization
- **docker:** update dcmtk extraction
- **release:** release v1.0.0
- **utils:** add python script to retrieve/export DICOM tags

### Feat
- **DICOM:** add DICOM DISME request
- **anonymizer:** add anonymization methods
- **aws:** connect dynamodb to frontend
- **aws:** add dynamodb
- **aws:** add db, glacier & s3 connectors
- **gui:** add web server
- **misc:** add nosql-workbench dynamodb schema
- **scheduler:** add scheduler function
- **storage:** create storage interface

### Fix
- **gui:** set correct ID when canceling editing
- **gui:** form error when validating add patient
- **gui:** add htmx 400 response add_patient


[Unreleased]: https://github.com/wyll-io/dicomizer/compare/v1.0.5...HEAD
[v1.0.5]: https://github.com/wyll-io/dicomizer/compare/v1.0.4...v1.0.5
[v1.0.4]: https://github.com/wyll-io/dicomizer/compare/v1.0.3...v1.0.4
[v1.0.3]: https://github.com/wyll-io/dicomizer/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/wyll-io/dicomizer/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/wyll-io/dicomizer/compare/v1.0.0...v1.0.1
