# Installation

## Docker

First pull the docker image:
```bash
docker pull ghcr.io/wyll-io/dicomizer:TAG # Replace TAG with the desired tag. `latest` is valid.
```

Then, run the image with the desired configuration
```bash
docker run --name dicomizer \
    -e JWT_SECRET=demo \
    -e ADMIN_PASSWORD=demo \
    -e AWS_ACCESS_KEY_ID=xxx \
    -e AWS_SECRET_ACCESS_KEY=xxx \
    -e AWS_REGION=eu-west-3 \
    ghcr.io/wyll-io/dicomizer:TAG \
    command # Replace with the desired command (start, anonymize, etc)
```

## Binary

Download the desired version from the [GitHub release section](https://github.com/wyll-io/dicomizer/releases),
place the extracted directory where you want and add it to your path.

Example:
```bash
wget https://github.com/wyll-io/dicomizer/releases/download/v1.0.4/dicomizer_Linux_x86_64.tar.gz
tar xzf dicomizer_Linux_x86_64.tar.gz -C /usr/local/dicomizer
export PATH=$PATH:/usr/local/dicomizer

dicomizer --help
```

### Systemd

After downloading and extracting the release (instructions above), create a new systemd service file:

```bash
sudo touch /etc/systemd/system/dicomizer.service
```

Copy the following content inside the newly created service file:

```
[Unit]
Description="Dicomizer service"
After=network.target

[Service]
Type=simple
Restart=on-failure
RestartSec=5
# Don't forget to customize the parameters
ExecStart=/usr/local/dicomizer/dicomizer start --aec AEC --aem DICOMIZER_QR --aet DICOMIZER_QR --pacs 0.0.0.0 --center "CENTER" --bind 0.0.0.0:80 --crontab "*/2 * * * *"
WorkingDirectory=/usr/local/dicomizer

[Install]
WantedBy=multi-user.target
```

Don't forget to enable your service:
```bash
sudo systemctl enable dicomizer.service
```

After, create an override file to set your environment variables:
```bash
sudo systemctl edit dicomizer.service
```

And append the following environment variables:
```
[Service]
Environment="AWS_ACCESS_KEY_ID=XXXX"
Environment="AWS_SECRET_ACCESS_KEY=XXX"
Environment="AWS_REGION=xxx"
Environment="JWT_SECRET=demo"
Environment="ADMIN_PASSWORD=demo"
```

And reload your systemd daemon to take changes:
```bash
sudo systemctl daemon-reload
```

You should be good to go!
