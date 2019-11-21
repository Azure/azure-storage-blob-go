# Azure Blob Storage for Golang

> see https://aka.ms/autorest

### Generation
```bash
cd swagger
autorest README.md --use=@microsoft.azure/autorest.go@v3.0.63
gofmt -w Go_BlobStorage/*
```

### Settings
``` yaml
input-file: ./blob.json
go: true
output-folder: Go_BlobStorage
namespace: azblob
go-export-clients: false
enable-xml: true
file-prefix: zz_generated_
```

### TODO: Get rid of StorageError since we define it