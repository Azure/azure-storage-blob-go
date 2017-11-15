# Azure Storage Blob SDK for Go
The Microsoft Azure Storage SDK for Go allows you to build applications that takes advantage of Azure's scalable cloud storage. 

This repository contains the open source Blob SDK for Go. Documentation and examples can be found [here](https://godoc.org/github.com/Azure/azure-storage-blob-go/2016-05-31/azblob).

## Features
* Blob Storage
	* Create/Read/Update/Delete Block Blobs
	* Create/Read/Update/Delete Page Blobs
	* Create/Read/Update/Delete Append Blobs

## Getting Started
* If you don't already have it, install [the Go Programming runtime](https://golang.org/dl/)
* Go get the SDK:

```go get -u https://github.com/Azure/azure-storage-blob-go```
		
## SDK Architecture

* Storage SDK for Go provides 2 set of APIs: high-level, and low-level APIs
	* ServiceURI, ContainerURI and BlobURI objects provide the low-level API functionality and maps one-to-one to the [Azure Storage Blob REST APIs](https://docs.microsoft.com/en-us/rest/api/storageservices/blob-service-rest-api)
	* A set of high-level APIs are provided in highlevel.go file. These functions provide high level abstractions for convenience like uploading a large stream to Blob storage using multiple PutBlock requests

## Code Samples
* [Blob Storage Samples](https://github.com/seguler/azure-storage-blob-go/blob/master/2016-05-31/azblob/zt_examples_test.go)

## License
This project is licensed under MIT.

## Contributing
This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
