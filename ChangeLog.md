# Change Log

> See [BreakingChanges](BreakingChanges.md) for a detailed list of API breaks.

## Version 0.10.0:
- Added support for CopyBlobFromURL (sync) and upgrade version to 2019-02-02.
- Provided default values for UploadStreamToBlockBlobOptions and refactored UploadStreamToBlockBlob.
- Added support for multiple start/expiry time formats.
- Added Solaris support.
- Enabled recovering from a unexpectedEOF error.

## Version 0.9.0:
- Updated go.mod to fix dependency issues.

## Version 0.8.0:
- Fixed error handling in high-level function DoBatchTransfer, and made it public for easy customization

## Version 0.7.0:
- Added the ability to obtain User Delegation Keys (UDK)
- Added the ability to create User Delegation SAS tokens from UDKs
- Added support for generating and using blob snapshot SAS tokens
- General secondary host improvements

## Version 0.3.0:
- Removed most panics from the library. Several functions now return an error.
- Removed 2016 and 2017 service versions.
- Added support for module.
- Fixed chunking bug in highlevel function uploadStream.