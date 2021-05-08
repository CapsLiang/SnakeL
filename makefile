APPS = logicserver rcenterserver roomserver
BASE = D:\PATH\myvendor
CAPSLIANG = D:\WorkSpace\GolandProjects\SnakeL

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

