# Program flow

![Diagram flow](./flow.png)

## Introduction

This document describes the flow of the program.

## GUI flow

### Adding a new user

1. The user fills the patient form with the patient's information (Fullname and DICOM filters).
2. The user clicks on the "Sauvegarder" button.
3. The server receives the form and creates a new patient in the database.
4. The patient is now available in the list of patients.

### Modifying a user

1. The user clicks on the "Modifier" button of the patient he wants to modify.
2. The user modifies the patient's information.
3. The user clicks on the "Sauvegarder" button.
4. The server receives the form and updates the patient in the database.
5. The patient's newly updated information are now available.

### Deleting a user

1. The user clicks on the "Supprimer" button of the patient he wants to delete.
2. The server receives the request and deletes the patient from the database.
3. The patient is no longer available in the list of patients.

## Cron job flow

1. The cron job fetches the list of patients from the database.
2. For each patient, the cron job fetches the DICOM data from the PACS server using
   the patient's DICOM filters.
3. The server check if the DICOM data is already in the database.
   - If the data is not in the database
     1. The server anonymizes the data.
     2. The server sends the anonymized file to the AWS S3 bucket.
     3. The server saves the data in the database.
   - If the data is already in the database
     1. Skip it.
