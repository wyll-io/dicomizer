import json
import sys

import jinja2
import requests
from bs4 import BeautifulSoup, ResultSet, Tag

DICOM_LIBRARY_TAGS = "https://www.dicomlibrary.com/dicom/dicom-tags"


def fetch_dicom_tags(output="./dicom_tags.json"):
    print("Fetching DICOM tags...")
    rsp = requests.get(DICOM_LIBRARY_TAGS)
    if not rsp.ok:
        print("Error: Failed to get DICOM tags")
        exit(1)

    print("Parsing DICOM tags...")
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
            continue

        tag_description = tds[2].text
        if len(tag_description) == 0:
            print(f"error: {tds[0].text} has no tag name")
            continue

        group, el = tds[0].text.replace("(", "").replace(")", "").split(",")
        if "x" in group or "x" in el:
            print(f"error: {tds[0].text} has invalid group or element")
            continue

        tag_name = (
            tag_description.replace("-", " ")
            .replace("(", " ")
            .replace(")", " ")
            .replace("'", "")
            .replace("/", " ")
        )
        tag_name = "".join(x for x in tag_name.title() if not x.isspace())

        tags.append(
            {
                "tag": {"group": group, "element": el},
                "vr": tds[1].text,
                "name": tag_name,
                "retired": tds[3].text.lower() == "retired",
            }
        )

    with open(output, "w") as f:
        json.dump(tags, f)

    print("Done")


def gen_go_file(
    package_name: str,
    input="./dicom_tags.json",
    output="./dicom_tags.go",
    accept_retired="false",
):
    print(f'Generating Go file with "{input}"...')
    env = jinja2.Environment(loader=jinja2.FileSystemLoader(searchpath="./"))

    template = env.get_template("go.jinja2")
    with open(input, "r") as f:
        with open(output, "w") as out:
            out.write(
                template.render(
                    {
                        "tags": json.load(f),
                        "package_name": package_name,
                        "accept_retired": accept_retired == "true",
                    }
                )
            )

    print(f'Go file generated at "{output}"')


if __name__ == "__main__":
    args = sys.argv[1:]
    if args[0] == "fetch_tags":
        fetch_dicom_tags(*args[1:])
    elif args[0] == "gen_go":
        gen_go_file(*args[1:])
