# hello-world

This is a very basic HTTP application which is aimed to be used as a test application running in Kubernetes, AWS ECS,
or similar.

* Echos back the details of request for every response
* Doesn't run as root in case `runAsNonRoot` is enabled
* Is able to check the server is running using the same binary (`/usr/local/bin/app --health`) to serve as the ECS
  health check.
