DESTDIR ?= /
VERSION ?= 34.0
TAG = v$(VERSION)
PREVTAG := $(shell git tag --sort=-creatordate | head -n 2 | tail -n 1)
COMMITS := $(shell git log --pretty=oneline --no-merges ${PREVTAG}..HEAD | wc -l)
GPGKEY ?= $(shell git config user.signingkey)
GITEMAIL := $(shell git config user.email)
GITNAME := $(shell git config user.name)
BUILDFLAGS ?= -ldflags="-X github.com/osbuild/weldr-client/cmd/composer-cli/root.Version=${VERSION}"
GOBUILDFLAGS ?= 

build: composer-cli
composer-cli:
	go build ${BUILDFLAGS} ${GOBUILDFLAGS} ./cmd/composer-cli

check:
	go vet ./... && golint -set_exit_status ./... && golangci-lint --build-tags=integration run ./...

test:
	go test ${BUILDFLAGS} ${GOBUILDFLAGS} -v -covermode=atomic -coverprofile=coverage.txt -coverpkg=./... ./...

integration: composer-cli-tests
composer-cli-tests:
	go test -c -tags=integration ${BUILDFLAGS} ${GOBUILDFLAGS} -o composer-cli-tests ./weldr/

install: composer-cli
	install -m 0755 -vd ${DESTDIR}/usr/bin/
	install -m 0755 -vp composer-cli ${DESTDIR}/usr/bin/
	install -m 0755 -vd ${DESTDIR}/etc/bash_completion.d/
	install -m 0755 -vp etc/bash_completion.d/composer-cli ${DESTDIR}/etc/bash_completion.d/
	install -m 0755 -vd ${DESTDIR}/usr/share/man/man1/
	./composer-cli doc ${DESTDIR}/usr/share/man/man1/

install-tests: composer-cli-tests
	install -m 0755 -vd ${DESTDIR}/usr/libexec/tests/composer-cli/
	install -m 0755 -vp composer-cli-tests ${DESTDIR}/usr/libexec/tests/composer-cli/

weldr-client.spec: weldr-client.spec.in
	sed -e "s/%%VERSION%%/$(VERSION)/" < $< > $@
	$(MAKE) -s changelog >> $@

tag:
	@if [ -z "$(GPGKEY)" ]; then echo "ERROR: The git config user.signingkey must be set" ; exit 1; fi
	git tag -u $(GPGKEY) -m "Tag as $(TAG)" -f $(TAG)
	@echo "Tagged as $(TAG)"

# Order matters, so run make for each step instead of declaring them as dependencies
release:
	@if [ -z "$(GPGKEY)" ]; then echo "ERROR: The git config user.signingkey must be set" ; exit 1; fi
	$(MAKE) test && $(MAKE) bumpver && $(MAKE) tag && $(MAKE) archive && $(MAKE) sign

sign:
	@if [ -z "$(GPGKEY)" ]; then echo "ERROR: The git config user.signingkey must be set" ; exit 1; fi
	gpg --armor --detach-sign -u $(GPGKEY) weldr-client-$(VERSION).tar.gz

changelog:
	@echo "* $(shell date '+%a %b %d %Y') ${GITNAME} <${GITEMAIL}> - ${VERSION}-1"
	@git log --no-merges --pretty="format:- %s (%ae)" ${PREVTAG}..HEAD |sed -e 's/@.*)/)/'

bumpver:
	@NEWSUBVER=$$((`echo $(VERSION) |cut -d . -f 2` + 1)) ; \
	NEWVERSION=`echo $(VERSION).$$NEWSUBVER |cut -d . -f 1,3` ; \
	sed -i "s/VERSION ?= $(VERSION)/VERSION ?= $$NEWVERSION/" Makefile ; \
	git add Makefile; \
	git commit -m "New release: $$NEWVERSION"

archive:
	git archive --prefix=weldr-client-$(VERSION)/ --format=tar.gz HEAD > weldr-client-$(VERSION).tar.gz

RPM_SPECFILE=rpmbuild/SPECS/weldr-client.spec
RPM_TARBALL=rpmbuild/SOURCES/weldr-client-$(VERSION).tar.gz
RPM_TARBALL_SIG=rpmbuild/SOURCES/weldr-client-$(VERSION).tar.gz.asc

$(RPM_SPECFILE): weldr-client.spec
	mkdir -p $(CURDIR)/rpmbuild/SPECS
	cp weldr-client.spec $(CURDIR)/rpmbuild/SPECS

$(RPM_TARBALL): archive sign
	mkdir -p $(CURDIR)/rpmbuild/SOURCES
	cp weldr-client-$(VERSION).tar.gz* rpmbuild/SOURCES/

srpm: $(RPM_SPECFILE) $(RPM_TARBALL)
	rpmbuild -bs \
		--define "_topdir $(CURDIR)/rpmbuild" \
		--define "commit $(VERSION)" \
		--with tests \
		$(RPM_SPECFILE)

rpm: $(RPM_SPECFILE) $(RPM_TARBALL)
	dnf builddep -y --spec -D 'with 1' $(RPM_SPECFILE)
	rpmbuild -bb \
		--define "_topdir $(CURDIR)/rpmbuild" \
		--define "commit $(VERSION)" \
		--with tests \
		$(RPM_SPECFILE)


.PHONY: build check test integration install srpm rpm weldr-client.spec
