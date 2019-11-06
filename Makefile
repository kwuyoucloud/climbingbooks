# Go parameters
VERSION=1.12.9
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOFILE=./spider
MAINGOFILE=cmd/main.go
MAINGO=main

all: build

run:
	@$(GORUN) $(MAINGOFILE)
build:
	@$(GOBUILD) $(MAINGOFILE)
	@mv $(MAINGO) $(GOFILE)
	@echo "Build correct."

clean:
	@rm $(GOFILE)

gitpush:
	# will get comment from args
	@git add *
	@git commit -m 'args'
	@git push origin master
