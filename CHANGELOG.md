<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.1"></a>
## [v1.0.1] - 2024-04-12
### Chore
- **config:** move nosq-workbench schema
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


[Unreleased]: https://github.com/wyll-io/dicomizer/compare/v1.0.1...HEAD
[v1.0.1]: https://github.com/wyll-io/dicomizer/compare/v1.0.0...v1.0.1
