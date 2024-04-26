BINARY_NAME := "dicomizer"
DOCKER_IMAGE_NAME := "ghcr.io/wyll-io/dicomizer"

latest-tag := `git describe --tags --abbrev=0 | cut -c 1-`

alias bb := build-binary
build-binary: build-css
  go build -o {{BINARY_NAME}} -v .

alias bd := build-docker
build-docker: build-css
  docker build -t {{DOCKER_IMAGE_NAME}} .

alias bc := build-css
build-css:
  pnpm exec tailwindcss -i ./input.css -o ./public/css/tailwind.css

publish version: (bump-files-version version)
  git-chglog --next-tag {{version}} --output CHANGELOG.md
  git add CHANGELOG.md main.go
  git commit -m "chore(changelog): release {{version}}"
  git tag -a {{version}} -m "{{version}}"
  git push --follow-tags

bump-files-version version:
  sed -i 's/{{latest-tag}}/{{version}}/g' main.go
