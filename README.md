# cyclonedx-merge
Tool to merge cyclonedx files

## Merge rules

If `ids` from both objects are the same we consider that the objects are equal and keep the first object.

|Type|Ids|Comment|
|---|---|---|
|Annotations|BomRef|   |
|Compositions|BomRef|   |
|ExternalReferences|URL & Type|   |
|Properties|Name & Value|If different files have properties with the same name, its impossible to merge them|
|Services|BomRef|   |

## Usage
```
Usage:
  -dir value
    	merges files in directory
  -file value
    	merges file
  -format value
    	output format - json/xml (default: json)
  -nested
    	nested merge
  -output value
    	output file (default: stdout)
```

## Run

```
docker run -v `pwd`/sbom/:/sbom/ fnxpt/cyclonedx-merge:latest --dir /sbom/ > output.json
```

## TODO:

- [ ] Add tests
- [x] Merge Annotations
- [x] Merge Compositions
- [x] Merge ExternalReferences
- [x] Merge Properties
- [x] Merge Services
- [ ] Make it generic
- [ ] Clean code
