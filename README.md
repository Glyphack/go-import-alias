# Go Import Alias

Go import alias is a linter to check consistency in import aliases

## Rules

consider import path: github.com/projectcontour/x/y/z/v\*

1. the alias name should be subset of `x[optional]_y[optional]_z[optional]_v*` where optional means it can be present or not.
2. one of `x` or `y` or `z` must be present in alias name.
3. If version exists in path, must be specified.
4. words like `apis` should be `api` in import alias

### Example

**Valid imports**

```go
import meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import api_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import api_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
```

**Invalid imports**

```go
import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import meta "k8s.io/apimachinery/pkg/apis/meta/v1"
import meta_api_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import api_meta "k8s.io/apimachinery/pkg/apis/meta/v1"
```
