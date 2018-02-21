PACKAGE  = 5sigma/spyder
BASE     = $(CURDIR)
WEB_PATH = $(BASE)/doc_template

all: clean web bin clean

bin:
	go install .

web:
	cd $(WEB_PATH); \
	npm run build; \
	mkdir dist; \
	rm build/static/js/*.map; \
	cp build/static/js/*.js dist/doc.js; \
	rm build/static/css/*.map; \
	cp build/static/css/*.css dist/doc.css; \
	go-bindata -o ../docgen/assets.go -pkg docgen -nocompress dist/
	
clean:
	rm -Rf $(WEB_PATH)/dist; \
	rm -Rf $(WEB_PATH)/build
