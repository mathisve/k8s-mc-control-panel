kubectl delete secret frontend-proxy-config -n mc || true
kubectl create secret generic frontend-proxy-config --from-file=oauth2-proxy.cfg -n mc
kubectl rollout restart deploy/frontend-proxy -n mc