BINARY_NAME := "dicomizer"
DOCKER_IMAGE_NAME := "ghcr.io/wyll-io/dicomizer"

alias bb := build-binary
build-binary: build-css
  @echo "Building binary..."
  go build -o {{BINARY_NAME}} -v .
  @echo "Done!"

alias bd := build-docker
build-docker: build-css
  @echo "Building docker image..."
  docker build -t {{DOCKER_IMAGE_NAME}} .
  @echo "Done!"

alias bc := build-css
build-css:
  @echo "Building CSS..."
  pnpm exec tailwindcss -i ./input.css -o ./public/css/tailwind.css