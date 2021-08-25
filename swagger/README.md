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
input-file: https://raw.githubusercontent.com/Azure/azure-rest-api-specs/storage-dataplane-preview/specification/storage/data-plane/Microsoft.BlobStorage/preview/2020-08-04/blob.json
go: true
output-folder: Go_BlobStorage
namespace: azblob
go-export-clients: false
enable-xml: true
file-prefix: zz_generated_
```

### Set blob metadata to the right type

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: >
      $["BlobMetadata"] = {
        "type": "object",
        "xml": {
          "name": "Metadata"
        },
        "additionalProperties": {
          "type": "string"
        }
      }
```

### Generate immutability policy correctly

```yaml
directive:
  - from: swagger-document
    where: $.parameters.ImmutabilityPolicyMode
    transform: >
      $.enum = [
        "Mutable",
        "Unlocked",
        "Locked"
      ]
```

### Add permissions to ListBlobsInclude

This *does* work, it's just not documented, or in the swagger unfortunately. Some directives jank it in.

These directives may need to be removed eventually, once these items are included.

For some reason, permissions gets added 3x, so a check for includes is added.

```yaml
directive:
  - from: swagger-document
    where: $.parameters.ListBlobsInclude
    transform: >
      if (!$.items.enum.includes("permissions"))
        $.items.enum.push("permissions")
```

```yaml
directive:
  - from: swagger-document
    where: $.definitions.BlobPropertiesInternal
    transform: >
      $.properties["Owner"] = {
        "type": "string"
      };
      $.properties["Group"] = {
        "type": "string"
      };
      $.properties["Permissions"] = {
        "type": "string"
      };
      $.properties["Acl"] = {
        "type": "string"
      };
```

### TODO: Get rid of StorageError since we define it
### TODO: rfc3339Format = "2006-01-02T15:04:05Z" //This was wrong in the generated code
