# Go parameters
GO ?=	go
GODEP ?=	godep
GOBUILD ?=	$(GO) build
GOCLEAN ?=	$(GO) clean
GOINSTALL ?=	$(GO) install
GOTEST ?=	$(GO) test
GOFMT ?=	gofmt -w -s
GODIR ?=	github.com/c14-cli

NAME =		c14

SOURCES :=	$(shell find . -type f -name "*.go" | grep -vE '^./vendor')
COMMANDS :=	$(shell go list ./... | grep -v /vendor/ | grep /cmd/)
PACKAGES :=	$(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)
REV =		$(shell git rev-parse --short HEAD 2> /dev/null || echo "commit")
LDFLAGS = "-X `go list ./pkg/version`.GITCOMMIT=$(REV) -s"

# Check go version
GOVERSIONMAJOR = $(shell go version | grep -o '[1-9].[0-9].' | cut -d '.' -f1)
GOVERSIONMINOR = $(shell go version | grep -o '[1-9].[0-9].' | cut -d '.' -f2)
VERSION_GE_1_6 = $(shell [ $(GOVERSIONMAJOR) -gt 1 -o $(GOVERSIONMINOR) -ge 6 ] && echo true)
ifneq ($(VERSION_GE_1_6),true)
 $(error Bad go version, please install a version greater than or equal to 1.6)
endif

CLEAN_LIST =		$(foreach int, $(COMMANDS) $(PACKAGES), $(int)_clean)
INSTALL_LIST =		$(foreach int, $(COMMANDS), $(int)_install)
TEST_LIST =		$(foreach int, $(COMMANDS) $(PACKAGES), $(int)_test)
COVERPROFILE_LIST =	$(foreach int, $(subst $(GODIR),./,$(PACKAGES)), $(int)/profile.out)

.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(IREF_LIST)

all: build
build: $(NAME)
clean: $(CLEAN_LIST)
	rm -rf $(NAME)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
fmt: $(FMT_LIST)

.git:
	touch $@

$(NAME): $(SOURCES)
	$(GOFMT) $(SOURCES)
	$(GO) tool vet --all=true $(SOURCES)
	$(GOBUILD) -ldflags $(LDFLAGS) ./cmd/c14

$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(subst $(GODIR),./,$*)

$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(subst $(GODIR),./,$*)

$(TEST_LIST): %_test:
	$(GOTEST) -ldflags $(LDFLAGS) -v $(subst $(GODIR),.,$*)


.PHONY: golint
golint:
	@$(GO) get github.com/golang/lint/golint
	@for dir in $(shell $(GO) list ./... | grep -v /vendor/); do golint $$dir; done


.PHONY: gocyclo
gocyclo:
	go get github.com/fzipp/gocyclo
	gocyclo -over 15 $(shell find . -name "*.go" -not -name "*test.go" | grep -v /vendor/)


.PHONY: godep-save
godep-save:
	go get github.com/tools/godep
	$(GODEP) save $(PACKAGES) $(COMMANDS)


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	$(GOENV) goconvey -cover -port=9042 -workDir="$(realpath .)/pkg" -depth=-1


.PHONY: cover
cover: profile.out


$(COVERPROFILE_LIST):: $(SOURCES)
	rm -f $@
	$(GOCOVER) -ldflags $(LDFLAGS) -coverpkg=./pkg/... -coverprofile=$@ ./$(dir $@)

profile.out:: $(COVERPROFILE_LIST)
	rm -f $@
	echo "mode: set" > $@
	cat ./pkg/*/profile.out | grep -v mode: | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $@

goxc:
	rm -rf dist/$(shell cat .goxc.json| jq -r .PackageVersion)
	rm -f dist/latest
	mkdir -p dist/$(shell cat .goxc.json| jq -r .PackageVersion)
	ln -s -f $(shell cat .goxc.json| jq -r .PackageVersion) dist/latest

	$(GOENV) goxc -build-ldflags $(LDFLAGS)

	-mv dist/latest/darwin_386/c14         dist/latest/c14-Darwin-i386
	-mv dist/latest/darwin_amd64/c14       dist/latest/c14-Darwin-amd64
	-mv dist/latest/freebsd_386/c14        dist/latest/c14-Freebsd-i386
	-mv dist/latest/freebsd_amd64/c14      dist/latest/c14-Freebsd-x86_64
	-mv dist/latest/freebsd_arm/c14        dist/latest/c14-Freebsd-arm
	-mv dist/latest/linux_386/c14          dist/latest/c14-Linux-i386
	-mv dist/latest/linux_amd64/c14        dist/latest/c14-Linux-x86_64
	-mv dist/latest/linux_arm/c14          dist/latest/c14-Linux-arm
	-mv dist/latest/windows_386/c14.exe    dist/latest/c14-Windows-i386.exe
	-mv dist/latest/windows_amd64/c14.exe  dist/latest/c14-Windows-x86_64.exe

	-cp dist/latest/c14-Linux-arm dist/latest/c14-Linux-armv7l

	@rmdir dist/latest/* || true

	@echo "Now you can run 'goxc publish-github'"

.PHONY: show_version
show_version:
	./c14 version
