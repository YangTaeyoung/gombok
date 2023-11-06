# Hello Gombok! 👋
안녕하세요! 해당 도구는 Golang에서 Java Library 중 하나인 [Lombok](https://projectlombok.org/)과 유사한 기능을 제공하는 도구입니다.
해당 도구를 사용하면, Golang에서 주석을 활용한 어노테이션을 이용해 Getter, Setter, Constructor, Builder 등 자주 사용하는 유틸함수를 쉽게 생성하여 사용할 수 있습니다.

[🇺🇸 English](./README_en.md)

# Install With Go
- [혹시 Go가 설치되어 있지 않으신가요?](https://go.dev/dl/)

터미널에 다음 명령어를 입력합니다.
```bash
$ go install github.com/YangTaeyoung/gombok@v1.1.0
```

# Usage
1. 어노테이션을 스캔할 폴더에 들어갑니다.
    ```bash
    cd <path>
    ```
   - 다음과 같은 파일이 있다고 가정합니다.
        ```go
        // in ./test/test.go
		
        // @Builder
        type Test struct {
            Name string
            Age  int
        }
        ```
2. 아래의 명령어를 실행합니다.
    ```bash
    $ gombok
    ```
3. gombok는 `test` 폴더에 있는 모든 Go 파일들을 스캔하여, 어노테이션을 찾아 다음과 같은 파일을 생성합니다.
    ```go
    // in ./test/test_gombok.go
   
   // DO NOT EDIT THIS FILE. THIS FILE WAS AUTO GENERATED BY GOMBOK.
    package test
    
    // TestBuilder is a builder for Test
    type TestBuilder struct {
        target *Test
    }
    
    // SetName sets the Name field of the target Test
    func (tb TestBuilder) WithName(name string) TestBuilder {
        tb.target.Name = name
    
        return tb
    }
    
    // SetAge sets the Age field of the target Test
    func (tb TestBuilder) WithAge(age int) TestBuilder {
        tb.target.Age = age
    
        return tb
    }
    
    // Build constructs a Test from the builder
    func (tb TestBuilder) Build() Test {
        return *tb.target
    }
    
    // NewTestBuilder creates a new builder instance for Test
    func NewTestBuilder() TestBuilder {
        return TestBuilder{target: &Test{}}
    }
    ```
4. 이제 다음과 같이 쉽게 만들어 진 함수를 사용할 수 있습니다.
    ```go
    func SomeMethod() {
        test := NewTestBuilder().WithName("Yang").WithAge(25).Build()
    }
    ```
   
# Annotations
| Annotation | Description                                                    |
| --- |----------------------------------------------------------------|
| `@AllArgsConstructor` | 모든 매개변수를 받는 Constructor를 생성합니다.                                |
| `@NoArgsConstructor` | 매개변수가 없는 Constructor를 생성합니다.                                   |
| `@RequiredArgsConstructor` | `validate:"required"` 태그가 붙은 필드만을 매개변수로 받는 Constructor를 생성합니다. |
| `@Builder` | Builder를 생성합니다.                                                |
| `@Getter` | Getter를 생성합니다.                                                 |
| `@Setter` | Setter를 생성합니다.                                                 |
| `@ToString` | ToString 함수를 생성합니다.                                            |
| `@Equals` | Equals 함수를 생성합니다.                                              | 

## Default Constructor
`// @{생성자 어노테이션}.Default`를 통해 해당 생성자를 패키지의 기본 생성자 `New()`로 만들 수 있습니다.
```go
// some_file.go
// Test
// @AllArgsConstructor.Default
type Test struct {
    Name string
    Age  int
}
```
```go
// some_file_gombok.go

// New
func New(name string, age int) Test {
   return Test{
      Name: name,
      Age:  age,
   }
}
```

# Tags
다음 태그를 이용하여 gombok을 통해 생성되는 함수의 동작을 변경할 수 있습니다.

| Tag         | Value | Description                                                                                                        |
|-------------| --- |--------------------------------------------------------------------------------------------------------------------|
| `validate`  | `required` | @RequiredArgsConstructor 어노테이션을 통해 생성되는 필드를 지정할 수 있습니다.                                                            |
| `constructor` | `ignore` | 해당 태그가 지정된 필드의 경우 `@AllArgsConstructor`, `@RequiredArgsConstructor` 어노테이션을 통해 생성되는 Constructor에서 제외됩니다.            |
| `builder`   | `ignore` | 해당 태그가 지정된 필드의 경우 `@Builder` 어노테이션을 통해 생성되는 `Builder`에서 `WithXXX()` 메서드가 생성되지 않습니다.                                |
| `builder` | `must` | 해당 태그가 지정된 필드의 경우 `@Builder` 어노테이션을 통해 생성되는 `WithXXX()` 메서드에서 해당 변수가 zero value이거나 포인터 타입인데 nil인 경우 panic을 발생시킵니다. |
| `getter`    | `ignore` | 해당 태그가 지정된 필드의 경우 `@Getter` 어노테이션을 통해 생성되는 해당 필드의 Getter 메서드가 생성되지 않습니다.                                           |
| `setter`    | `ignore` | 해당 태그가 지정된 필드의 경우 `@Setter` 어노테이션을 통해 생성되는 해당 필드의 Setter 메서드가 생성되지 않습니다.                                           |
| `to_string` | `ignore` | 해당 태그가 지정된 필드의 경우 `@ToString` 어노테이션을 통해 생성되는 `String()` 메서드에서 제외됩니다.                                               |


## Example 
```go
// @Builder
type Test struct {
    Name string `builder:"ignore"`
    Age  int
}
```

## Trouble Shooting 👊

```bash
$ gombok
> zsh: command not found: gombok
```

go로 설치한 프로그램을 실행할 때 발생하는 에러입니다. ~/.zshrc 파일(혹은 ~/.bashrc)의 하단에 다음과 같이 환경변수를 추가합니다.
```bash
# ...
export PATH="$HOME/go/bin:$PATH"
```
