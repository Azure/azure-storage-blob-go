# Change Log

> See [BreakingChanges](BreakingChanges.md) for a detailed list of API breaks.

## Version XX.XX.XX:
- Added the ability to obtain User Delegation Keys (UDK)
- Added the ability to create User Delegation SAS tokens from UDKs

## Version 0.3.0:
- Removed most panics from the library. Several functions now return an error.
- Removed 2016 and 2017 service versions.
- Added support for module.
- Fixed chunking bug in highlevel function uploadStream.