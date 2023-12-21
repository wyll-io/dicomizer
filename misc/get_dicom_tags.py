import json

import requests
from bs4 import BeautifulSoup, ResultSet, Tag

dicom_tags_url = "https://www.dicomlibrary.com/dicom/dicom-tags"

rsp = requests.get(dicom_tags_url)
if not rsp.ok:
    print("Error: Failed to get DICOM tags")
    exit(1)

soup = BeautifulSoup(rsp.text, "lxml")
table = soup.find(attrs={"id": "table1"})
if type(table) is not Tag:
    print("Error: error type")
    exit(1)


trs: ResultSet[Tag] = table.find_all("tr")
# * delete
del trs[0:2]

tags = []

for tr in trs:
    tds: ResultSet[Tag] = tr.find_all("td")
    if len(tds) > 4:
        print(f"error: {tds[0].text} has more than 4 data elements")
        exit(1)

    group, el = tds[0].text.replace("(", "").replace(")", "").split(",")
    description = tds[2].text.replace("-", " ")
    description = "".join(x for x in description.title() if not x.isspace())

    tags.append(
        {
            "tag": {"group": group, "element": el},
            "vr": tds[1].text,
            "description": description,
            "retired": tds[3].text.lower() == "retired",
        }
    )

with open("tags.json", "w") as f:
    json.dump(tags, f)
