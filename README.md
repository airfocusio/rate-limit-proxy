# rate-limit-proxy

A small go HTTP proxy that can run as sidecar in [Kubernetes](https://kubernetes.io/de/) pods to add HTTP rate limiting:

* Use shared request counter based on [Redis](https://redis.io/).
* Different user identification strategies which can be combined:
    * IP which is the default fallback.
    * [JWT](https://jwt.io/) which you can configure to fit your environment:
        * Where to extract the JWT from? So far Authorization Bearer header custom query parameter are supported.
        * Which signature algorithm to use? So far HSxxx, RSxxx and ESxxx are supported.
        * Which JWT claim to use for identification?
        * Which JWT [kid](https://tools.ietf.org/html/rfc7515#section-4.1.4) to match?
    * more can be easily added
* Specify different limits for different user:
    * One default limit for anonymous users.
    * One default limit for all identified users.
    * Special limits depending on the identified user.

## Try it out

```bash
export JWT_USER="eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyIn0.okfJTi3nwcSI2WITtYXRo8NX7JLd-xqW9iYP7smS2Co"
export JWT_SYSTEM="eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzeXN0ZW0ifQ.-L6_PMWjva1HxRnhGN1ZhfI5PGnmHNrGwA11ndZD6fI"

# make sure go is installed
go version
make run
curl -s -I -XGET -H "Host: golang.org" localhost:8080
curl -s -I -XGET -H "Host: golang.org" -H "Authorization: Bearer $JWT_USER" localhost:8080
curl -s -I -XGET -H "Host: golang.org" -H "Authorization: Bearer $JWT_SYSTEM" localhost:8080

# make sure minikube is running
minikube start
kubectl apply -f example.kubernetes.yaml
kubectl port-forward svc/nginx 8080:http-public
curl -s -I -XGET localhost:8080
curl -s -I -XGET -H "Authorization: Bearer $JWT_USER" localhost:8080
curl -s -I -XGET -H "Authorization: Bearer $JWT_SYSTEM" localhost:8080
kubectl port-forward svc/nginx 8080:http
curl -s -I -XGET localhost:8080
```
