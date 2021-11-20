build:
	go build

clean:
	-rm jsonl-graph
	-rm testdata/*.svg

docs: clean build
	cat testdata/gin.jsonl | ./jsonl-graph > testdata/gin_nocolor.dot
	cat testdata/small.jsonl | ./jsonl-graph > testdata/small.dot

.PHONY: build clean docs 
