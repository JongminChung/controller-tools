# controller-tools 기여 로드맵

이 문서는 Kubernetes 경험은 있지만 Go는 초급인 기여자를 기준으로,
`controller-tools`를 읽고 작은 기여 PR까지 가는 진행 체크리스트다.

## 진행 상태

- 현재 단계: `1주차: CLI와 공통 generator 흐름`
- 다음 행동: `cmd/controller-gen/main.go`에서 CLI option이 `genall.Runtime`으로 바뀌는 흐름 설명하기
- 기여 목표: `pkg/rbac`, `pkg/webhook`, `pkg/markers` 중 하나에 작고 검증 가능한 변경 만들기
- 기여 후보: 아직 미정

상태 표기는 다음처럼 관리한다.

- `[ ]`: 아직 하지 않음
- `[x]`: 완료
- `(진행 중)`: 아직 완료하지 않았지만 현재 보고 있는 항목
- `(재확인 필요)`: 실행은 했지만 설명이나 근거 정리가 더 필요한 항목

## 최종 목표

- [ ] `controller-gen`이 CLI 옵션을 받아 generator를 실행하는 흐름을 설명할 수 있다.
- [ ] marker comment가 Go 값으로 파싱되고 generator 입력으로 쓰이는 과정을 설명할 수 있다.
- [ ] `go/packages`, `go/ast`, `go/types`가 이 프로젝트에서 어떤 역할을 하는지 설명할 수 있다.
- [ ] `rbac`, `webhook`, `crd` generator 중 하나에 작은 변경과 테스트를 추가할 수 있다.
- [ ] golden output 변경이 필요한 PR에서 입력, 생성 결과, diff 이유를 설명할 수 있다.
- [ ] 작은 PR 후보를 고르고, 테스트와 검증 결과를 포함해 기여 준비를 끝낸다.

## controller-tools는 무엇인가

- [x] `controller-tools`의 역할을 한 문장으로 설명한다.

`controller-tools`는 Kubernetes API/Controller 프로젝트에서 필요한 생성물을
Go 소스 코드와 marker 주석을 읽어서 만들어 주는 도구 모음이다.

대표 실행 파일은 `controller-gen`이다. `controller-gen`은 controller를 직접 실행하는
런타임이 아니라, 개발 또는 빌드 시점에 실행하는 generator다. 사람이 직접 길게 작성하기
어려운 Kubernetes YAML이나 반복적인 Go 코드를 프로젝트의 Go 타입과 marker에서 만들어 낸다.

주요 generator는 다음과 같다.

- `crd`: Go API type을 읽어서 `CustomResourceDefinition` YAML을 만든다.
- `rbac`: `+kubebuilder:rbac` marker를 읽어서 `ClusterRole` 또는 `Role` YAML을 만든다.
- `object`: Kubernetes runtime object에 필요한 `DeepCopy`, `DeepCopyInto`, `DeepCopyObject` 코드를 만든다.
- `webhook`: webhook marker를 읽어서 webhook configuration YAML을 만든다.
- `applyconfiguration`: Server-Side Apply에 쓰는 apply configuration Go 코드를 만든다.
- `schemapatch`: 기존 CRD manifest에 새 schema를 patch한다.

큰 흐름은 다음과 같다.

```text
Go source code
+ marker comments
        |
        v
controller-gen
        |
        v
CRD YAML / RBAC YAML / Webhook YAML / DeepCopy Go code
```

## 실제 사용하는 느낌

- [x] CRD generator가 Go type과 marker에서 CRD YAML을 만드는 흐름을 이해한다.
- [x] RBAC generator가 `+kubebuilder:rbac` marker에서 RBAC YAML을 만드는 흐름을 이해한다.
- [x] DeepCopy generator가 반복 Go 코드를 생성하는 위치를 이해한다.
- [x] `Taskfile.yml`로 실제 generator 실행 결과를 확인한다.

Kubernetes API type을 작성할 때 Go struct 위에 marker를 붙인다.

```go
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=cronjobs,scope=Namespaced
// +kubebuilder:subresource:status
type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronJobSpec   `json:"spec,omitempty"`
	Status CronJobStatus `json:"status,omitempty"`
}
```

필드에는 OpenAPI validation에 해당하는 marker를 붙일 수 있다.

```go
type CronJobSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas int32 `json:"replicas"`
}
```

그다음 CRD generator를 실행한다.

```bash
controller-gen crd paths=./api/... output:crd:dir=config/crd/bases
```

그러면 대략 다음 위치에 CRD YAML이 생성된다.

```text
config/crd/bases/example.my.domain_cronjobs.yaml
```

RBAC도 같은 방식이다. controller 코드 근처에 필요한 권한을 marker로 적는다.

```go
// +kubebuilder:rbac:groups=example.my.domain,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.my.domain,resources=cronjobs/status,verbs=get;update;patch
```

그리고 RBAC generator를 실행한다.

```bash
controller-gen rbac:roleName=manager-role paths=./controllers/... output:rbac:dir=config/rbac
```

이 저장소에서 직접 흐름을 볼 때는 `Taskfile.yml`의 학습 태스크를 사용한다.

```bash
task learning:main:rbac-stdout
task learning:main:rbac-dir
task learning:main:crd-dir
```

확인한 생성물:

- [x] `artifacts/controller-gen-learning/rbac/role.yaml`
- [x] `artifacts/controller-gen-learning/crd/testdata.kubebuilder.io_cronjobs.yaml`

## 1주차: CLI와 공통 generator 흐름

목표: `controller-gen` CLI가 raw args를 받아 어떤 generator를 실행할지 결정하고,
`genall.Runtime`을 통해 실제 생성 작업을 호출하는 흐름을 설명한다.

### 읽을 파일

- [ ] `cmd/controller-gen/main.go` (진행 중)
- [ ] `pkg/genall/options.go` (진행 중)
- [ ] `pkg/genall/genall.go`
- [ ] `pkg/genall/output.go`

### 확인할 것

- [x] `allGenerators`가 어떤 generator들을 CLI 이름에 연결하는지 확인한다.
- [x] `allOutputRules`가 출력 방식을 어떻게 등록하는지 확인한다.
- [ ] `optionsRegistry`에 generator, output rule, common option이 marker definition으로 들어가는 흐름을 따라간다. (진행 중)
- [ ] `genall.FromOptionsWithConfig`가 raw CLI args를 `Runtime`으로 바꾸는 지점을 찾는다. (진행 중)
- [ ] `rt.Run()`이 generator별 `Generate`를 호출하는 흐름을 설명한다.
- [x] generator 없이 `paths`만 주면 `no generators specified`가 나는 것을 확인한다.

### 핵심 타입

- [x] `Generator`: marker 등록과 artifact 생성을 담당하는 인터페이스.
- [ ] `Runtime`: root package, marker collector, output rule, generator 목록을 묶는 실행 단위.
- [ ] `GenerationContext`: 각 generator가 공유해서 쓰는 입력, 출력, type checker 컨텍스트.
- [ ] `OutputRule`: YAML이나 Go 파일을 어디에 쓸지 결정하는 출력 전략.

### 실행 체크

```bash
task learning:main:source
task learning:main:test-genall
task learning:main:help
task learning:main:no-generator
task learning:main:markers-crd
```

- [x] `task learning:main:test-genall` 실행
- [x] `task learning:main:no-generator` 실행
- [x] `task learning:main:markers-crd` 실행
- [ ] `task learning:main:source` 출력에서 `main.go`와 `genall/options.go` 연결 설명하기
- [ ] `task learning:main:help` 출력에서 등록된 generators/options 설명하기

### 완료 기준

- [ ] `controller-gen crd paths=... output:stdout`가 어떤 코드 경로로 실행되는지 말로 설명한다.
- [ ] CLI option도 내부적으로 marker parser를 사용한다는 점을 설명한다.
- [ ] 새 generator를 CLI에 추가하려면 어디를 수정해야 하는지 설명한다.

## 2주차: marker와 loader

목표: source comment marker가 Go 값으로 파싱되고, package/type/field 단위로 수집되는 흐름을 이해한다.

### `pkg/markers`

읽는 순서:

- [ ] `pkg/markers/doc.go`
- [ ] `pkg/markers/reg.go`
- [ ] `pkg/markers/parse.go`
- [ ] `pkg/markers/collect.go`
- [ ] `pkg/markers/parse_test.go`
- [ ] `pkg/markers/collect_test.go`

확인할 것:

- [ ] `+kubebuilder:validation:...` 같은 주석이 어떤 grammar로 파싱되는지 확인한다.
- [ ] `Definition`, `Registry`, `Collector`, `MarkerValues`의 책임을 구분한다.
- [ ] struct tag `marker:",optional"`과 pointer optional field가 어떻게 다르게 쓰이는지 확인한다.
- [ ] package, type, field marker가 같은 이름을 가질 수 있는 이유를 설명한다.

### `pkg/loader`

읽는 순서:

- [ ] `pkg/loader/doc.go`
- [ ] `pkg/loader/loader.go`
- [ ] `pkg/loader/visit.go`
- [ ] `pkg/loader/refs.go`
- [ ] `pkg/loader/loader_test.go`

확인할 것:

- [ ] `LoadRoots`가 `go/packages` 결과를 controller-tools의 `Package`로 감싸는 방식을 확인한다.
- [ ] `NeedSyntax()`와 `NeedTypesInfo()`가 lazy하게 AST와 type info를 채우는 흐름을 확인한다.
- [ ] `TypeChecker.Check()`가 필요한 package만 partial type-checking하는 이유를 설명한다.
- [ ] `Package.AddError`가 파일 위치가 있는 error로 바뀌는 방식을 확인한다.

실행 체크:

```bash
go test ./pkg/markers ./pkg/loader
go test ./pkg/markers -run TestDoesNotExist
```

- [ ] `go test ./pkg/markers ./pkg/loader` 실행
- [ ] Ginkgo suite와 일반 `testing`의 차이를 확인한다.

## 3주차: 작은 generator로 기여 감각 만들기

목표: `pkg/rbac` 또는 `pkg/webhook` 중 하나를 골라 marker 입력이 Kubernetes object로 변환되는 흐름을 따라간다.

추천 시작점은 `pkg/rbac`이다. 코드가 상대적으로 작고, golden output이 명확하다.

### `pkg/rbac`

- [ ] `RuleDefinition`이 package marker로 등록되는 방식을 확인한다.
- [ ] `Rule` struct의 marker tag와 Kubernetes `PolicyRule` 변환 흐름을 확인한다.
- [ ] rule merge, dedup, sort가 golden output 안정성에 왜 중요한지 설명한다.
- [ ] `parser_integration_test.go`가 testdata module을 로드하고 `role.yaml`과 비교하는 방식을 확인한다.

실행 체크:

```bash
go test ./pkg/rbac
go run ./cmd/controller-gen rbac:roleName=manager-role paths=./pkg/rbac/testdata output:stdout
task learning:main:rbac-stdout
task learning:main:rbac-dir
```

- [x] `task learning:main:rbac-stdout` 실행
- [x] `task learning:main:rbac-dir` 실행
- [ ] `go test ./pkg/rbac` 실행
- [ ] `pkg/rbac/testdata/role.yaml`과 생성 결과를 비교한다.

### `pkg/webhook`

- [ ] webhook marker가 `MutatingWebhookConfiguration` 또는 `ValidatingWebhookConfiguration`으로 변환되는 흐름을 확인한다.
- [ ] validation error가 package error로 누적되는 방식을 확인한다.
- [ ] `testdata/*/manifests.yaml`이 golden file 역할을 하는 방식을 확인한다.

실행 체크:

```bash
go test ./pkg/webhook
```

- [ ] `go test ./pkg/webhook` 실행

### 완료 기준

- [ ] `rbac` 또는 `webhook` 중 하나에서 입력 marker, 내부 struct, 출력 YAML의 연결을 설명한다.
- [ ] 작은 변경 후보 하나를 고른다.
- [ ] 변경 후보가 test-only인지, behavior 변경인지, help 문구 변경인지 분류한다.

## 4주차: CRD generator 깊게 읽기

목표: 입력 Go type 하나가 CRD YAML schema로 변환되는 경로를 추적한다.

`pkg/crd`는 controller-tools의 핵심이자 가장 복잡한 영역이다. 바로 전체를 외우려고 하지 말고,
입력 Go type 하나가 CRD YAML이 되는 경로를 추적한다.

읽는 순서:

- [ ] `pkg/crd/gen.go`
- [ ] `pkg/crd/parser.go`
- [ ] `pkg/crd/spec.go`
- [ ] `pkg/crd/schema.go`
- [ ] `pkg/crd/flatten.go`
- [ ] `pkg/crd/markers`
- [ ] `pkg/crd/parser_integration_test.go`

따라갈 대표 흐름:

- [ ] `Generator.Generate`가 `Parser`를 만들고 root package를 `NeedPackage`로 등록한다.
- [ ] `FindKubeKinds`가 Kubernetes resource type을 찾는다.
- [ ] `NeedCRDFor`가 group-kind에 대한 CRD spec을 만든다.
- [ ] `NeedSchemaFor`와 `infoToSchema`가 Go type을 OpenAPI schema로 바꾼다.
- [ ] `applyMarkers`가 validation marker를 schema에 적용한다.
- [ ] `FlattenedSchemata`가 CRD에 들어갈 reference 없는 schema로 정리된다.
- [ ] `WriteYAML`이 최종 YAML을 쓴다.

실습 체크:

- [ ] `pkg/crd/testdata/cronjob_types.go`에서 필드 하나를 고른다.
- [ ] 그 필드의 Go type, JSON tag, kubebuilder marker를 확인한다.
- [ ] 생성된 `pkg/crd/testdata/testdata.kubebuilder.io_cronjobs.yaml`에서 해당 schema 위치를 찾는다.
- [ ] `schema.go`의 `typeToSchema`, `structToSchema`, `applyMarkers` 중 어디를 거쳤는지 추적한다.

실행 체크:

```bash
go test ./pkg/crd
go test ./pkg/crd -run TestDoesNotExist
task learning:main:crd-dir
```

- [x] `task learning:main:crd-dir` 실행
- [ ] `go test ./pkg/crd` 실행
- [ ] CRD golden output 변경이 필요한 경우를 설명한다.

## 기여 후보 고르기

처음 기여는 영향 범위가 좁고 검증이 쉬운 일을 고른다.

- [ ] `pkg/rbac` test 추가 또는 marker help 개선
- [ ] `pkg/webhook` validation edge case test 추가
- [ ] `pkg/markers` parser edge case test 추가
- [ ] `pkg/crd/markers`의 좁은 marker 동작 test 추가
- [ ] `pkg/crd/schema.go`의 실제 schema generation 변경

처음부터 `loader`나 `crd` 전체 구조를 바꾸는 PR은 피한다. 이 영역은 영향 범위가 넓고,
기존 partial type-checking과 golden output 안정성에 대한 이해가 필요하다.

기여 후보를 고를 때 기록할 것:

- [ ] 어떤 package를 바꾸는가
- [ ] 새 동작이 marker parse, generator transform, output serialization 중 어디에 속하는가
- [ ] golden file이 있는가
- [ ] 실패 케이스 test가 필요한가
- [ ] help text 또는 generated marker help 갱신이 필요한가

## PR 준비 체크리스트

변경 전에:

- [ ] 관련 package의 integration test 구조를 먼저 읽는다.
- [ ] golden file이 있는지 확인한다.
- [ ] 새 동작이 marker parse, generator transform, output serialization 중 어디에 속하는지 정한다.
- [ ] 작은 실패 재현 또는 test case를 먼저 만든다.

변경 중:

- [ ] 기존 사용자 변경사항을 되돌리지 않는다.
- [ ] 변경 범위를 선택한 package와 testdata에 가깝게 유지한다.
- [ ] error message를 바꾸면 실패 케이스 test를 함께 추가한다.
- [ ] marker public behavior를 바꾸면 help text 갱신 필요 여부를 확인한다.

변경 후:

- [ ] 좁은 테스트를 먼저 실행한다.
- [ ] golden output이 바뀌면 README 지시에 따라 재생성하고 diff를 검토한다.
- [ ] 관련 package test를 실행한다.
- [ ] 최종 검증 명령을 실행하거나, 실행하지 못한 이유를 기록한다.

최종 확인:

```bash
go test ./pkg/markers ./pkg/loader ./pkg/genall
go test ./pkg/rbac ./pkg/webhook ./pkg/crd
./test.sh
```

`./test.sh`는 envtest binary 다운로드, `go generate`, lint, race test를 포함하므로
PR 전 최종 검증으로 사용한다.

## 진행 기록

| 날짜 | 단계 | 한 일 | 다음 행동 |
| --- | --- | --- | --- |
| 2026-05-27 | 사전 이해 | `controller-tools` 역할, 실제 사용 흐름, `Taskfile.yml` 기반 실행 흐름 정리 | `main.go`에서 `FromOptionsWithConfig`와 `Runtime.Run` 연결 설명 |
