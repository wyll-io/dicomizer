# GUI

A web server is shipped with the `Dicomizer` package. It is used to manage patients

## Introduction

`Dicomizer` have a GUI (graphical user interface) to manage new patients and how to retrieve their data from the PACS server.

## How to use it

1. Setup the CLI [(docs)](./cli.md)

2. Run (customize the cron expression to your needs and change the IP address to which ip address the server should listen to and its port):

```bash
dicomizer start "0 0 * * *"
```

3. Open your browser and go to: <http://localhost>

4. Login with the password you set in the `.env` file (`ADMIN_PASSWORD`).

## Actions

- Right column, you can add a new patient with an `Identity` and DICOM tags to filter the patient on the PACS server (`;` separated). Filter syntax can either be `(GROUP,ELEMENT)=VALUE` or `GROUP,ELEMENT=VALUE`. Example:

```
Fullname: John Doe
Tags: (0010,0020)=12345678;(0010,0010)=DOE^JOHN
```

- Left column, you can see the list of patients that will be retrieved from the database (not the PACS). You can search for a patient by typing in the search bar. You can also delete a patient by clicking on the trash icon. Or modify its information by clicking on the pencil icon.

- You can disconnect yourself by clicking on the `Se deconnecter` button.
