# gocommon

Common stuff for all my Go projects. Internally, it uses [miso](https://github.com/CurtisNewbie/miso).

## Configuration

GoAuth Resource/Path Report Configuration:

```yaml
goauth:
  client:
    enabled: true
  resource:
    - name: "Manage files"
      code: "manage-files"
      path:
        - url: "/open/api/file/upload/duplication/preflight"
          desc: "Preflight check for duplicate file uploads"
        - url: "/open/api/file/parent"
          desc: "User fetch parent file info"
        - url: "/open/api/file/move-to-dir"
          desc: "Upload file"
    - name: "Admin file service"
      code: "admin-file-service"
  path:
    - url: "/open/api/file/info"
      desc: "Fetch file info"
```