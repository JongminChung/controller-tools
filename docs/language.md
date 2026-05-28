# controller-tools 언어

```text
CLI option
  -> package loading / marker parsing
  -> generation
  -> output
  -> verification
```

이 흐름은 주로 `controller-gen` 기준이다. 다만 `cmd/` 아래에는 `controller-gen`만 있는 것이 아니라
생성 도구를 보조하는 CLI도 함께 있다.

- `controller-gen`: CRD/RBAC/Webhook/DeepCopy/applyconfiguration 등을 생성하는 대표 CLI
- `helpgen`: marker type의 godoc을 읽어 marker help Go code를 생성하는 보조 CLI
- `type-scaffold`: Kubernetes API type 골격을 stdout으로 생성하는 scaffold CLI

예를 들어 다음 명령은:

```bash
controller-gen rbac:roleName=manager-role paths=./controllers/... output:stdout
```

다음처럼 읽는다.

- `rbac:roleName=manager-role`: RBAC generator와 generator option
- `paths=./controllers/...`: 읽을 Go package 범위
- `output:stdout`: 생성 artifact를 쓸 위치

## 용어 테이블

| 흐름 단계        | 관련 영역                                             | 용어                           | 의미                                                                   | 예시                                                                                      | 헷갈리면 안 되는 것                   | 구분 코멘트                                                                             |
|--------------|---------------------------------------------------|------------------------------|----------------------------------------------------------------------|-----------------------------------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------------------|
| CLI          | `cmd/controller-gen`                              | `controller-gen`             | `controller-tools`의 대표 CLI 실행 파일                                     | `controller-gen crd paths=./api/...`                                                    | controller process            | 클러스터 안에서 reconcile하는 controller가 아니라 로컬/CI에서 실행하는 생성 CLI다.                         |
| CLI          | `cmd/helpgen`                                     | `helpgen`                    | marker type의 godoc에서 marker help Go code를 생성하는 보조 CLI                | `go run ./cmd/helpgen paths=../../pkg/... generate:headerFile=../../boilerplate.go.txt` | `controller-gen`              | 사용자 manifest를 만드는 CLI가 아니라 `zz_generated.markerhelp.go` 같은 help metadata를 만든다.     |
| CLI          | `cmd/type-scaffold`                               | `type-scaffold`              | Kubernetes API type의 기본 Go type 골격을 stdout으로 출력하는 scaffold CLI       | `type-scaffold --kind Foo --resource foos`                                              | `controller-gen object`       | `object` generator는 DeepCopy code를 만들고, `type-scaffold`는 새 API type skeleton을 만든다. |
| CLI          | `cmd/controller-gen/main.go`                      | Raw Option                   | 사용자가 CLI에 넘긴 원본 option 문자열                                           | `crd`, `paths=./api/...`, `output:stdout`                                               | Go flag만으로 파싱되는 값             | Cobra flag도 있지만 generator/output/path option은 marker registry를 통해 파싱된다.            |
| CLI          | `pkg/genall/options.go`                           | Option Marker                | CLI option을 marker 문법으로 해석한 값                                        | `crd`가 내부적으로 `+crd`처럼 처리됨                                                               | source marker                 | CLI option marker는 Go 주석이 아니라 command line 문자열에서 온다.                               |
| CLI          | `pkg/genall/options.go`                           | Common Option                | 모든 generator가 공유하는 option                                            | `paths=./api/...`                                                                       | generator-specific option     | `paths`는 특정 generator만의 설정이 아니라 root package loading 범위를 정한다.                      |
| CLI          | `cmd/controller-gen/main.go`, `pkg/genall`        | Generator                    | marker를 등록하고 artifact를 생성하는 실행 단위                                    | `crd.Generator`, `rbac.Generator`                                                       | controller-runtime reconciler | Generator는 파일/YAML을 만드는 빌드 타임 객체이고, Reconciler는 클러스터 상태를 맞추는 런타임 객체다.              |
| CLI          | `pkg/genall/genall.go`                            | Generation Runtime           | 선택된 generator, root package, collector, output rule을 묶은 실행 단위        | `genall.Runtime`                                                                        | Kubernetes runtime            | 여기서 Runtime은 Kubernetes object runtime이 아니라 한 번의 generation 실행 묶음이다.               |
| CLI          | `pkg/genall/genall.go`                            | Generation Context           | generator 하나가 실행될 때 받는 공유 입력 묶음                                      | `Roots`, `Collector`, `Checker`, `OutputRule`                                           | request context               | 취소/timeout용 context가 아니라 generator 실행에 필요한 입력과 출력 전략이다.                            |
| Marker       | `pkg/markers`                                     | Marker                       | 생성기에 의미를 전달하는 `// +...` 형식의 특수 주석                                    | `// +kubebuilder:validation:Minimum=1`                                                  | Go struct tag                 | Struct tag는 ``json:"name"``처럼 field 뒤에 붙고, marker는 source comment로 붙는다.            |
| Marker       | `pkg/markers`                                     | Marker Definition            | marker 이름, 적용 대상, 파싱 결과 Go type을 정의한 것                               | `markers.Definition`                                                                    | marker value                  | Definition은 parser 규칙이고, value는 실제 source에서 읽은 값이다.                                |
| Marker       | `pkg/markers`                                     | Registry                     | 사용 가능한 marker definition을 등록해 두는 저장소                                 | `markers.Registry`                                                                      | Collector                     | Registry는 “무엇을 파싱할 수 있는지”를 알고, Collector는 “source에서 무엇을 찾았는지”를 모은다.                |
| Marker       | `pkg/markers`                                     | Collector                    | Go package에서 marker 주석을 찾아 Go 값으로 모으는 객체                             | `markers.Collector`                                                                     | loader                        | Collector는 marker를 수집하고, loader는 Go package와 AST/type 정보를 로드한다.                    |
| Marker       | `pkg/markers`                                     | Marker Target                | marker가 붙을 수 있는 위치                                                   | package, type, field                                                                    | Kubernetes target object      | 여기서 target은 YAML 대상 리소스가 아니라 Go source 안의 부착 위치다.                                  |
| Marker       | `pkg/markers`                                     | Marker Value                 | 파싱된 marker 결과 Go 값                                                   | `Minimum=1`이 숫자 값으로 변환됨                                                                 | raw marker string             | Raw string을 Definition으로 parse한 뒤 generator가 쓰는 구조화된 값이다.                          |
| Marker       | `pkg/markers`, `pkg/crd`                          | Package Marker               | package 전체에 적용되는 marker                                              | `// +groupName=batch.example.com`                                                       | Go package name               | Go package 이름과 별개로 generator가 package scope에서 읽는 marker다.                          |
| Marker       | `pkg/markers`, `pkg/crd`                          | Type Marker                  | type 선언에 적용되는 marker                                                 | `// +kubebuilder:resource:path=cronjobs`                                                | field marker                  | Type 전체의 resource 성격이나 CRD 설정을 표현한다.                                               |
| Marker       | `pkg/markers`, `pkg/crd`                          | Field Marker                 | struct field에 적용되는 marker                                            | `// +kubebuilder:validation:MaxLength=63`                                               | type marker                   | 특정 field의 validation/schema 속성을 표현한다.                                              |
| Marker       | `cmd/helpgen`, `pkg/*/zz_generated.markerhelp.go` | Marker Help                  | marker option 설명과 field 설명을 담은 generated help metadata               | `Help() *markers.DefinitionHelp`                                                        | marker definition             | Definition은 marker parse 규칙이고, Marker Help는 사용자에게 보여 줄 설명 문구다.                     |
| Loader       | `pkg/loader`                                      | Root Package                 | 사용자가 `paths`로 지정한 시작 package                                         | `paths=./api/...`                                                                       | output directory              | `paths`는 입력 source 위치이고, 출력 위치는 `output:...`이 정한다.                                 |
| Loader       | `pkg/loader`                                      | Package                      | `go/packages` 결과를 감싼 controller-tools 내부 package 객체                  | `loader.Package`                                                                        | 배포 package, Java package      | Go package loading 단위다. 배포 단위는 Go module이고 Java package와도 규칙이 다르다.                 |
| Loader       | `pkg/loader`                                      | Syntax                       | Go source의 AST 정보                                                    | `NeedSyntax()`                                                                          | type info                     | Syntax는 문법 트리이고, type info는 type checker가 계산한 의미 정보다.                              |
| Loader       | `pkg/loader`                                      | Type Info                    | Go type checker가 계산한 타입 정보                                           | `NeedTypesInfo()`                                                                       | schema                        | Type Info는 Go 타입 분석 결과이고, Schema는 CRD YAML에 들어가는 OpenAPI 표현이다.                     |
| Loader       | `pkg/loader`                                      | Type Checker                 | 필요한 package와 type만 부분적으로 type-checking하는 도구                          | `loader.TypeChecker`                                                                    | `go test`                     | 전체 테스트 실행기가 아니라 generator가 필요한 타입 정보를 계산하는 도구다.                                    |
| Loader       | `pkg/loader`                                      | Node Filter                  | type-checking할 관심 node를 고르는 필터                                       | CRD generator가 필요한 type 중심으로 검사                                                         | Kubernetes field selector     | Go AST/type graph에서 검사 범위를 좁히는 필터다.                                                |
| Loader       | `pkg/loader`                                      | Package Error                | package에 누적된 위치 기반 error                                             | marker validation 실패, parsing error                                                     | process exit error            | source 위치를 포함해 나중에 출력되는 package-level error다.                                      |
| Generator    | `pkg/crd`                                         | CRD Generator                | Go API type에서 `CustomResourceDefinition` YAML을 만드는 generator         | `controller-gen crd`                                                                    | API server                    | CRD를 생성하는 빌드 타임 도구이고, CRD를 검증/서빙하는 것은 Kubernetes API server다.                      |
| Generator    | `pkg/rbac`                                        | RBAC Generator               | RBAC marker에서 `Role`/`ClusterRole` YAML을 만드는 generator               | `controller-gen rbac:roleName=manager-role`                                             | Kubernetes authorization      | 권한을 판정하지 않고 권한 manifest를 생성한다.                                                     |
| Generator    | `pkg/webhook`                                     | Webhook Generator            | webhook marker에서 admission webhook configuration YAML을 만드는 generator | `controller-gen webhook`                                                                | webhook server                | Webhook 서버를 실행하지 않고 webhook configuration manifest를 만든다.                           |
| Generator    | `pkg/deepcopy`                                    | Object Generator             | Kubernetes object용 DeepCopy Go code를 만드는 generator                   | `controller-gen object`                                                                 | Kubernetes object instance    | 객체 인스턴스를 만드는 것이 아니라 object type에 필요한 복사 메서드를 생성한다.                                 |
| Generator    | `pkg/applyconfiguration`                          | ApplyConfiguration Generator | Server-Side Apply용 apply configuration Go code를 만드는 generator        | `controller-gen applyconfiguration paths=./api/...`                                     | CRD schema                    | CRD YAML이 아니라 apply configuration 타입 코드를 생성한다.                                     |
| Generator    | `pkg/schemapatcher`                               | Schema Patcher               | 기존 CRD manifest에 새 schema를 patch하는 generator                         | `controller-gen schemapatch:manifests=./manifests`                                      | CRD Generator                 | 새 CRD 전체를 생성하기보다 기존 manifest의 schema를 갱신하는 흐름이다.                                   |
| Generator    | `pkg/typescaffold`                                | Type Scaffold                | Kind/Resource/Namespaced 설정으로 Kubernetes API type skeleton을 만드는 기능   | `typescaffold.ScaffoldOptions`                                                          | CRD Generator                 | CRD YAML을 생성하지 않고 CRD의 입력이 될 수 있는 Go type 골격을 만든다.                                 |
| Generator    | `pkg/crd`                                         | Schema                       | CRD에 들어가는 OpenAPI v3 schema                                          | `spec.versions[].schema.openAPIV3Schema`                                                | Go type                       | Go type은 입력 모델이고, Schema는 Kubernetes API server가 읽는 YAML 표현이다.                     |
| Generator    | `pkg/crd`                                         | Plural Path                  | Kubernetes REST resource path에 쓰는 복수형 이름                             | `Widget`의 `widgets`                                                                     | Kind, singular name           | `Kind`는 `Widget`, singular는 `widget`, plural path는 API 경로의 `widgets`다.             |
| Generator    | `pkg/rbac`                                        | Rule                         | Kubernetes RBAC policy rule                                          | `verbs=get;list;watch`                                                                  | business rule                 | 여기서는 권한 manifest의 `PolicyRule`을 뜻한다. 도메인 정책 규칙이 아니다.                               |
| Generator    | `pkg/genall`, generator packages                  | Artifact                     | generator가 만들어 내는 결과물                                                | CRD YAML, RBAC YAML, DeepCopy Go code                                                   | source input                  | Go source와 marker는 입력이고 artifact는 생성된 출력이다.                                        |
| Output       | `pkg/genall/output.go`                            | Output Rule                  | artifact를 어디에 쓸지 결정하는 전략                                             | `output:stdout`, `output:dir=...`                                                       | input `paths`                 | `paths`는 읽을 위치이고, Output Rule은 쓸 위치다.                                              |
| Output       | `pkg/genall/output.go`                            | Default Output               | generator별 지정이 없을 때 쓰는 출력 규칙                                         | `output:stdout`                                                                         | per-generator output          | 모든 generator에 fallback으로 적용된다.                                                     |
| Output       | `pkg/genall/output.go`                            | Per-generator Output         | 특정 generator에만 적용하는 출력 규칙                                            | `output:crd:dir=config/crd/bases`                                                       | default output                | `crd`처럼 지정된 generator에만 우선 적용된다.                                                   |
| Output       | `pkg/genall/output.go`                            | Config Artifact              | Go package에 속하지 않는 Kubernetes 설정 생성물                                 | CRD YAML, RBAC YAML                                                                     | code artifact                 | Kubernetes manifest는 보통 `config/` 아래에 놓이고 Go compile 대상이 아니다.                      |
| Output       | `pkg/genall/output.go`                            | Code Artifact                | Go package에 속하는 생성 코드                                                | `zz_generated.deepcopy.go`                                                              | config artifact               | Go package 안에 생성되어 컴파일 대상이 된다.                                                     |
| Output       | `pkg/genall/output.go`                            | Output Artifacts             | config와 code 출력을 나눠서 배치하는 출력 방식                                      | `output:artifacts:config=config/crd/bases`                                              | `output:dir`                  | `output:dir`은 단일 디렉터리 전략이고, artifacts는 config/code 배치를 분리한다.                       |
| Verification | `pkg/*/testdata`                                  | Testdata Module              | generator나 loader를 검증하기 위한 샘플 Go module                              | `pkg/crd/testdata/go.mod`                                                               | production module             | 테스트 fixture로 만든 module이며 루트 module과 배포 의미가 다르다.                                    |
| Verification | `pkg/*/testdata`                                  | Golden File                  | 생성 결과와 비교하는 기대 출력 파일                                                 | `pkg/rbac/testdata/role.yaml`                                                           | generated source of truth     | Golden file은 기대값이다. source of truth는 보통 Go type과 marker 입력이다.                      |
| Verification | `pkg/*/*_test.go`                                 | Integration Test             | 실제 package loading과 generation을 포함하는 테스트                             | `parser_integration_test.go`                                                            | unit test                     | parser 함수만 보는 테스트보다 실제 생성 흐름을 더 많이 포함한다.                                           |
| Verification | local command                                     | Narrow Test                  | 변경한 package만 빠르게 검증하는 테스트                                            | `go test ./pkg/rbac`                                                                    | full verification             | 빠른 피드백용이고 PR 전 전체 검증을 대체하지 않는다.                                                    |
| Verification | `test.sh`                                         | Full Verification            | PR 전 넓은 검증                                                           | `./test.sh`                                                                             | narrow test                   | envtest, generate, lint, race test까지 포함할 수 있는 최종 검증 흐름이다.                          |
| Verification | `go generate`, `.run-controller-gen.sh`           | Regeneration                 | golden output이나 generated code를 다시 만드는 작업                            | `go generate ./pkg/crd/testdata`                                                        | manual edit                   | 생성 파일은 손으로 맞추기보다 generator로 재생성하고 diff를 설명한다.                                      |

## 짧은 예시

### CLI option에서 Generation Runtime까지

```bash
controller-gen crd paths=./api/... output:crd:dir=config/crd/bases
```

이 명령은 다음처럼 읽는다.

```text
controller-gen은 CRD Generator를 선택한다.
paths option은 Root Package 범위를 정한다.
output:crd:dir은 CRD Generator에만 적용되는 Per-generator Output이다.
genall.Runtime은 이 정보를 묶고 Generator.Generate를 호출한다.
```

### Marker가 CRD schema로 바뀌는 흐름

Go source:

```go
// +kubebuilder:resource:path=widgets,scope=Namespaced
type Widget struct {
Spec WidgetSpec `json:"spec,omitempty"`
}

type WidgetSpec struct {
// +kubebuilder:validation:Minimum=1
Replicas int32 `json:"replicas"`
}
```

생성되는 CRD schema 일부:

```yaml
spec:
  names:
    kind: Widget
    plural: widgets
  scope: Namespaced
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                replicas:
                  type: integer
                  minimum: 1
```

여기서 `path=widgets`는 plural path가 되고, `Minimum=1` marker value는 OpenAPI schema의
`minimum: 1`이 된다.

### RBAC marker가 ClusterRole YAML로 바뀌는 흐름

Go source:

```go
// +kubebuilder:rbac:groups=example.my.domain,resources=widgets,verbs=get;list;watch
// +kubebuilder:rbac:groups=example.my.domain,resources=widgets/status,verbs=get;update;patch
```

CLI:

```bash
controller-gen rbac:roleName=manager-role paths=./controllers/... output:stdout
```

생성되는 RBAC YAML 일부:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: manager-role
rules:
  - apiGroups:
      - example.my.domain
    resources:
      - widgets
    verbs:
      - get
      - list
      - watch
```

여기서 `roleName` option은 artifact의 `metadata.name`이 되고, marker의 `groups`, `resources`,
`verbs`는 RBAC `PolicyRule`이 된다.

### Output Rule이 출력 위치를 정하는 흐름

```bash
controller-gen rbac:roleName=manager-role crd paths=./api/... output:crd:dir=config/crd/bases output:stdout
```

이 명령은 다음처럼 읽는다.

```text
CRD artifact는 config/crd/bases 디렉터리에 쓴다.
RBAC artifact는 별도 per-generator output이 없으므로 default output인 stdout에 쓴다.
```

## 용어 사용 규칙

- `Generator`는 `crd`, `rbac`, `object` 같은 생성 기능 단위를 가리킬 때만 쓴다.
- `Runtime`이라고만 쓰면 Kubernetes runtime과 헷갈릴 수 있으므로, `genall.Runtime`을 뜻할 때는 `Generation Runtime`이라고 쓴다.
- `Marker`는 Go struct tag가 아니라 `// +...` 형태의 comment marker를 뜻한다.
- `Package`는 Go package loading 단위를 뜻한다. Go module, 배포 package, Java package와 구분한다.
- `Artifact`는 생성 결과물이다. Go source와 marker input을 artifact라고 부르지 않는다.
- `Schema`는 CRD OpenAPI schema를 뜻할 때가 많다. Go type 자체와 구분한다.
- `Plural path`는 Kubernetes REST resource path에 쓰는 복수형 이름이다. 예: `Widget`의 plural path는 `widgets`.
