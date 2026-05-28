# controller-tools의 K8s 생태계 사용처

`controller-tools`는 Kubernetes API/Controller 프로젝트에서 쓰는 빌드 타임 생성 도구다.
클러스터 안에서 실행되는 컨트롤러나 admission webhook 서버가 아니라, 개발자가 로컬이나 CI에서
`controller-gen` CLI를 실행할 때 동작한다.

핵심 역할은 Go 코드와 `// +kubebuilder:...` marker 주석을 읽어서 Kubernetes가 사용할 수 있는
YAML manifest와 반복적인 Go 코드를 생성하는 것이다.

```text
Go API type + marker comments
        |
        v
controller-gen
        |
        v
CRD / RBAC / Webhook YAML + DeepCopy Go code
        |
        v
kubectl / kustomize / OLM / GitOps로 클러스터에 배포
```

## 왜 필요한가

Kubernetes extension 프로젝트에서는 보통 다음 생성물이 필요하다.

- `CustomResourceDefinition` YAML
- CRD OpenAPI v3 schema
- controller가 사용할 `Role` 또는 `ClusterRole` YAML
- admission webhook configuration YAML
- Kubernetes runtime object 구현에 필요한 `DeepCopy` Go 코드

이 파일들을 사람이 직접 계속 수정하면 Go type과 YAML schema가 쉽게 어긋난다. 그래서 Go API type을
source of truth로 두고, `controller-gen`으로 필요한 생성물을 다시 만든다.

## 대표 사용처

### Kubebuilder 프로젝트

Kubebuilder 프로젝트에서 가장 흔한 사용 방식은 `make generate`와 `make manifests`다.

```bash
make generate
make manifests
```

일반적으로 `make generate`는 `controller-gen object`를 호출해서 DeepCopy 코드를 만들고,
`make manifests`는 `controller-gen crd rbac webhook`을 호출해서 Kubernetes YAML을 만든다.

결과물은 보통 다음 위치에 생긴다.

```text
api/v1/zz_generated.deepcopy.go
config/crd/bases/*.yaml
config/rbac/role.yaml
config/webhook/*.yaml
```

### Operator SDK / OpenShift Operator 개발

Operator SDK나 OpenShift Operator 개발 흐름에서도 같은 패턴을 사용한다. 먼저 API type을 만든 뒤,
`*_types.go` 파일에 `Spec`, `Status`, validation marker를 작성한다.

```go
type MemcachedSpec struct {
	// +kubebuilder:validation:Minimum=0
	Size int32 `json:"size"`
}
```

그 다음 생성 명령을 실행한다.

```bash
make generate
make manifests
```

이 흐름에서 `controller-gen`은 `zz_generated.deepcopy.go`와 CRD manifest를 갱신한다.
생성된 manifest는 이후 `kustomize`, Operator bundle, OLM, GitOps 파이프라인에서 사용된다.

### controller-runtime 기반 커스텀 컨트롤러

Kubebuilder나 Operator SDK scaffold를 쓰지 않고 직접 `controller-runtime` 기반 컨트롤러를 만들 때도
`controller-gen`만 별도로 사용할 수 있다.

이 경우에도 패턴은 같다.

```bash
controller-gen object paths="./..."
controller-gen crd paths=./api/... output:crd:dir=config/crd/bases
controller-gen rbac:roleName=manager-role paths=./controllers/... output:rbac:dir=config/rbac
```

즉, `controller-tools`는 특정 scaffold 도구에만 묶인 도구가 아니라 CRD 기반 API 프로젝트에서
Go type과 Kubernetes manifest를 맞춰 주는 생성 도구로 사용할 수 있다.

## CLI 사용 형태

`controller-gen` 명령은 세 가지를 조합한다.

```text
controller-gen <generator들> paths=<읽을 Go package> output:<출력 방식>
```

예를 들어 CRD만 생성한다.

```bash
controller-gen crd paths=./api/... output:crd:dir=config/crd/bases
```

RBAC만 생성한다.

```bash
controller-gen rbac:roleName=manager-role paths=./controllers/... output:rbac:dir=config/rbac
```

DeepCopy 코드를 생성한다.

```bash
controller-gen object paths="./..."
```

CRD, RBAC, webhook을 한 번에 생성한다.

```bash
controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

각 부분의 의미는 다음과 같다.

- `crd`, `rbac`, `object`, `webhook`: 어떤 generator를 실행할지 정한다.
- `paths=...`: 어떤 Go package들을 읽을지 정한다.
- `output:...`: 생성 결과를 어디에 쓸지 정한다.
- `rbac:roleName=...`: 특정 generator에 넘기는 옵션이다.

## 생성물별 역할

### CRD YAML

`controller-gen crd`는 Go API type과 validation marker를 읽어서 CRD YAML을 만든다.

```bash
controller-gen crd paths=./api/... output:crd:dir=config/crd/bases
```

예상 결과:

```text
config/crd/bases/example.my.domain_widgets.yaml
```

이 파일의 OpenAPI v3 schema는 Kubernetes API server가 Custom Resource를 생성하거나 수정할 때
입력값 검증에 사용한다.

### RBAC YAML

`controller-gen rbac`은 controller 코드 근처의 RBAC marker를 읽어서 권한 manifest를 만든다.

```go
// +kubebuilder:rbac:groups=example.my.domain,resources=widgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.my.domain,resources=widgets/status,verbs=get;update;patch
```

```bash
controller-gen rbac:roleName=manager-role paths=./controllers/... output:rbac:dir=config/rbac
```

예상 결과:

```text
config/rbac/role.yaml
```

### Webhook YAML

`controller-gen webhook`은 webhook marker를 읽어서 admission webhook configuration manifest를 만든다.

```bash
controller-gen webhook paths=./api/... output:webhook:dir=config/webhook
```

예상 결과:

```text
config/webhook/*.yaml
```

### DeepCopy Go 코드

`controller-gen object`는 Kubernetes runtime object에 필요한 DeepCopy 메서드를 생성한다.

```bash
controller-gen object paths=./api/...
```

예상 결과:

```text
api/v1/zz_generated.deepcopy.go
```

이 코드는 API type이 Kubernetes runtime에서 안전하게 복사될 수 있도록 한다.

## 정리

K8s 생태계에서 `controller-tools`의 위치는 다음과 같다.

- Kubebuilder, Operator SDK, OpenShift Operator 개발 흐름에서 생성 단계를 담당한다.
- CRD 기반 API 프로젝트에서 Go type을 source of truth로 삼게 해 준다.
- 개발자가 작성한 marker 주석을 Kubernetes YAML과 Go 코드로 변환한다.
- 클러스터 런타임 컴포넌트가 아니라 로컬/CI에서 실행되는 빌드 타임 도구다.

## 참고 자료

- [Kubebuilder Book: controller-gen CLI](https://book.kubebuilder.io/reference/controller-gen.html)
- [Kubebuilder Book: Generating CRDs](https://book.kubebuilder.io/reference/generating-crd.html)
- [Operator SDK: generate kustomize manifests](https://sdk.operatorframework.io/docs/cli/operator-sdk_generate_kustomize_manifests/)
- [OpenShift: Developing Operators](https://docs.redhat.com/en/documentation/openshift_container_platform/4.11/html/operators/developing-operators)
- [kubernetes-sigs/controller-tools GitHub README](https://github.com/kubernetes-sigs/controller-tools)
