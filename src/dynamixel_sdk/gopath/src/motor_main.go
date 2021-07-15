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
