# cyclonedx-merge
Tool to merge cyclonedx files

## Assumptions

1. Whenever a bomref exists, we use it as an identifier, so if two seperate sbom's have the same id they will be merged
2. When merging we consider that the objects are the same and we use the first object

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
- [ ] Merge Annotations
- [ ] Merge Compositions
- [ ] Merge ExternalReferences
- [ ] Merge Properties
- [ ] Merge Services
- [ ] Make it generic
- [ ] Clean code
