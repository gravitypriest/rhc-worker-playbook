PYTHON		?= python3.9

PKGNAME=rhc-worker-playbook
PKGVER = $(shell $(PYTHON) setup.py --version | tr -d '\n')
_SHORT_COMMIT = $(shell git rev-parse --short HEAD | tr -d '\n')
_LATEST_TAG = $(shell git describe --tags --abbrev=0 --always | tr -d '\n')
_NUM_COMMITS_SINCE_LATEST_TAG = $(shell git rev-list $(_LATEST_TAG)..HEAD --count | tr -d '\n')
RELEASE = $(shell printf "99.%s.git.%s" $(_NUM_COMMITS_SINCE_LATEST_TAG) $(_SHORT_COMMIT))

DESTDIR		?= /
PREFIX		?= $(DESTDIR)usr/local
LIBDIR		?= $(PREFIX)/lib
LIBEXECDIR	?= $(PREFIX)/libexec
SYSCONFDIR	?= $(PREFIX)/etc
CONFIG_FILE	?= $(SYSCONFDIR)/rhc/workers/rhc-worker-playbook.toml
WORKER_LIB_DIR ?= $(LIBDIR)/$(PKGNAME)

build: rhc_worker_playbook/constants.py
	$(PYTHON) setup.py sdist
	$(PYTHON) -m pip wheel -c constraints.txt --wheel-dir=wheel .

rhc_worker_playbook/constants.py: rhc_worker_playbook/constants.py.in
	sed \
		-e 's,[@]CONFIG_FILE[@],$(CONFIG_FILE),g' \
		-e 's,[@]WORKER_LIB_DIR[@],$(WORKER_LIB_DIR),g' \
		$^ > $@

.PHONY: install
install: 
	for wheel in $(shell ls wheel/ansible*) $(shell ls wheel/grpcio*) $(shell ls wheel/protobuf*) ; do \
		$(PYTHON) -m pip install $$wheel --no-index --find-links wheel --target $(LIBDIR)/$(PKGNAME) ;\
	done
	$(PYTHON) -m pip install $(shell ls wheel/rhc_worker_playbook*) --no-index --no-deps --prefix=$(PREFIX)
	install -Dm 755 $(PREFIX)/bin/rhc-worker-playbook.worker $(LIBEXECDIR)/rhc/workers/rhc-worker-playbook.worker
	rm -rf $(PREFIX)/bin/rhc-worker-playbook.worker
	[[ -e $(CONFIG_FILE) ]] || install -Dm 644 ./rhc-worker-playbook.toml $(CONFIG_FILE)

.PHONY: uninstall
uninstall:
	rm -rf $(LIBEXECDIR)/rhc/$(PKGNAME).worker
	rm -rf $(LIBDIR)/python*/site-packages/$(PKGNAME)*
	rm -rf $(LIBDIR)/$(PKGNAME)
	$(PYTHON) -m pip -y uninstall rhc-worker-playbook
	rm -f $(LIBEXECDIR)/rhc/workers/rhc-worker-playbook.worker


.PHONY: clean
clean:
	rm -f rhc_worker_playbook/constants.py
	rm -f rhc-worker-playbook.spec
	rm -rf wheel

rhc-worker-playbook.spec: rhc-worker-playbook.spec.in
	sed \
		-e 's,[@]PKGVER[@],$(PKGVER),g' \
		-e 's,[@]RELEASE[@],$(RELEASE),g' \
		$< > $@.tmp && mv $@.tmp $@
