# storage-benchmark-go
 ## Why this benchmark?
 We need to test a specific read operation of large files from network-based storage solutions. This benchmark tool measures the sequential and parallel read performance of large files in a particular way.
 
 ## Usage

    Usage:
      storage-benchmark-go[options]

    Options:
      -dataChunkSize int
          Optianal data chunk size argument between 1024-524288(512KB) default is 10000 (default 10000)
      -filePath string 
          filePath for the read test. A file larger than 100GB is recommended. (default "256GB.dummy")
      -offsetOption int
          An optional offset number as number of chunks between 0-10000. Otherwise it is random (default -1)
      -parallelOption
          Optional parallel read flag.
      -h, --help 
          Show help and usage information

 ## Example Read Scenario on a 100GB File 
 ![Example Read Scenario on a 100GB File](images/read_behavior.jpg)


 ## Direct Download
 You can get the latest release from [here](https://github.com/muhsingurel/storage-benchmark-go/releases)


 ## How to compile?
 1) Install GO if you don't already have it https://golang.org/dl/
 2) Download the repo
 3) From your terminal set your working directory to the repo
 4) Type ```go build .\storage-benchmark-go.go``` to create an executable file or ```go run .\storage-benchmark-go.go``` to run it without building the code.
