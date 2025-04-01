podman-run:
	@echo "Running CodeFlick with Podman..."
	@make podman-build-client-image
	@make podman-build-server-image
	@make podman-pod-create

podman-build-client-image:
	@echo "Building client image with Podman..."
	@podman build -t localhost/codeflick-client:latest ./client/

podman-build-server-image:
	@echo "Building server image with Podman..."
	@podman build -t localhost/codeflick-server:latest ./server/

podman-pod-create:
	@echo "Creating CodeFlick pod..."
	@podman play kube podman-pod.yaml

podman-pod-delete:
	@echo "Deleting CodeFlick pod..."
	@podman pod rm -f codeflick-pod