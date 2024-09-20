# k8s-mc-control-panel
A control panel for a Minecraft server running on K8s

## frontend-control
A react front end that polls the backend server.

## backend
A golang backend that talks to the KubeAPI server and responds to HTTP requests from the frontend-control.

## Current TODOs
- Remove static ips, replace with something better
- Templatize deployment with Helm
- Add RCON proxy container to talk to the minecraft server directly
- Add user list in UI