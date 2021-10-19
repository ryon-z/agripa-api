# agripa-api
아그리파 API
api.agripa.kr

## 1. Swagger 문서 갱신법
- Swagger란
    - Swagger는 자동으로 API문서를 생성해주는 프로그램
    - 소스코드에 변경이 생겼을 시 아래 명령어를 실행하여 문서를 갱신하여야 함
- Swagger cmd 설치
    -  ```
        go get -u github.com/swaggo/swag/cmd/swag
        ```
- Swagger 문서 갱신
    - ```
        swag init
        ```
- Swagger 문서 확인
    - gin API 서버가 실행 중이어야 Swagger 문서 확인 가능
    - 브라우저를 켜고 아래 사이트 방문  
      http://localhost:8080/swagger/index.html
        

## 2. Docker를 사용하여 Local 개발환경 구축
- Docker 설치 및 실행
    - 개발 pc의 운영체제 맞는 Docker를 설치 및 실행

- Docker image 생성
    - 터미널의 현재 위치는 agripa-api 디렉토리이며, 이미지 이름은 agripa-api라고 가정
    - ```
        docker build -t agripa-api .
        ```
- Docker container 생성 및 bash shell 실행
    - 터미널의 현재 위치는 agripa-api 디렉토리라고 가정
    - container의 8080 포트를 개발 pc의 8080 포트와 연결
    - container는 종료 시 자동 삭제
    - 현재 개발 pc 터미널 위치와 container의 /agripa-api 위치를 동기화
    - ```
        # 윈도우
        docker run -it --rm -p 8080:8080 -v %cd%/:/agripa-api agripa-api /bin/bash 

        # 리눅스 및 맥os
        docker run -it --rm -p 8080:8080 -v ./:/agripa-api agripa-api /bin/bash 
        ```

## 3. gin API 서버 로컬 실행
- 그냥 실행
    - ```
        go run main.go 
        ```  
- air로 실행
    - air란
        - 지정된 파일의 소스코드가 변경되었을 때 자동으로 지정된 명령어를 실행하는 프로그램
        - 세부 설정은 .air.conf에서 수정 가능 
    - air 설치
        - ```
            go get -u github.com/cosmtrek/air
            ```
    - air 실행
        - ```
            air
            ```

## 4. gin API 서버 로컬 실행 확인
- 단순 실행 확인
    - 브라우저를 켠다.
    - 주소창에 http://localhost:8080 을 입력한다.
    - 즉시 응답이 오는지 확인한다 (ex. 404 page not found)
    - 즉시 응답이 오지않고, 오랫동안 요청을 기다린다면 API 서버가 정상 실행되지 않은 것
    - gin을 실행하고 있는 터미널에서 요청을 받았는지 확인한다.
    - 받은 요청이 없으면 정상 실행되지 않은 것
- 세부 실행 확인
    - 브라우저를 켠다.
    - 주소창에 라우터 등록한 URL을 입력한다.  
      (ex. http://localhost:8080/media/news?itemCode=111)
    - 이후 내용은 "단순 실행 확인"과 동일하다.
   