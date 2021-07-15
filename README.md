# Motor control API wrapping


# SWIG

- Go 1.16.5
- SWIG 4.1.0
- Ubuntu 18.04

### SWIG installation

[https://blog.naver.com/chandong83/221130068042](https://blog.naver.com/chandong83/221130068042)

# How to Call Motor API to Golang

!!!!!!!!경로 및 상대경로 확인해서 참조할 것!!!!!!!!

### from Motor API (C++) to Golang Architecture

                                           `src ( .h)` + `include ( .cpp)`

                                             ↗️  (reference to)

`Interface file`       +      `dynamixel_sdk.h (header files)`  

                                  ⬇️ (generate)         

`generated.go (package)`      ⬅️      `main.go` 

                                           (reference to)

### 1. Compile C++ src

[DynamixelSDK.zip](Motor%20control%20API%20wrapping%2070c75e701f5148239cfafb21273f63b8/DynamixelSDK.zip)

- C++ 헤더 코드들을 먼저 컴파일 → 컴파일된 헤더들을 참조하는 방식
- 소스 파일들을 수정하였기에 아래 방법으로 빌드:
    1. 첨부된 압축 파일 다운
    2. 압축 해제 후 :

    ```bash
    cd /DynamixelSDK/c++/build/linux64

    make clean && make

    sudo make uninstall && sudo make install
    ```

### 2. Make Interface file

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

Wrapping과정에서 헤더 파일 수정을 최소화하는 것을 목표

헤더 파일 수정시 새로 컴파일 및 wrapping해야되어 일이 커짐

대부분의 이슈들은 swig interface file에서 해결 가능 

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

C++ 헤더들을 swig의 interface file을 통해 link시켜 Go package를 생성해줌

대부분의 작업들은 interface file에서 이루어짐 

**Interface file 문법**

- `%module example` : package 명을 example로 설정
- `%include header` : 헤더 선언
- `%extend` : 추가로 함수 확장하는 등에 사용
- `%feature` : director 권한 부여시 사용(c++의 상속 개념)
- `%import sdkjfh.i` : sdkjfh.i interface file을 참조
- `%insert` : 추가로 삽입할 내용이 있을 경우 사용
- `%apply` : 기존의 존재하는 변수 맵핑을 이용하여 변수 변환, default
- `%rename (foo) name` : foo를 name으로 이름 변경
- `%ignore print` : ignore all declarations named `print`
- Interface file 내에서 사용하는 헤더 존재 ex) %include <typemaps.i>
- `typemap(gotype)` : default로 제공하는 swig 변수 변환외에도 추가할 수 있음

    ex) %typemap(gotype) bool, const bool & "bool"

**Interface file**

```cpp
// Interface file
// Motor.i

%module(director="1") Motor

%include <typemaps.i>

%{
#include <fcntl.h>
#include <termios.h>
#include <unistd.h>

#include <stdlib.h>
#include <stdio.h>

#include "dynamixel_sdk.h"
#include "port_handler.h"
#include "packet_handler.h"
%}

%feature("director") PortHandlerLinux;
%feature("director") Protocol2PacketHandler;

%apply int {uint8_t}
%apply int {uint16_t}
%apply long long {uint32_t}
%apply double {float}

%include "dynamixel_sdk.h"
%include "port_handler.h"
%include "packet_handler.h"
```

### 3. SWIG Header files with Interface file

- `GOPATH` : 실제 작업한 소스가 존재하는 경로, 유저가 임의로 지정

       * 해당 프로젝트에선 DynamixelSDK/c++/src/dynamixel_sdk/gopath로 지정

       * 지정한 GOPATH, GOROOT를 확인은 `go env` 명령어를 통해 확인  

       * .profile 파일을 수정했는데도 경로가 변하지 않는다면 `export GOPATH=` 방법을 통해 설정 

- `GOROOT` : /usr/local/go/
- `Motor` : 이번 프로젝트에서 사용된 package명
- `SRCDIR` : $GOPATH/src/

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

경로 설정 시 `.profile` , `.bash_profile` , `.zsh_profile` 등을 사용해서 위의 경로를 설정할 수 있지만 한가지 방법만 사용하는 것이 좋음 → 중구난방으로 사용하면 나중에 어떤 방법으로 경로를 설정한지 찾기 어려워 path가 제대로 수정되지 않음

[https://mingrammer.com/go-codelab/install-and-env/](https://mingrammer.com/go-codelab/install-and-env/)

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

```bash
cd $GOPATH/src/Motor/

swig -go -cgo -c++ -intgosize 64 Motor.i

go mod init Motor

go build

cd ..

sudo cp -r Motor /usr/local/go/src/
```

### 4. make reference to Generated.go

- wrapping 후 헤더에 있는 function, class, variable 등등을 참조할 때 대문자를 사용 → Go에서는 지역변수, 전역변수 개념 x → 대문자와 소문자로 구분
- wrapping되면서 함수명이 변경됨 → 변경되는 규칙은 manual 참조

    ex) New, Set, Get, Delete 등이 추가될 수 있음

- 이외의 다른 rule은 SWIG manual 참고

**Main code**

```go
// Main Code
// read_write_main.go

package main

import (
	"Motor"
	"fmt"
)

const (
	ADDR_TORQUE_ENABLE     = 64
	ADDR_GOAL_POSITION     = 116
	ADDR_PRESENT_POSITION  = 132
	ADDR_GOAL_PWM          = 100
	ADDR_GOAL_VELOCITY     = 104
	ADDR_DRIVE_MODE        = 10
	MINIMUM_POSITION_LIMIT = 0
	MAXIMUM_POSITION_LIMIT = 4095
	BAUDRATE               = 57600

	PROTOCOL_VERSION = 2.0

	DXL_ID = 1
	// DXL_ID_2 = 2
	// DXL_ID_3 = 3

	DEVICENAME = "/dev/ttyUSB0"

	TORQUE_ENABLE               = 1
	TORQUE_DISABLE              = 0
	DXL_MOVING_STATUS_THRESHOLD = 20
	ESC_ASCII_VALUE             = 0x1b
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func main() {
	portHandler := Motor.PortHandlerGetPortHandler(DEVICENAME)
	packetHandler := Motor.PacketHandlerGetPacketHandler(PROTOCOL_VERSION)

	var index int = 0
	var (
		dxl_comm_result int = Motor.COMM_TX_FAIL
		dxl_error       int = 0
	)

	var (
		dxl_goal_position          = [2]int64{MINIMUM_POSITION_LIMIT, MAXIMUM_POSITION_LIMIT}
		dxl_present_position int64 = 0

		dxl_goal_PWM      int   = 200
		dxl_goal_Velocity int64 = 200
	)

	// Open Port
	if portHandler.OpenPort() {
		fmt.Println("Succeeded to open the port!")
	} else {
		fmt.Println("Failed to open the port!")
		fmt.Println("Press any key to terminate...")
		Motor.Getch()
	}

	// Set Port Baudrate
	if portHandler.SetBaudRate(BAUDRATE) {
		fmt.Println("Succeeded to change the baudrate!")
	} else {
		fmt.Println("Failed to change the baudrate!")
		fmt.Println("Press any key to terminate...")
		Motor.Getch()
	}

	// Enable Torque
	dxl_comm_result = packetHandler.Write1ByteTxRx(portHandler, DXL_ID, ADDR_TORQUE_ENABLE, TORQUE_ENABLE, &dxl_error)
	if dxl_comm_result != Motor.COMM_SUCCESS {
		fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
	} else if dxl_error != 0 {
		fmt.Println(packetHandler.GetRxPacketError(dxl_error))
	} else {
		fmt.Println("Succeeded enabling DYNAMIXEL Torque.")
	}

	fmt.Println("goal position : ", dxl_goal_position[0])

	for {

		// Write goal position
		dxl_comm_result = packetHandler.Write4ByteTxRx(portHandler, DXL_ID, ADDR_GOAL_POSITION, dxl_goal_position[index], &dxl_error)
		if dxl_comm_result != Motor.COMM_SUCCESS {
			fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
		} else if dxl_error != 0 {
			fmt.Println(packetHandler.GetRxPacketError(dxl_error))
		}

		// Write goal PWM
		dxl_comm_result = packetHandler.Write2ByteTxRx(portHandler, DXL_ID, ADDR_GOAL_PWM, dxl_goal_PWM, &dxl_error)
		if dxl_comm_result != Motor.COMM_SUCCESS {
			fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
		} else if dxl_error != 0 {
			fmt.Println(packetHandler.GetRxPacketError(dxl_error))
		}

		// Write goal Velocity
		dxl_comm_result = packetHandler.Write4ByteTxRx(portHandler, DXL_ID, ADDR_GOAL_VELOCITY, dxl_goal_Velocity, &dxl_error)
		if dxl_comm_result != Motor.COMM_SUCCESS {
			fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
		} else if dxl_error != 0 {
			fmt.Println(packetHandler.GetRxPacketError(dxl_error))
		}

		fmt.Println("ID : ", DXL_ID, "Goal position", dxl_goal_position[index], "Present position", dxl_present_position)

		// Read the Present position
		dxl_comm_result = packetHandler.Read4ByteTxRx(portHandler, DXL_ID, ADDR_PRESENT_POSITION, &dxl_present_position, &dxl_error)
		if dxl_comm_result != Motor.COMM_SUCCESS {
			fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
		} else if dxl_error != 0 {
			fmt.Println(packetHandler.GetRxPacketError(dxl_error))
		}

		fmt.Println("ID : ", DXL_ID, "Goal position : ", dxl_goal_position[index], "Present position : ", dxl_present_position)

		for abs(dxl_goal_position[index]-dxl_present_position) > DXL_MOVING_STATUS_THRESHOLD {
			dxl_comm_result = packetHandler.Read4ByteTxRx(portHandler, DXL_ID, ADDR_PRESENT_POSITION, &dxl_present_position, &dxl_error)
			if dxl_comm_result != Motor.COMM_SUCCESS {
				fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
			} else if dxl_error != 0 {
				fmt.Println(packetHandler.GetRxPacketError(dxl_error))
			}

			fmt.Println("ID : ", DXL_ID, "Goal position : ", dxl_goal_position[index], "Present position : ", dxl_present_position)
		}

		// Switch the Goal position
		if index == 0 {
			index = 1
		} else {
			index = 0
		}
	}

	// Disable Torque
	dxl_comm_result = packetHandler.Write1ByteTxRx(portHandler, DXL_ID, ADDR_TORQUE_ENABLE, TORQUE_DISABLE, &dxl_error)
	if dxl_comm_result != Motor.COMM_SUCCESS {
		fmt.Println(packetHandler.GetTxRxResult(dxl_comm_result))
	} else if dxl_error != 0 {
		fmt.Println(packetHandler.GetRxPacketError(dxl_error))
	} else {
		fmt.Println("Succeeded disabling DYNAMIXEL Torque.")
	}

	// Close Port
	portHandler.ClosePort()
}
```

### 5. go run main.go

```bash
cd $SRCDIR

sudo chmod a+rw /dev/ttyUSB0

go run main.go
```



------
## Reference
[http://www.swig.org/Doc4.0/SWIGDocumentation.html](http://www.swig.org/Doc4.0/SWIGDocumentation.html)

[https://github.com/swig/swig/tree/master/Examples/go](https://github.com/swig/swig/tree/master/Examples/go)

[http://zacg.github.io/blog/2013/06/06/calling-c-plus-plus-code-from-go-with-swig/](http://zacg.github.io/blog/2013/06/06/calling-c-plus-plus-code-from-go-with-swig/)

[https://tv.naver.com/v/16972079](https://tv.naver.com/v/16972079)
---

by The Legendary Dragon🐉