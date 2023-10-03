FROM registry.k8s.io/e2e-test-images/agnhost:2.40 as builder

FROM gcr.io/distroless/static-debian12

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /agnhost /agnhost

USER 65534

# https://pkg.go.dev/k8s.io/kubernetes/test/images/agnhost#section-readme
ENTRYPOINT ["/agnhost"]
CMD ["netexec"]