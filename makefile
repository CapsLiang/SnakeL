APPS = logicserver rcenterserver roomserver
BASE = /../SnakeLBase
CAPSLIANG = /root/workspace/go/SnakeL

BUILDVER = '0.1'
BUILDDATE = 'date +%F'

all: install

install: 
	export GOPATH=$(PWD):$(BASE):$(CAPSLIANG)\
	&& for ser in $(APPS);\
	do \
		go install -x $$ser;\
		if [ "$$?" != "0" ]; then\
			exit 1;\
		fi;\
	done

