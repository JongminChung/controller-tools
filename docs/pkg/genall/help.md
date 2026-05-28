# `pkg/genall/help` 분석

`pkg/genall/help`는 generator 실행 로직이 아니라 help 출력용 모델과 포매터를 제공하는 패키지다.
생성물을 만드는 핵심 경로가 아니라, `controller-gen -h`, `controller-gen -w`, `helpgen`을 이해할 때 보면 된다.

## 한 줄 요약

```text
markers.Registry에 등록된 marker definition + help metadata
        |
        v
pkg/genall/help
        |
        v
controller-gen -h / -w 출력 또는 JSON help
```

## 언제 쓰이나

`controller-gen`에서 사용자가 help를 요청할 때 쓰인다.

```bash
controller-gen -h
controller-gen crd -w
controller-gen crd -ww
controller-gen crd -wwww
```

`-h`는 CLI option 도움말이다.

```text
generators
+crd ...
+rbac ...
+webhook ...

generic
+paths=...

output rules
+output:stdout
+output:dir=...
```

`-w`는 선택한 generator가 이해하는 marker 도움말이다.

```text
+kubebuilder:resource
+kubebuilder:validation:Minimum
+kubebuilder:subresource:status
```

## 코드에서의 연결

`cmd/controller-gen/main.go`는 help level에 따라 marker registry를 help 문서로 바꾼다.

```go
helpInfo := help.ByCategory(reg, sorter)
```

터미널 출력은 `pretty` 하위 패키지가 담당한다.

```go
contents := prettyhelp.MarkersSummary(cat.Category, cat.Markers)
contents := prettyhelp.MarkersDetails(fullDetail, cat.Category, cat.Markers)
```

정렬 기준은 두 가지가 있다.

- `help.SortByOption`: `controller-gen -h`처럼 generator option 중심으로 묶는다.
- `help.SortByCategory`: `controller-gen crd -w`처럼 marker category 중심으로 묶는다.

## 디렉터리 구조

```text
pkg/genall/help/
  doc.go        help package 설명
  types.go     MarkerDoc, FieldHelp, CategoryDoc 같은 help용 데이터 모델
  sort.go      help 항목 정렬/그룹핑 규칙
  pretty/
    doc.go     terminal formatting package 설명
    help.go    marker help를 터미널용 summary/detail로 포맷
    print.go   색상, Span, Table 출력 유틸
    table.go   terminal table width 계산
```

## 중요한 개념

`pkg/genall/help`의 데이터는 generator 입력값이 아니라 사용자에게 보여줄 문서화용 중간 표현이다.

예를 들어 `MarkerDoc`은 marker definition과 help metadata를 합친 표시용 모델이다.

```go
type MarkerDoc struct {
	Name     string
	Target   string
	Category string
	Fields   []FieldHelp
}
```

이 값은 `crd.Generator.Generate` 같은 생성 로직에 직접 들어가는 값이 아니다. 사용자가 `-h` 또는 `-w`를 실행했을 때
터미널이나 JSON으로 보여주기 위해 만든 구조다.

## `helpgen`과의 관계

`cmd/helpgen`도 `pkg/genall/help`를 사용한다.

`helpgen`은 marker type의 godoc을 읽어서 generated help 파일을 만든다.

```text
pkg/crd/zz_generated.markerhelp.go
pkg/rbac/zz_generated.markerhelp.go
pkg/webhook/zz_generated.markerhelp.go
```

전체 흐름은 다음과 같다.

```text
marker type의 Go doc
        |
        v
cmd/helpgen
        |
        v
zz_generated.markerhelp.go
        |
        v
markers.Registry에 help metadata 등록
        |
        v
pkg/genall/help.ByCategory
        |
        v
pretty.MarkersSummary / pretty.MarkersDetails
        |
        v
controller-gen -h / -w 출력
```

## `-h`와 `-w`를 읽는 기준

`controller-gen -h`는 CLI option 관점이다.

```bash
controller-gen -h
```

이때는 다음을 확인한다.

- 어떤 generator를 실행할 수 있는가
- 공통 option은 무엇인가
- output rule은 어떤 것이 있는가

`controller-gen crd -w`는 marker 관점이다.

```bash
controller-gen crd -w
```

이때는 다음을 확인한다.

- CRD generator가 어떤 marker를 이해하는가
- marker가 package/type/field 중 어디에 붙는가
- marker argument가 필수인지 optional인지

## 학습 우선순위

`pkg/genall/help`는 presentation layer에 가깝다. 생성 흐름을 처음 읽을 때는 우선순위가 낮다.

먼저 다음을 읽는 것이 좋다.

```text
pkg/genall/options.go
pkg/genall/genall.go
pkg/genall/output.go
pkg/markers
pkg/loader
pkg/rbac
```

그 다음 `controller-gen -h`, `controller-gen -w`, `helpgen`, `zz_generated.markerhelp.go`가 궁금할 때
`pkg/genall/help`를 보면 된다.

## 헷갈리면 안 되는 것

| 용어 | 헷갈리면 안 되는 것 | 구분 |
| --- | --- | --- |
| `pkg/genall/help` | generator 실행 엔진 | help 모델과 표시 로직만 다룬다. |
| `MarkerDoc` | marker parse 결과 | 사용자가 읽는 help 문서용 모델이다. |
| `pretty` | generator output pretty print | CLI help를 보기 좋게 출력하는 terminal formatter다. |
| `helpgen` | `controller-gen` | `helpgen`은 marker help Go code를 생성하는 보조 CLI다. |
| `zz_generated.markerhelp.go` | CRD/RBAC/Webhook manifest | marker 설명 metadata를 담는 generated Go code다. |
