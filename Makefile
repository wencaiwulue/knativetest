.PHONY: all
all: admissionwebhook knative

.PHONY: admissionwebhook
admissionwebhook:
	./build/build.sh admissionwebhook

.PHONY: knative
knative:
	./build/build.sh knative
