# cyclonedx-merge
![Coverage](https://img.shields.io/badge/Coverage-11.8%25-red)
Tool to merge cyclonedx files (json/xml)

## Install

```
go install github.com/fnxpt/cyclonedx-merge@latest
```

## Run with docker

```
docker run -v `pwd`/sbom/:/sbom/ fnxpt/cyclonedx-merge:latest --dir /sbom/ > output.json
```

## Usage
```
Usage:
  -bomref value
        BOMRef of the merged parent component
  -dir value
        merges files in directory
  -file value
        merges file
  -format value
        output format - json/xml (default: json)
  -group value
        group of the merged parent component
  -mode value
        merge mode - normal/flat/smart (default: normal)
  -name value
        name of the merged parent component
  -output value
        output file (default: stdout)
  -type value
        type of the aggregator component
  -version value
        version of the merged parent component
```

## Modes

|Mode|Description|Example|
|---|---|---|
|Normal|This merge the sboms and keep the relationships, this may lead to wrong dependencies on the graph|SBOM1 has libA that depends on libB:1.0;<br />SBOM2 has libA that depends on libB:2.0<br /><br />Outcome: Merged SBOM will have libA that dependes on libB:1.0 and libB:2.0|
|Flat|This is merge the sboms and sets all relationships on the second level, this leads to a simplified version of the graph, losing most of the relationships|SBOM1 has libA that depends on libB:1.0 that dependends on libC<br />SBOM2 has libA that depends on libB:2.0 that dependends on libC<br /><br />Outcome: Merged SBOM will have libA that depends on libB:1.0 and libC in main component of SBOM1 and libA that depends on libB:2.0 and libC in main component of SBOM2|
|Smart|To be implemented||

## Merge rules

If `ids` from both objects are the same we consider that the objects are equal and keep the first object.

|Type|Ids|Comment|
|---|---|---|
|Annotations|BomRef|   |
|Components|BomRef|   |
|Compositions|BomRef|   |
|ExternalReferences|URL & Type|   |
|Properties|Name & Value|If different files have properties with the same name, its impossible to merge them|
|Services|BomRef|   |
