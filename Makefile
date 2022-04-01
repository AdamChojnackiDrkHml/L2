
.PHONY: build buildD runD buildC runC cleanO clean cleanC cleanD

all: build

build: buildC buildD

buildD: 
	@go build -o build/decoder/decoder cmd/decoderMain/decoderMain.go

runD:
	@./build/decoder/decoder  $(IN) $(OUT)

buildC: 
	@go build -o build/coder/coder cmd/coderMain/coderMain.go

runC:
	@./build/coder/coder $(IN) $(OUT)

cleanO: 
	@rm data/output/*

cleanD:
	@rm build/decoder/* 

cleanC:
	@rm build/coder/*

clean: cleanC cleanD
