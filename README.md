# cyclonedx-merge
Tool to merge cyclonedx files

## Merge rules

If `ids` from both objects are the same we consider that the objects are equal and keep the first object.

### Annotations

Merge is performed if BomRef is the same on both objects

### Compositions

Merge is performed if BomRef is the same on both objects

### ExternalReferences

Merge is performed if url and type are the same on both objects

### Properties

Merge is performed if name is the same on both objects

### Services

Merge is performed if BomRef is the same on both objects


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
