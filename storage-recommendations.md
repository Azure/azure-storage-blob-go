# Recommendations for the Azure Storage Blob Go panics

> Note: All comments are in regards to https://github.com/Azure/azure-storage-blob-go/tree/0.2.0 
> (commit: bb46532f68b79e9e1baca8fb19a382ef5d40ed33)

## Panic to protect from data in an unexpected format

> **Recommendation:** Propogate parsing failures as errors instead of panics.

These panics are at the core of our conversation, and really hit at the heart of the conversation. There are two mechanisms being discussed for communication failure: `panic` and `error`. You can see the two styles demonstrated in the Appendix. As discussed ad nauseam, there are advatages and disadvantages, but the community has long been adopting the error pattern. If adding errors to the getter signatures is too cumbersome, consider moving the Header parsing upstream so that the error would be propogated as part of initially handling the response instead discovering the misformatting lazily.

<details>

| Location | Panic Text |
|:-------------------------------------------|:-------------------------------------------------------------------------------|
| [2016-05-31/azblob/zz_generated_models.go:320](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L320) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:338](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L338) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:432](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L432) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:485](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L485) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:503](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L503) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:546](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L546) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:589](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L589) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:607](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L607) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:703](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L703) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:726](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L726) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:764](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L764) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:792](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L792) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:863](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L863) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:881](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L881) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:899](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L899) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:947](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L947) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:970](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L970) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1013](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1013) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1036](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1036) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1084](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1084) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1102](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1102) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1145](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1145) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1163](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1163) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1181](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1181) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1232](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1232) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1255](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1255) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1303](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1303) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1353](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1353) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1371](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1371) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1389](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1389) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1454](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1454) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1472](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1472) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1515](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1515) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1558](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1558) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1576](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1576) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1637](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1637) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1655](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1655) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1726](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1726) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1744](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1744) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1762](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1762) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1819](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1819) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1837](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1837) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1880](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1880) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:1898](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L1898) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2016](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2016) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2044](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2044) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2082](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2082) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2105](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2105) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2191](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2191) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2306](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2306) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2324](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2324) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2367](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2367) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2385](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2385) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2408](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2408) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2453](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2453) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2466](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2466) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2484](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2484) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2585](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2585) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2603](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2603) | panic(err) |
| [2016-05-31/azblob/zz_generated_models.go:2688](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_models.go#L2688) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:675](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L675) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:718](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L718) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:736](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L736) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:784](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L784) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:802](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L802) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:815](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L815) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:858](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L858) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:876](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L876) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:924](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L924) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:942](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L942) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:990](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L990) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1066](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1066) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1079](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1079) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1117](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1117) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1130](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1130) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1148](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1148) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1186](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1186) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1219](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1219) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1277](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1277) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1295](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1295) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1338](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1338) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1356](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1356) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1404](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1404) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1417](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1417) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1435](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1435) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1478](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1478) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1501](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1501) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1584](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1584) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1602](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1602) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1645](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1645) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1696](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1696) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1709](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1709) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1732](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1732) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1775](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1775) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1788](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1788) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1836](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1836) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1849](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1849) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1872](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1872) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1917](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1917) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1935](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1935) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:1953](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L1953) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2050](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2050) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2068](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2068) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2116](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2116) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2134](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2134) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2147](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2147) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2190](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2190) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2208](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2208) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2256](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2256) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2274](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2274) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2317](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2317) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2378](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2378) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2396](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2396) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2454](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2454) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2472](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2472) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2515](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2515) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2533](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2533) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2581](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2581) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2599](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2599) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2642](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2642) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2660](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2660) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2743](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2743) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2756](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2756) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2769](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2769) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2807](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2807) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2820](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2820) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2843](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2843) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2881](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2881) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:2904](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2904) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3003](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3003) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3061](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3061) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3166](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3166) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3179](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3179) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3192](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3192) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3210](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3210) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3263](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3263) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3281](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3281) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3324](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3324) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3337](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3337) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3360](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3360) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3403](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3403) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3416](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3416) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3434](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3434) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3477](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3477) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3490](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3490) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3508](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3508) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3551](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3551) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3564](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3564) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3577](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3577) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3600](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3600) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3645](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3645) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3658](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3658) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3676](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3676) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3777](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3777) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3795](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3795) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:3882](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L3882) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:426](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L426) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:439](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L439) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:452](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L452) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:470](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L470) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:513](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L513) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:526](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L526) | panic(err) |
| [2017-07-29/azblob/zz_generated_models.go:549](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L549) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:689](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L689) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:702](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L702) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:715](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L715) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:738](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L738) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:781](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L781) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:794](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L794) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:822](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L822) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:865](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L865) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:913](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L913) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:936](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L936) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:984](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L984) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1007](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1007) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1020](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1020) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1063](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1063) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1086](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1086) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1134](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1134) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1157](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1157) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1205](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1205) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1265](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1265) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1341](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1341) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1364](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1364) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1377](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1377) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1415](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1415) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1428](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1428) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1446](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1446) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1484](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1484) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1497](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1497) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1535](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1535) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1679](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1679) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1702](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1702) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1745](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1745) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1768](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1768) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1816](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1816) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1829](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1829) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1852](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1852) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1895](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1895) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:1923](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1923) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2011](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2011) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2034](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2034) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2077](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2077) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2133](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2133) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2146](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2146) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2174](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2174) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2217](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2217) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2230](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2230) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2283](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2283) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2296](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2296) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2349](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2349) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2362](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2362) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2390](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2390) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2435](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2435) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2453](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2453) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2476](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2476) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2534](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2534) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2557](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2557) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2605](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2605) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2628](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2628) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2641](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2641) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2684](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2684) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2707](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2707) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2755](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2755) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2778](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2778) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2821](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2821) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2874](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2874) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2945](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2945) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:2978](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L2978) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3079](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3079) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3102](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3102) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3145](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3145) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3168](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3168) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3216](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3216) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3239](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3239) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3282](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3282) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3305](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3305) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3388](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3388) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3401](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3401) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3414](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3414) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3452](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3452) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3465](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3465) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3488](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3488) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3526](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3526) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3554](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3554) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3653](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3653) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3716](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3716) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3831](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3831) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3844](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3844) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3857](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3857) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3880](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3880) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3933](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3933) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3956](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3956) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:3999](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3999) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4012](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4012) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4040](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4040) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4083](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4083) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4096](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4096) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4119](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4119) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4162](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4162) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4175](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4175) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4198](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4198) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4241](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4241) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4254](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4254) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4267](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4267) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4295](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4295) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4340](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4340) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4353](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4353) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4376](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4376) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4438](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4438) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4539](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4539) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4562](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4562) | panic(err) |
| [2018-03-28/azblob/zz_generated_models.go:4660](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L4660) | panic(err) |
| [2016-05-31/azblob/zz_response_helpers.go:57](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_response_helpers.go#L57) | panic(err) |
| [2016-05-31/azblob/zz_response_helpers.go:109](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_response_helpers.go#L109) | panic(err) |
| [2016-05-31/azblob/credential_shared_key.go:23](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/credential_shared_key.go#L23) | panic(err) |
| [2016-05-31/azblob/credential_shared_key.go:167](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/credential_shared_key.go#L167) | panic(err) |
| [2016-05-31/azblob/sas_service.go:40](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/sas_service.go#L40) | panic(err) |
| [2016-05-31/azblob/sas_service.go:48](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/sas_service.go#L48) | panic(err) |
| [2017-07-29/azblob/sas_service.go:40](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/sas_service.go#L40) | panic(err) |
| [2017-07-29/azblob/sas_service.go:48](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/sas_service.go#L48) | panic(err) |
| [2017-07-29/azblob/zc_credential_shared_key.go:23](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_credential_shared_key.go#L23) | panic(err) |
| [2017-07-29/azblob/zc_credential_shared_key.go:167](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_credential_shared_key.go#L167) | panic(err) |
| [2017-07-29/azblob/zc_sas_account.go:35](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_sas_account.go#L35) | panic(err) |
| [2018-03-28/azblob/zc_sas_account.go:35](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_sas_account.go#L35) | panic(err) |
| [2018-03-28/azblob/zc_credential_shared_key.go:167](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_credential_shared_key.go#L167) | panic(err) |
</details>

## Panics to wrap unexpected errors

> **Recommendation:** Errors that are important enough to check should be propogated upwards, errors that are not important enough to check need not take any action.

There are plenty of examples in my code where I don't bother to check an error, because the chance of it failing is so low, and the failure is so inconsequential. For example, I almost never check for errors when I call the [`fmt.Printf` function](https://godoc.org/fmt#Printf).

In cases like the UUID panics detailed below, failure is very rare. However, failure to read random bits totally impairs UUIDs ability to generate. In this scenario, informing the caller that this operation has failed is the correct course of action. However, simply propogating the error is currently the community's expected behavior.

<details>

| Location | Panic Text |
|:---------|:-----------|
| [2016-05-31/azblob/uuid.go:26](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/uuid.go#L26) | panic("ran.Read failed") |
| [2017-07-29/azblob/zc_uuid.go:26](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_uuid.go#L26) | panic("ran.Read failed") |
| [2018-03-28/azblob/zc_uuid.go:26](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_uuid.go#L26) | panic("ran.Read failed") |
| [2016-05-31/azblob/mmf_unix.go:25](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/mmf_unix.go#L25) | panic(err) |
| [2017-07-29/azblob/zc_mmf_windows.go:36](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_mmf_windows.go#L36) | panic(err) |
| [2017-07-29/azblob/zc_mmf_unix.go:25](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_mmf_unix.go#L25) | panic(err) |
| [2018-03-28/azblob/zc_mmf_windows.go:36](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_mmf_windows.go#L36) | panic(err) |
| [2018-03-28/azblob/zc_mmf_unix.go:25](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_mmf_unix.go#L25) | panic(err) |
| [2017-07-29/azblob/zc_policy_retry.go:179](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_policy_retry.go#L179) | panic(err) |
| [2016-05-31/azblob/policy_retry.go:169](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/policy_retry.go#L169) | panic(err) |
| [2018-03-28/azblob/zc_policy_retry.go:180](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L180) | panic(err) |
| [2016-05-31/azblob/mmf_windows.go:36](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/mmf_windows.go#L36) | panic(err) |
| [2017-07-29/azblob/zt_url_blob_test.go:204](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zt_url_blob_test.go#L204) | panic("ran.Read failed") |
| [2017-07-29/azblob/zc_util_validate.go:45](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_util_validate.go#L45) | panic("failed to seek stream") |
| [2018-03-28/azblob/sas_service.go:41](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/sas_service.go#L41) | panic(err) |
| [2018-03-28/azblob/sas_service.go:49](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/sas_service.go#L49) | panic(err) || [2018-03-28/azblob/zc_util_validate.go:45](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_util_validate.go#L45) | panic("failed to seek stream") |
| [2018-03-28/azblob/zc_util_validate.go:57](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_util_validate.go#L57) | panic(err) |
| [2017-07-29/azblob/zc_util_validate.go:57](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_util_validate.go#L57) | panic(err) |
| [2018-03-28/azblob/zt_url_blob_test.go:205](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zt_url_blob_test.go#L205) | panic("ran.Read failed") |
| [2018-03-28/azblob/zc_policy_retry.go:267](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L267) | panic("invalid state, response should not be nil when the operation is executed successfully") |

</details>

## Panics for bounds checking

> **Recommendation:** Return an error when paramters are out of bounds, instead of panicking.

When an invalid value is passed to a function, it is expected that it communicates back failure in some capacity. Some standard library functions error, and some panic. However, it would be highly unusual for a thrid-party library to panic.

<details>

| Location | Panic Text |
|:-------------------------------------------|:-------------------------------------------------------------------------------|
| [2016-05-31/azblob/sas_account.go:28](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/sas_account.go#L28) | panic("Account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType") |
| [2016-05-31/azblob/sas_account.go:35](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/sas_account.go#L35) | panic(err) |
| [2016-05-31/azblob/url_page_blob.go:31](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L31) | panic("p can't be nil") |
| [2016-05-31/azblob/url_page_blob.go:55](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L55) | panic("sequenceNumber must be greater than or equal to 0") |
| [2016-05-31/azblob/url_page_blob.go:105](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L105) | panic("Size must be a multiple of PageBlobPageBytes (512)") |
| [2016-05-31/azblob/url_page_blob.go:116](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L116) | panic("sequenceNumber must be greater than or equal to 0") |
| [2016-05-31/azblob/url_page_blob.go:144](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L144) | panic("PageRange's Start value must be greater than or equal to 0") |
| [2016-05-31/azblob/url_page_blob.go:147](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L147) | panic("PageRange's End value must be greater than 0") |
| [2016-05-31/azblob/url_page_blob.go:150](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L150) | panic("PageRange's Start value must be a multiple of 512") |
| [2016-05-31/azblob/url_page_blob.go:153](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L153) | panic("PageRange's End value must be 1 less than a multiple of 512") |
| [2016-05-31/azblob/url_page_blob.go:156](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L156) | panic("PageRange's End value must be after the start") |
| [2016-05-31/azblob/url_page_blob.go:190](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L190) | panic("Ifsequencenumberlessthan can't be less than -1") |
| [2016-05-31/azblob/url_page_blob.go:193](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L193) | panic("IfSequenceNumberLessThanOrEqual can't be less than -1") |
| [2016-05-31/azblob/url_page_blob.go:196](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_page_blob.go#L196) | panic("IfSequenceNumberEqual can't be less than -1") |
| [2016-05-31/azblob/policy_retry.go:60](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/policy_retry.go#L60) | panic("RetryPolicy must be RetryPolicyExponential or RetryPolicyFixed") |
| [2016-05-31/azblob/policy_retry.go:63](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/policy_retry.go#L63) | panic("MaxTries must be >= 0") |
| [2016-05-31/azblob/policy_retry.go:66](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/policy_retry.go#L66) | panic("TryTimeout, RetryDelay, and MaxRetryDelay must all be >= 0") |
| [2016-05-31/azblob/policy_retry.go:69](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/policy_retry.go#L69) | panic("RetryDelay must be <= MaxRetryDelay") |
| [2016-05-31/azblob/policy_retry.go:72](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/policy_retry.go#L72) | panic("Both RetryDelay and MaxRetryDelay must be 0 or neither can be 0") |
| [2017-07-29/azblob/atomicmorph.go:13](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/atomicmorph.go#L13) | panic("target and morpher mut not be nil") |
| [2017-07-29/azblob/atomicmorph.go:32](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/atomicmorph.go#L32) | panic("target and morpher mut not be nil") |
| [2017-07-29/azblob/atomicmorph.go:51](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/atomicmorph.go#L51) | panic("target and morpher mut not be nil") |
| [2017-07-29/azblob/atomicmorph.go:70](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/atomicmorph.go#L70) | panic("target and morpher mut not be nil") |
| [2017-07-29/azblob/zc_pipeline.go:25](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_pipeline.go#L25) | panic("c can't be nil") |
| [2017-07-29/azblob/zc_policy_retry.go:70](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_policy_retry.go#L70) | panic("RetryPolicy must be RetryPolicyExponential or RetryPolicyFixed") |
| [2017-07-29/azblob/zc_policy_retry.go:73](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_policy_retry.go#L73) | panic("MaxTries must be >= 0") |
| [2017-07-29/azblob/zc_policy_retry.go:76](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_policy_retry.go#L76) | panic("TryTimeout, RetryDelay, and MaxRetryDelay must all be >= 0") |
| [2017-07-29/azblob/zc_policy_retry.go:79](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_policy_retry.go#L79) | panic("RetryDelay must be <= MaxRetryDelay") |
| [2017-07-29/azblob/zc_policy_retry.go:82](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_policy_retry.go#L82) | panic("Both RetryDelay and MaxRetryDelay must be 0 or neither can be 0") |
| [2018-03-28/azblob/url_page_blob.go:30](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L30) | panic("p can't be nil") |
| [2018-03-28/azblob/url_page_blob.go:54](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L54) | panic("sequenceNumber must be greater than or equal to 0") |
| [2018-03-28/azblob/url_page_blob.go:115](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L115) | panic("Size must be a multiple of PageBlobPageBytes (512)") |
| [2018-03-28/azblob/url_page_blob.go:126](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L126) | panic("sequenceNumber must be greater than or equal to 0") |
| [2018-03-28/azblob/url_page_blob.go:154](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L154) | panic("PageRange's Start value must be greater than or equal to 0") |
| [2018-03-28/azblob/url_page_blob.go:157](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L157) | panic("PageRange's End value must be greater than 0") |
| [2018-03-28/azblob/url_page_blob.go:160](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L160) | panic("PageRange's Start value must be a multiple of 512") |
| [2018-03-28/azblob/url_page_blob.go:163](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L163) | panic("PageRange's End value must be 1 less than a multiple of 512") |
| [2018-03-28/azblob/url_page_blob.go:166](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L166) | panic("PageRange's End value must be after the start") |
| [2018-03-28/azblob/url_page_blob.go:206](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L206) | panic("Ifsequencenumberlessthan can't be less than -1") |
| [2018-03-28/azblob/url_page_blob.go:209](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L209) | panic("IfSequenceNumberLessThanOrEqual can't be less than -1") |
| [2018-03-28/azblob/url_page_blob.go:212](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_page_blob.go#L212) | panic("IfSequenceNumberEqual can't be less than -1") |
| [2018-03-28/azblob/url_append_blob.go:94](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_append_blob.go#L94) | panic("IfAppendPositionEqual can't be less than -1") |
| [2018-03-28/azblob/url_append_blob.go:97](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_append_blob.go#L97) | panic("IfMaxSizeLessThanOrEqual can't be less than -1") |
| [2018-03-28/azblob/url_container.go:268](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_container.go#L268) | panic("MaxResults must be >= 0") |
| [2018-03-28/azblob/highlevel.go:69](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L69) | panic(fmt.Sprintf("BlockSize option must be > 0 and <= %d", BlockBlobMaxUploadBlobBytes)) |
| [2018-03-28/azblob/highlevel.go:101](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L101) | panic(fmt.Sprintf("The buffer's size is too big or the BlockSize is too small; the number of blocks must be <= %d", BlockBlobMaxBlocks)) |
| [2018-03-28/azblob/highlevel.go:192](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L192) | panic("BlockSize option must be >= 0") |
| [2018-03-28/azblob/highlevel.go:199](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L199) | panic("offset option must be >= 0") |
| [2018-03-28/azblob/highlevel.go:203](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L203) | panic("count option must be >= 0") |
| [2018-03-28/azblob/highlevel.go:220](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L220) | panic(fmt.Errorf("the buffer's size should be equal to or larger than the request count of bytes: %d", count)) |
| [2018-03-28/azblob/highlevel.go:273](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/highlevel.go#L273) | panic("file must not be nil") |
| [2018-03-28/azblob/url_block_blob.go:34](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_block_blob.go#L34) | panic("p can't be nil") |
| [2018-03-28/azblob/url_blob.go:18](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_blob.go#L18) | panic("p can't be nil") |
| [2018-03-28/azblob/url_blob.go:38](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_blob.go#L38) | panic("p can't be nil") |
| [2018-03-28/azblob/zc_retry_reader.go:59](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_retry_reader.go#L59) | panic("getter must not be nil") |
| [2018-03-28/azblob/zc_retry_reader.go:62](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_retry_reader.go#L62) | panic("info.Count must be >= 0") |
| [2018-03-28/azblob/zc_retry_reader.go:65](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_retry_reader.go#L65) | panic("o.MaxRetryRequests must be >= 0") |
| [2018-03-28/azblob/url_service.go:27](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_service.go#L27) | panic("p can't be nil") |
| [2018-03-28/azblob/url_service.go:103](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/url_service.go#L103) | panic("MaxResults must be >= 0") |
| [2018-03-28/azblob/zc_sas_account.go:28](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_sas_account.go#L28) | panic("Account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType") |
| [2017-07-29/azblob/zc_util_validate.go:59](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_util_validate.go#L59) | panic(errors.New("stream must be set to position 0")) |
| [2017-07-29/azblob/url_page_blob.go:30](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L30) | panic("p can't be nil") |
| [2017-07-29/azblob/url_page_blob.go:54](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L54) | panic("sequenceNumber must be greater than or equal to 0") |
| [2017-07-29/azblob/url_page_blob.go:115](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L115) | panic("Size must be a multiple of PageBlobPageBytes (512)") |
| [2017-07-29/azblob/url_page_blob.go:126](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L126) | panic("sequenceNumber must be greater than or equal to 0") |
| [2017-07-29/azblob/url_page_blob.go:154](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L154) | panic("PageRange's Start value must be greater than or equal to 0") |
| [2017-07-29/azblob/url_page_blob.go:157](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L157) | panic("PageRange's End value must be greater than 0") |
| [2017-07-29/azblob/url_page_blob.go:160](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L160) | panic("PageRange's Start value must be a multiple of 512") |
| [2017-07-29/azblob/url_page_blob.go:163](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L163) | panic("PageRange's End value must be 1 less than a multiple of 512") |
| [2017-07-29/azblob/url_page_blob.go:166](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L166) | panic("PageRange's End value must be after the start") |
| [2017-07-29/azblob/url_page_blob.go:200](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L200) | panic("Ifsequencenumberlessthan can't be less than -1") |
| [2017-07-29/azblob/url_page_blob.go:203](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L203) | panic("IfSequenceNumberLessThanOrEqual can't be less than -1") |
| [2017-07-29/azblob/url_page_blob.go:206](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_page_blob.go#L206) | panic("IfSequenceNumberEqual can't be less than -1") |
| [2017-07-29/azblob/url_container.go:269](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_container.go#L269) | panic("MaxResults must be >= 0") |
| [2017-07-29/azblob/highlevel.go:66](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L66) | panic(fmt.Sprintf("BlockSize option must be > 0 and <= %d", BlockBlobMaxUploadBlobBytes)) |
| [2017-07-29/azblob/highlevel.go:84](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L84) | panic(fmt.Sprintf("The buffer's size is too big or the BlockSize is too small; the number of blocks must be <= %d", BlockBlobMaxBlocks)) |
| [2017-07-29/azblob/highlevel.go:176](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L176) | panic("BlockSize option must be >= 0") |
| [2017-07-29/azblob/highlevel.go:183](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L183) | panic("offset option must be >= 0") |
| [2017-07-29/azblob/highlevel.go:187](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L187) | panic("count option must be >= 0") |
| [2017-07-29/azblob/highlevel.go:204](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L204) | panic(fmt.Errorf("the buffer's size should be equal to or larger than the request count of bytes: %d", count)) |
| [2017-07-29/azblob/highlevel.go:257](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/highlevel.go#L257) | panic("file must not be nil") |
| [2017-07-29/azblob/url_block_blob.go:33](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_block_blob.go#L33) | panic("p can't be nil") |
| [2017-07-29/azblob/url_blob.go:17](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_blob.go#L17) | panic("p can't be nil") |
| [2017-07-29/azblob/url_blob.go:37](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_blob.go#L37) | panic("p can't be nil") |
| [2017-07-29/azblob/zc_retry_reader.go:59](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_retry_reader.go#L59) | panic("getter must not be nil") |
| [2017-07-29/azblob/zc_retry_reader.go:62](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_retry_reader.go#L62) | panic("info.Count must be >= 0") |
| [2017-07-29/azblob/zc_retry_reader.go:65](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_retry_reader.go#L65) | panic("o.MaxRetryRequests must be >= 0") |
| [2017-07-29/azblob/url_service.go:27](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_service.go#L27) | panic("p can't be nil") |
| [2017-07-29/azblob/url_service.go:103](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_service.go#L103) | panic("MaxResults must be >= 0") |
| [2017-07-29/azblob/zc_sas_account.go:28](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_sas_account.go#L28) | panic("Account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType") |
| [2016-05-31/azblob/url_container.go:21](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_container.go#L21) | panic("p can't be nil") |
| [2016-05-31/azblob/url_container.go:93](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_container.go#L93) | panic("the IfMatch and IfNoneMatch access conditions must have their default values because they are ignored by the service") |
| [2016-05-31/azblob/url_container.go:112](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_container.go#L112) | panic("the IfUnmodifiedSince, IfMatch, and IfNoneMatch must have their default values because they are ignored by the blob service") |
| [2016-05-31/azblob/url_container.go:184](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_container.go#L184) | panic("the IfMatch and IfNoneMatch access conditions must have their default values because they are ignored by the service") |
| [2016-05-31/azblob/url_container.go:258](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_container.go#L258) | panic("MaxResults must be >= 0") |
| [2016-05-31/azblob/credential_token.go:34](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/credential_token.go#L34) | panic("Token credentials require a URL using the https protocol scheme.") |
| [2016-05-31/azblob/highlevel.go:66](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/highlevel.go#L66) | panic(fmt.Sprintf("BlockSize option must be > 0 and <= %d", BlockBlobMaxPutBlockBytes)) |
| [2016-05-31/azblob/highlevel.go:89](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/highlevel.go#L89) | panic(fmt.Sprintf("The streamSize is too big or the BlockSize is too small; the number of blocks must be <= %d", BlockBlobMaxBlocks)) |
| [2016-05-31/azblob/highlevel.go:203](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/highlevel.go#L203) | panic("getBlob must not be nil") |
| [2016-05-31/azblob/url_block_blob.go:32](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_block_blob.go#L32) | panic("p can't be nil") |
| [2016-05-31/azblob/url_blob.go:21](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_blob.go#L21) | panic("p can't be nil") |
| [2016-05-31/azblob/url_blob.go:41](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_blob.go#L41) | panic("p can't be nil") |
| [2016-05-31/azblob/url_blob.go:100](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_blob.go#L100) | panic("The blob's range Offset must be >= 0") |
| [2016-05-31/azblob/url_blob.go:103](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_blob.go#L103) | panic("The blob's range Count must be >= 0") |
| [2018-03-28/azblob/zc_util_validate.go:59](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_util_validate.go#L59) | panic(errors.New("stream must be set to position 0")) |
| [2018-03-28/azblob/zc_policy_retry.go:71](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L71) | panic("RetryPolicy must be RetryPolicyExponential or RetryPolicyFixed") |
| [2018-03-28/azblob/zc_policy_retry.go:74](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L74) | panic("MaxTries must be >= 0") |
| [2018-03-28/azblob/zc_policy_retry.go:77](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L77) | panic("TryTimeout, RetryDelay, and MaxRetryDelay must all be >= 0") |
| [2018-03-28/azblob/zc_policy_retry.go:80](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L80) | panic("RetryDelay must be <= MaxRetryDelay") |
| [2018-03-28/azblob/zc_policy_retry.go:83](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_policy_retry.go#L83) | panic("Both RetryDelay and MaxRetryDelay must be 0 or neither can be 0") |
| [2016-05-31/azblob/url_service.go:34](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_service.go#L34) | panic("c can't be nil") |
| [2016-05-31/azblob/url_service.go:65](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_service.go#L65) | panic("p can't be nil") |
| [2016-05-31/azblob/url_service.go:150](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_service.go#L150) | panic("MaxResults must be >= 0") |
| [2016-05-31/azblob/url_append_blob.go:89](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_append_blob.go#L89) | panic("IfAppendPositionEqual can't be less than -1") |
| [2016-05-31/azblob/url_append_blob.go:92](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/url_append_blob.go#L92) | panic("IfMaxSizeLessThanOrEqual can't be less than -1") |
| [2018-03-28/azblob/zc_util_validate.go:23](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_util_validate.go#L23) | panic("The range offset must be >= 0") |
| [2018-03-28/azblob/zc_util_validate.go:26](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_util_validate.go#L26) | panic("The range count must be >= 0") |
| [2017-07-29/azblob/zc_util_validate.go:23](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_util_validate.go#L23) | panic("The range offset must be >= 0") |
| [2017-07-29/azblob/zc_util_validate.go:26](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_util_validate.go#L26) | panic("The range count must be >= 0") || [2018-03-28/azblob/atomicmorph.go:15](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/atomicmorph.go#L15) | panic(targetAndMorpherMustNotBeNil) |
| [2018-03-28/azblob/atomicmorph.go:34](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/atomicmorph.go#L34) | panic(targetAndMorpherMustNotBeNil) |
| [2018-03-28/azblob/atomicmorph.go:53](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/atomicmorph.go#L53) | panic(targetAndMorpherMustNotBeNil) |
| [2018-03-28/azblob/atomicmorph.go:72](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/atomicmorph.go#L72) | panic(targetAndMorpherMustNotBeNil) |
| [2018-03-28/azblob/zc_pipeline.go:25](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_pipeline.go#L25) | panic("c can't be nil") |
| [2017-07-29/azblob/url_append_blob.go:88](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_append_blob.go#L88) | panic("IfAppendPositionEqual can't be less than -1") |
| [2017-07-29/azblob/url_append_blob.go:91](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_append_blob.go#L91) | panic("IfMaxSizeLessThanOrEqual can't be less than -1") |
| [2017-07-29/azblob/zz_generated_models.go:80](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L80) | panic("s wasn't a slice or array") |
| [2018-03-28/azblob/zz_generated_models.go:80](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L80) | panic("s wasn't a slice or array") |
| [2017-07-29/azblob/sas_service.go:32](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/sas_service.go#L32) | panic("sharedKeyCredential can't be nil") |
| [2016-05-31/azblob/sas_service.go:32](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/sas_service.go#L32) | panic("sharedKeyCredential can't be nil") |
| [2017-07-29/azblob/url_container.go:20](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_container.go#L20) | panic("p can't be nil") |
| [2017-07-29/azblob/url_container.go:92](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_container.go#L92) | panic("the IfMatch and IfNoneMatch access conditions must have their default values because they are ignored by the service") |
| [2017-07-29/azblob/url_container.go:112](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_container.go#L112) | panic("the IfUnmodifiedSince, IfMatch, and IfNoneMatch must have their default values because they are ignored by the blob service") |
| [2017-07-29/azblob/url_container.go:184](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_container.go#L184) | panic("the IfMatch and IfNoneMatch access conditions must have their default values because they are ignored by the service") |
| [2017-07-29/azblob/url_container.go:244](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/url_container.go#L244) | panic("snapshots are not supported in this listing operation") |
| [2017-07-29/azblob/zc_credential_token.go:120](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zc_credential_token.go#L120) | panic("Token credentials require a URL using the https protocol scheme.") |
| [2018-03-28/azblob/sas_service.go:33](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/sas_service.go#L33) | panic("sharedKeyCredential can't be nil") |
| [2018-03-28/azblob/zt_test.go:40](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zt_test.go#L40) | panic(errors.New("RetryOption's RetryReadsFromSecondaryHost field must exist in the Blob SDK - uncomment it and make sure the field is returned from the retryReadsFromSecondaryHost() method too!")) |
| [2018-03-28/azblob/zc_credential_token.go:130](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zc_credential_token.go#L130) | panic("Token credentials require a URL using the https protocol scheme.") |

</details>

## Panics to protect from generation errors

> **Recommendation:** Run these checks as tests which fail at generation time, instead of shipping this logic to customers.

Detecting generation errors is important, and ensuring that data doesn't get further corrupted is important. However, it is unlikely that the shape of these types gets mutated. For that reason, only running these as tests at generation time would reduce the amount of concern customers had when inspecting our code base, while still protecting them from code-generation errors.

As an alternative, at least moving these panics into an `init` phase for the package would gaurentee that if a panic was going to occur, it would occur immediately instead of at random times during application lifecycle.

<details>

| Location | Panic Text |
|:-------------------------------------------|:-------------------------------------------------------------------------------|
| [2016-05-31/azblob/zz_generated_marshalling.go:62](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L62) | panic("size mismatch between AccessPolicy and accessPolicy") |
| [2016-05-31/azblob/zz_generated_marshalling.go:71](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L71) | panic("size mismatch between AccessPolicy and accessPolicy") |
| [2016-05-31/azblob/zz_generated_marshalling.go:110](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L110) | panic("size mismatch between BlobProperties and blobProperties") |
| [2016-05-31/azblob/zz_generated_marshalling.go:119](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L119) | panic("size mismatch between BlobProperties and blobProperties") |
| [2016-05-31/azblob/zz_generated_marshalling.go:136](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L136) | panic("size mismatch between Blob and blob") |
| [2016-05-31/azblob/zz_generated_marshalling.go:145](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L145) | panic("size mismatch between Blob and blob") |
| [2016-05-31/azblob/zz_generated_marshalling.go:164](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L164) | panic("size mismatch between ContainerProperties and containerProperties") |
| [2016-05-31/azblob/zz_generated_marshalling.go:173](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L173) | panic("size mismatch between ContainerProperties and containerProperties") |
| [2016-05-31/azblob/zz_generated_marshalling.go:188](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L188) | panic("size mismatch between GeoReplication and geoReplication") |
| [2016-05-31/azblob/zz_generated_marshalling.go:197](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zz_generated_marshalling.go#L197) | panic("size mismatch between GeoReplication and geoReplication") |
| [2017-07-29/azblob/zz_generated_models.go:2945](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2945) | panic("size mismatch between GeoReplication and geoReplication") |
| [2017-07-29/azblob/zz_generated_models.go:2954](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2954) | panic("size mismatch between GeoReplication and geoReplication") |
| [2017-07-29/azblob/zz_generated_models.go:632](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L632) | panic("size mismatch between BlobProperties and blobProperties") |
| [2017-07-29/azblob/zz_generated_models.go:641](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L641) | panic("size mismatch between BlobProperties and blobProperties") |
| [2017-07-29/azblob/zz_generated_models.go:2007](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2007) | panic("size mismatch between ContainerProperties and containerProperties") |
| [2017-07-29/azblob/zz_generated_models.go:2016](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L2016) | panic("size mismatch between ContainerProperties and containerProperties") |
| [2018-03-28/azblob/zz_generated_models.go:641](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L641) | panic("size mismatch between AccessPolicy and accessPolicy") |
| [2018-03-28/azblob/zz_generated_models.go:650](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L650) | panic("size mismatch between AccessPolicy and accessPolicy") |
| [2018-03-28/azblob/zz_generated_models.go:1636](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1636) | panic("size mismatch between BlobProperties and blobProperties") |
| [2018-03-28/azblob/zz_generated_models.go:1645](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L1645) | panic("size mismatch between BlobProperties and blobProperties") |
| [2018-03-28/azblob/zz_generated_models.go:3036](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3036) | panic("size mismatch between ContainerProperties and containerProperties") |
| [2018-03-28/azblob/zz_generated_models.go:3045](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3045) | panic("size mismatch between ContainerProperties and containerProperties") |
| [2018-03-28/azblob/zz_generated_models.go:3595](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3595) | panic("size mismatch between GeoReplication and geoReplication") |
| [2018-03-28/azblob/zz_generated_models.go:3604](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zz_generated_models.go#L3604) | panic("size mismatch between GeoReplication and geoReplication") |
| [2017-07-29/azblob/zz_generated_models.go:378](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L378) | panic("size mismatch between AccessPolicy and accessPolicy") |
| [2017-07-29/azblob/zz_generated_models.go:387](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zz_generated_models.go#L387) | panic("size mismatch between AccessPolicy and accessPolicy") |
</details>

## Panics in Tests

> **Recommendation:** Instead of panicking during tests, use `t.Fail()`, `t.FailNow()`, or even `t.Skip()` to allow other tests to execute and complete.

This is low priority, but still a source of curiosity. In instances where tests are panicking (and not being caught by a `recover`), they are actually bypassing many more rich options for communicating test failure. `t.Fail()` will allow for
the other tests in the suite to continue, potentially allowing for 

<details>

| Location | Panic Text |
|:-------------------------------------------|:-------------------------------------------------------------------------------|
| [2016-05-31/azblob/zt_url_blob_test.go:217](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zt_url_blob_test.go#L217) | panic("ran.Read failed") |
| [2016-05-31/azblob/zt_test.go:210](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2016-05-31/azblob/zt_test.go#L210) | panic("ACCOUNT_NAME and ACCOUNT_KEY environment vars must be set before running tests") |
| [2017-07-29/azblob/zt_test.go:39](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2017-07-29/azblob/zt_test.go#L39) | panic(errors.New("RetryOption's RetryReadsFromSecondaryHost field must exist in the Blob SDK - uncomment it and make sure the field is returned from the retryReadsFromSecondaryHost() method too!")) |
| [2018-03-28/azblob/zt_url_blob_test.go:205](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zt_url_blob_test.go#L205) | panic("ran.Read failed") |
| [2018-03-28/azblob/zt_test.go:40](https://github.com/Azure/azure-storage-blob-go/blob/0.2.0/2018-03-28/azblob/zt_test.go#L40) | panic(errors.New("RetryOption's RetryReadsFromSecondaryHost field must exist in the Blob SDK - uncomment it and make sure the field is returned from the retryReadsFromSecondaryHost() method too!")) |

</details>

# Appendix

## Demonstration of Panic and Error patterns

*panic:*
``` Go
// LeaseTime returns the value for header x-ms-lease-time.
func (blr BlobsLeaseResponse) LeaseTime() int32 {
	s := blr.rawResponse.Header.Get("x-ms-lease-time")
	if s == "" {
		return -1
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(i)
}
```

*error:*
``` Go
// ErrorEmptyHeader is returned when no value of a header was able to be parsed.
struct ErrorEmptyHeader string

func (err ErrorEmptyHeader) Error() string {
    return fmt.Sprintf("value of header %q was empty", string(err))
}

// LeaseTime returns the value for header x-ms-lease-time.
func (blr BlobsLeaseResponse) LeaseTime() (int32, error) {
    const leaseTime = "x-ms-lease-time"
    s := blr.rawResponse.Header.Get(leaseTime)
    if s == "" {
        return -1, ErrorEmptyHeader(leaseTime)
    }
    
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
        return -1, err
    }

    return int32(i), nil
}
```

Both effectively communicate that the contents of the header "x-ms-lease-time" are unable to be recognized as an 
integer. Between the two, it is also clear that declaring a panic is more concise. However, when considering the code
consuming these two strategies, the advantages and disadvantages are less clear. (Author's note: I firmly believe that
the subjectivity of deciding the merits of the two options below are what is driving the disagreement here.)

*error:*
``` Go
func WaitLeaseTime(ctx context.Context, blobResp azblob.BlobsLeaseResponse) error {
    numSeconds, err := blobResp.LeaseTime()
    if err != nil {
        return err
    }

    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(numSeconds * time.Second):
        return nil
    }
}
```

*panic:*
``` Go
func WaitLeaseTime(ctx context.Context, blobResp azblob.BlobLeaseResponse) error {
    defer func(){
        if r := recover(); r != nil {
            return r
        }
    }()

    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(blobResp.LeaseTime() * time.Second)
        return nil
    }
}
```

The key differences to notice, is the addition of two more variables declared in local scope when consuming with
errors. Besides needing to name these values, it also means that we cannot yse the value of `blobResp.LeaseTime()`
inline to compose a function call. However, the error check is done before the `select` clause is entered. Each exit
condition for this function is extremely obvious.

With the panic function, if we are to recover from the panic and propogate it up, we must add a `defer` which will
execute once a `return` statement is called, or a `panic` has been induced by code we consume. We then mutate the
return value that was given to the return statement to match the error value handed to us if we were terminating
because of a panic. Mutating return values isn't super complicated, but it can be unreadable, as seen here:
https://play.golang.org/p/NuNtqFcO01

