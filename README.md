# OpenAPI 3.0.x spec external refs resolve tool

### About
Tool collects local external refs of your openapi specification into single file

### State
- [x] Resolve schemas defined in spec file
- [ ] Validate generated spec
- [ ] Resolve schemas defined in external files
- [ ] Resolve requestBodies
- [ ] Resolve responses
- [ ] Resolve securitySchemas

### Usage
```
  -f string
        path to spec file
  -o string
        path to out spec file (default "/dev/stdout")
```

# Example
Organize your files like:
```
any-api-root-dir/
├─ schemas/
│  ├─ external_schema1.yml
│  ├─ external_schema2.yml
├─ schema.openapi.yml
```

Use references in your spec:
```yaml
paths:
    responses:
        200:
            schema:
                $ref: schemas/external_schema1.yml#/ExternalSchema
```

And run
```shell
    spec-ref-resolve -f any-api-root-dir/schema.openapi.yml
```