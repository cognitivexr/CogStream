PIP_CMD ?= pip3
VENV_BIN ?= python3 -m venv
VENV_DIR ?= .venv
VENV_ACTIVATE = $(VENV_DIR)/bin/activate
VENV_RUN = . $(VENV_ACTIVATE)

$(VENV_ACTIVATE): requirements.txt
	test -d $(VENV_DIR) || $(VENV_BIN) $(VENV_DIR)
	$(VENV_RUN); $(PIP_CMD) install --upgrade pip setuptools wheel
	touch $(VENV_ACTIVATE)

venv: $(VENV_ACTIVATE)	## create a new (empty) venv

install: venv
	$(VENV_RUN); $(PIP_CMD) install -r requirements.txt

start:
	$(VENV_RUN); python -m debug.main


